# hapredis

A simple adapter to use [go-redis](https://github.com/go-redis/redis) as storage backend for [brutella/hap](https://github.com/brutella/hap).

## Usage

Download the library within the context of your project:

```sh
go get github.com/brutella/hap
go get github.com/go-redis/redis/v8
go get github.com/mologie/hapredis
```

Integrate it by creating a Redis client and wrapping into a Store interface.

It is important to give your client a unique prefix to allow multiple different
instances to connect to use a single Redis database.

```go
package main

import (
	"context"
	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	"github.com/go-redis/redis/v8"
	"github.com/mologie/hapredis"
	"log"
)

func main() {
	ctx := context.Background()

	// create storage backend using this library
	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	store := hapredis.NewStore(ctx, redisClient, "unique-instance-prefix:")

	// create accessory and its server
	acc := accessory.NewSwitch(accessory.Info{Name: "Lamp"})
	server, err := hap.NewServer(store, acc.A)
	if err != nil {
		// Redis errors can get you here, e.g. when the server cannot load its
		// UUID due to connection or permission problems.
		log.Panic(err)
	}

	// run accessory server forever
	server.ListenAndServe(ctx)
}
```
