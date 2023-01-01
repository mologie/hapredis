package hapredis

import (
	"context"
	"errors"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStore(t *testing.T) {
	req := require.New(t)

	db, mock := redismock.NewClientMock()
	store := NewStore(context.Background(), db, "prefix:")

	mock.ExpectGet("prefix:foo").RedisNil()
	r, err := store.Get("foo")
	req.NoError(mock.ExpectationsWereMet())
	req.NoError(err)
	req.Nil(r)

	mock.ExpectSet("prefix:foo", []byte("bar"), 0).SetVal("OK")
	mock.ExpectGet("prefix:foo").SetVal("bar")
	req.NoError(store.Set("foo", []byte("bar")))
	r, err = store.Get("foo")
	req.NoError(mock.ExpectationsWereMet())
	req.NoError(err)
	req.Equal([]byte("bar"), r)

	expectedErr := errors.New("FAIL")
	mock.ExpectGet("prefix:err").SetErr(expectedErr)
	r, err = store.Get("err")
	req.NoError(mock.ExpectationsWereMet())
	req.ErrorIs(err, expectedErr)
	req.Nil(r)

	mock.ExpectDel("prefix:foo").SetVal(1)
	req.NoError(store.Delete("foo"))
	req.NoError(mock.ExpectationsWereMet())

	// Deletion of keys is unspecified and not implemented consistently in hap.
	// mem_store ignores unknown keys during deletion, fs_store returns the
	// error from os.Remove, which will fail when the destination file is gone.
	mock.ExpectDel("prefix:baz").SetVal(0)
	req.NoError(store.Delete("baz"))
	req.NoError(mock.ExpectationsWereMet())

	const match = `prefix:*f\?o`
	mock.ExpectScan(0, match, 0).SetVal([]string{"a", "b"}, 42)
	mock.ExpectScan(42, match, 0).SetVal([]string{"c", "a"}, 0)
	keys, err := store.KeysWithSuffix("f?o")
	req.NoError(mock.ExpectationsWereMet())
	req.NoError(err)
	req.Equal([]string{"a", "b", "c"}, keys)
}

func TestRedisEscape(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  string
	}{
		{input: "foo", want: "foo"},
		{input: "f??", want: `f\?\?`},
		{input: "f*x?x", want: `f\*x\?x`},
		{input: "f[^*]x?x", want: `f\[^\*\]x\?x`},
	} {
		escaped := redisEscape(tc.input)
		if escaped != tc.want {
			t.Errorf("escaping %v should result in %v, but got %v", tc.input, tc.want, escaped)
		}
	}
}
