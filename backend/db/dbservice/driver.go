package main

import (
	"github.com/garyburd/redigo/redis"
	"gopkg.in/mgo.v2"
	"log"
	"os"
	"time"
)

type driver struct {
	redisClient  *redis.Pool
	mongoSession *mgo.Session
}

type redisConfig struct {
	host        string
	db          int
	maxIdle     int
	maxActive   int
	idleTimeout time.Duration
}

func (d *driver) init(minfo *mgo.DialInfo, rcfg *redisConfig) {
	// init mongodb client
	var err error
	d.mongoSession, err = mgo.DialWithInfo(minfo)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	d.mongoSession.SetMode(mgo.Monotonic)

	// init redis client with pool
	d.redisClient = &redis.Pool{
		MaxIdle:     rcfg.maxIdle,
		MaxActive:   rcfg.maxActive,
		IdleTimeout: rcfg.idleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", rcfg.host)
			if err != nil {
				log.Fatal(err)
				os.Exit(-1)
			}
			// select redis db
			c.Do("SELECT", rcfg.db)
			return c, nil
		},
	}
}
