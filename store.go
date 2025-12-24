package hapredis

import (
	"context"
	"slices"
	"strings"

	"github.com/brutella/hap"
	"github.com/go-redis/redis/v8"
)

// Store is a Redis storage adapter for hap.Server.
type Store struct {
	ctx    context.Context
	client *redis.Client
	prefix string
}

// NewStore creates a Redis adapter that conforms to hap.Store.
//
// Set its context to the server's context to cancel stuck Redis operations when
// the server is shut down, or to context.Background() when no ordinary shutdown
// is implemented.
//
// The client instance must be created through go-redis.
//
// The prefix must be unique for your hap server instance. It is convention to
// use a colon symbol as separator in Redis keys. The prefix length is
// irrelevant with the amount of data stored by hap. An example for a short and
// nice prefix is "fooapp:barhost:", where barhost is a hostname.
func NewStore(ctx context.Context, client *redis.Client, prefix string) *Store {
	return &Store{ctx: ctx, client: client, prefix: prefix}
}

func (s Store) key(key string) string {
	return s.prefix + key
}

func (s Store) Set(key string, value []byte) error {
	return s.client.Set(s.ctx, s.key(key), value, 0).Err()
}

func (s Store) Get(key string) ([]byte, error) {
	return s.client.Get(s.ctx, s.key(key)).Bytes()
}

func (s Store) Delete(key string) error {
	return s.client.Del(s.ctx, s.key(key)).Err()
}

func (s Store) KeysWithSuffix(suffix string) ([]string, error) {
	match := s.prefix + "*" + redisEscape(suffix)
	keyMap := make(map[string]bool)
	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = s.client.Scan(s.ctx, cursor, match, 0).Result()
		if err != nil {
			return nil, err
		}
		for _, key := range keys {
			keyMap[key] = true
		}
		if cursor == 0 {
			break
		}
	}
	var result []string
	for key := range keyMap {
		result = append(result, key)
	}
	slices.Sort(result)
	return result, nil
}

// redisEscape escapes glob string characters mentioned in Redis' KEYS docs.
func redisEscape(s string) (result string) {
	for _, c := range s {
		if strings.ContainsRune("?*[]", c) {
			result += "\\"
		}
		result += string(c)
	}
	return
}

// Test that Store conforms to hap.Store at compile-time.
var _ hap.Store = (*Store)(nil)
