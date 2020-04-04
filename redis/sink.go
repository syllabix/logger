package redis

import (
	redigo "github.com/garyburd/redigo/redis"
)

// A Pool is used to retrieve pooled connections via it's Get method. When the pool
// is no longer in use - Close should be called to ensure all resources are released
type Pool interface {
	Get() redigo.Conn
	Close() error
}

// Sink can be used to sync zap logs with redis
type Sink struct {
	pool Pool
	key  string
}

func (s *Sink) Write(p []byte) (n int, err error) {
	conn := s.pool.Get()
	defer conn.Close()

	_, err = conn.Do("RPUSH", s.key, string(p))
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

// Close closes the underlying connection pull. All
// future calls to a closed instance will fail
func (s *Sink) Close() error {
	return s.pool.Close()
}

// NewSink construct a useful instance of zap logger sync
// that can be used to write logs to a remote redis instance
func NewSink(key string, pool Pool) *Sink {
	return &Sink{
		pool: pool,
		key:  key,
	}
}
