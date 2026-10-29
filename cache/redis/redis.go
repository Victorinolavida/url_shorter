package redis

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"url_shorter/cache"
	"url_shorter/config"
	"url_shorter/internal/entity"

	redisCache "github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

const defaultTTL = time.Hour

type RedisCache struct {
	cache *redisCache.Cache
}

func NewRedisClient(conf config.Redis) (cache.Cache, error) {
	addr, err := validateConfig(conf)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{Addr: addr})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	cache := redisCache.New(&redisCache.Options{Redis: client})
	return &RedisCache{cache: cache}, nil
}

func validateConfig(conf config.Redis) (string, error) {
	host := strings.TrimSpace(conf.Host)
	if host == "" {
		return "", fmt.Errorf("redis host is required")
	}

	port, err := strconv.Atoi(strings.TrimSpace(conf.Port))
	if err != nil || port < 1 || port > 65535 {
		return "", fmt.Errorf("invalid redis port %q", conf.Port)
	}

	return net.JoinHostPort(host, strconv.Itoa(port)), nil
}

func (r *RedisCache) Set(ctx context.Context, link *entity.Link) error {
	ttl := defaultTTL
	if link.ExpiresAt != nil {
		ttl = time.Until(*link.ExpiresAt)
		if ttl < time.Second {
			return r.cache.Delete(ctx, link.Code)
		}
	}

	return r.cache.Set(
		&redisCache.Item{
			Ctx:   ctx,
			Key:   link.Code,
			Value: link,
			TTL:   ttl,
		},
	)
}

func (r *RedisCache) Get(ctx context.Context, code string) (*entity.Link, error) {
	link := new(entity.Link)
	if err := r.cache.Get(ctx, code, link); err != nil {
		if errors.Is(err, redisCache.ErrCacheMiss) {
			return nil, cache.ErrCacheMiss
		}
		return nil, err
	}
	return link, nil
}

func (r *RedisCache) Delete(ctx context.Context, code string) error {
	return r.cache.Delete(ctx, code)
}
