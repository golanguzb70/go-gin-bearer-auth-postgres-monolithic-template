package redis

import "github.com/gomodule/redigo/redis"

type RedisRepo struct {
	Rds *redis.Pool
}

func NewRedisRepo(rds *redis.Pool) InMemoryStorageI {
	return &RedisRepo{
		Rds: rds,
	}
}

func (r *RedisRepo) Exists(key string) (interface{}, error) {
	conn := r.Rds.Get()
	defer conn.Close()

	return conn.Do("EXISTS", key)
}

func (r *RedisRepo) Set(key, value string) error {
	conn := r.Rds.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)

	return err
}

func (r *RedisRepo) SetWithTTL(key, value string, seconds int) error {
	conn := r.Rds.Get()
	defer conn.Close()

	_, err := conn.Do("SETEX", key, seconds, value)
	return err
}

func (r *RedisRepo) Get(key string) (interface{}, error) {
	conn := r.Rds.Get()
	defer conn.Close()

	return conn.Do("GET", key)
}

func (r *RedisRepo) Del(key string) (interface{}, error) {
	conn := r.Rds.Get()
	defer conn.Close()

	return conn.Do("DEL", key)
}

func (r *RedisRepo) Keys(pattern string) (interface{}, error) {
	conn := r.Rds.Get()
	defer conn.Close()

	return conn.Do("KEYS", pattern)
}
