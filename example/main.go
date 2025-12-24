package main

import (
	"context"
	"log"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	"github.com/go-redis/redis/v8"
	"github.com/mologie/hapredis"
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
