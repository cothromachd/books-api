package repo

import (
	"context"
	"time"

	"github.com/cothromachd/books-api/internal/config"
	"github.com/cothromachd/books-api/internal/entity"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	rc *redis.Client
}

func NewRedisCache(cfg *config.Config) *RedisCache {
	rc := redis.NewClient(&redis.Options{
		Addr:     cfg.RDB.Conn,
		Password: "",
		DB:       0,
	})

	return &RedisCache{rc: rc}
}

func (c *RedisCache) HasBook(id string) (bool, error) {
	exists, err := c.rc.Exists(context.Background(), id).Result()
	if err != nil {
		return false, err
	}

	if exists == 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func (c *RedisCache) GetBook(id string) (entity.Book, error) {
	bookJson, err := c.rc.Get(context.Background(), id).Result()
	if err != nil {
		return entity.Book{}, err
	}

	book, err := entity.Unmap(bookJson)
	if err != nil {
		return entity.Book{}, err
	}

	return book, nil
}

func (c *RedisCache) SetBook(id string, book entity.Book) error {
	bookJson, err := book.Map()
	if err != nil {
		return err
	}

	err = c.rc.Set(context.Background(), id, bookJson, time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisCache) DeleteBook(id string) error {
	err := c.rc.Del(context.Background(), id).Err()
	if err != nil {
		return err
	}

	return nil
}
