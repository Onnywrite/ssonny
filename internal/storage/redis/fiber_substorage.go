package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/Onnywrite/ssonny/internal/storage/repo"
	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
)

const namespaceSeparator = ":"

// fiberStorage implements the fiber.Storage interface using
// Redis as the backend. It allows you to organize data
// within a Redis instance by using namespaces,
// effectively creating separate storage areas within
// the same Redis database.
type fiberStorage struct {
	rdb       *redis.Client
	namespace string
}

func (r *RedisStorage) FiberSubstorage(namespace string) fiber.Storage {
	return &fiberStorage{
		rdb:       r.db,
		namespace: namespace,
	}
}

func (f *fiberStorage) namespaced(key string) string {
	return f.namespace + namespaceSeparator + key
}

// Get retrieves a value from the sub-storage.
func (f *fiberStorage) Get(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	bytes, err := f.rdb.Get(ctx, f.namespaced(key)).Bytes()
	if err == redis.Nil {
		return nil, nil // Key not found
	}

	if err != nil {
		return nil, fmt.Errorf("error getting key %s: %w", key, err)
	}

	return bytes, nil
}

// Set stores a value in the sub-storage with an optional expiration.
func (f *fiberStorage) Set(key string, val []byte, exp time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	err := f.rdb.Set(ctx, f.namespaced(key), val, exp).Err()
	if err != nil {
		return fmt.Errorf("error setting key %s: %w", key, err)
	}

	return nil
}

// Delete removes a value from the sub-storage.
func (f *fiberStorage) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	err := f.rdb.Del(ctx, f.namespaced(key)).Err()
	if err != nil {
		return fmt.Errorf("error deleting key %s: %w", key, err)
	}

	return nil
}

// Exists checks if a key exists in the sub-storage.
func (f *fiberStorage) Exists(key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	exists, err := f.rdb.Exists(ctx, f.namespaced(key)).Result()
	if err != nil {
		return false, fmt.Errorf("error checking if key %s exists: %w", key, err)
	}

	return exists > 0, nil
}

// Reset clears all data from the sub-storage.
func (f *fiberStorage) Reset() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	iter := f.rdb.Scan(ctx, 0, f.namespaced("*"), 100).Iterator()
	for iter.Next(ctx) {
		err := f.rdb.Del(ctx, iter.Val()).Err()
		if err != nil {
			return fmt.Errorf("%w: error deleting key %s: %w", repo.ErrInternal, iter.Val(), err)
		}
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("%w: error scanning keys: %w", repo.ErrInternal, err)
	}

	return nil
}

// Close does nothing.
func (f *fiberStorage) Close() error {
	return nil
}
