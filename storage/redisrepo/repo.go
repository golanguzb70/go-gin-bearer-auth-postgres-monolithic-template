package redisrepo


type InMemoryStorageI interface {
	Set(key, value string) error
	SetWithTTL(key, value string, seconds int) error
	Get(key string) (interface{}, error)
	Exists(key string) (interface{}, error)
	Del(key string) (interface{}, error)
	Keys(pattern string) (interface{}, error)
}
