package driver

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
)

type Config struct {
	Host        string
	DB          int
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

type RedisDriver struct {
	pool   *redis.Pool
	config *Config
}

func (d *RedisDriver) Init(config *Config) {
	pool := &redis.Pool{
		MaxIdle:     config.MaxIdle,
		MaxActive:   config.MaxActive,
		IdleTimeout: config.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Host)
			if err != nil {
				log.Error(err)
			} else {
				// select redis db
				c.Do("SELECT", config.DB)
			}
			return c, err
		},
	}

	d.pool = pool
	d.config = config
}

func (d *RedisDriver) Execute(f func(conn redis.Conn) error) error {
	conn := d.pool.Get()
	defer conn.Close()
	return f(conn)
}
