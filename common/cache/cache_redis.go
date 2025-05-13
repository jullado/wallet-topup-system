package cache

import (
	"context"
	"fmt"
	"log"
	"time"
	"wallet-topup-system/config"

	"github.com/redis/go-redis/v9"
)

type appCache struct {
	cache *redis.Client
}

func NewAppRedisCache() AppCache {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// parse redis uri
	opts, err := redis.ParseURL(config.Env.CacheURI)
	if err != nil {
		panic(err)
	}

	// create redis client
	client := redis.NewClient(opts)
	if err := client.Ping(ctx); err != nil && err.String() != "ping: PONG" {
		panic(err)
	}

	return &appCache{client}
}

func (c *appCache) Get(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := c.cache.Get(ctx, key).Bytes()
	if err != nil {

		// if key does not exist
		if err == redis.Nil {
			return nil, nil
		}

		return nil, err
	}

	return res, nil
}

func (c *appCache) Set(key string, val []byte, exp time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.cache.Set(ctx, key, val, exp).Err()
}

func (c *appCache) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.cache.Del(ctx, key).Err()
}

func (c *appCache) ExpiredEvent(callback func(key string) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// subscribe to the Redis event
	pubsub, err := c.subscribeToExpiredEvent(ctx)
	if err != nil {
		return err
	}

	// listen for events and process them
	return c.processExpiredEvent(pubsub, callback)
}

func (c *appCache) subscribeToExpiredEvent(ctx context.Context) (*redis.PubSub, error) {
	// subscribe to the expired key event
	pubsub := c.cache.Subscribe(ctx, fmt.Sprintf("__keyevent@%v__:expired", c.cache.Options().DB))

	// wait for the subscription to be set up
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to expired event: %w", err)
	}

	return pubsub, nil
}

func (c *appCache) processExpiredEvent(pubsub *redis.PubSub, callback func(key string) error) error {
	ch := pubsub.Channel()

	for event := range ch {
		if err := callback(event.Payload); err != nil {
			log.Printf("error processing expired event for key %s: %v", event.Payload, err)
		}
	}

	// if the channel is closed or an error occurs -> return the error
	if err := pubsub.Close(); err != nil {
		return fmt.Errorf("error closing pubsub: %w", err)
	}

	return nil
}
