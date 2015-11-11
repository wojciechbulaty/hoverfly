package main

import (
	"fmt"
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
)

const prefix = "genproxy:"

type Cache struct {
	pool *redis.Pool
}

// set records a key in cache (redis)
func (c *Cache) set(key string, value interface{}) error {

	client := c.pool.Get()
	defer client.Close()

	_, err := client.Do("SET", fmt.Sprintf(prefix+key), value)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"key":   fmt.Sprintf(prefix + key),
		}).Error("Failed to SET key...")
	} else {
		log.WithFields(log.Fields{
			"key": fmt.Sprintf(prefix + key),
		}).Info("Key/value SET successfuly!")
	}

	return err
}

// get returns key from cache
func (c *Cache) get(key string) (interface{}, error) {

	client := c.pool.Get()
	defer client.Close()

	value, err := client.Do("GET", fmt.Sprintf(prefix+key))

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"key":   fmt.Sprintf(prefix + key),
		}).Error("Failed to GET key...")
	} else {
		log.WithFields(log.Fields{
			"key": fmt.Sprintf(prefix + key),
		}).Info("Key found!")
	}

	return value, err
}

// getRedisPool returns thread safe Redis connection pool
func getRedisPool() *redis.Pool {

	// getting redis connection
	maxConnections := 10
	mc := os.Getenv("MaxConnections")
	if mc != "" {
		maxCons, err := strconv.Atoi(mc)
		if err != nil {
			maxConnections = 10
		} else {
			maxConnections = maxCons
		}
	}
	// getting redis client for state storing
	redisPool := redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", AppConfig.redisAddress)

		if err != nil {
			log.WithFields(log.Fields{"Error": err.Error()}).Panic("Failed to create Redis connection pool!")
			return nil, err
		}
		if AppConfig.redisPassword != "" {
			if _, err := c.Do("AUTH", AppConfig.redisPassword); err != nil {
				log.WithFields(log.Fields{
					"Error":        err.Error(),
					"PasswordUsed": AppConfig.redisPassword,
				}).Panic("Failed to authenticate to Redis!")
				c.Close()
				return nil, err
			} else {
				log.Info("Authenticated to Redis successfully! ")
			}
		}

		return c, err
	}, maxConnections)

	return redisPool
}
