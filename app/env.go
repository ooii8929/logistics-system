package main

import (
    "os"
    "strconv"
)

type Environment struct {
    DBUsername    string
    DBPassword    string
    DBHost        string
    DBPort        int
    DBName        string
    RedisAddr     string
    RedisPassword string
    RedisDB       int
}

func (e *Environment) load() error {
    e.DBUsername = os.Getenv("root")
    e.DBPassword = os.Getenv("my-mysql-password")
    e.DBHost = os.Getenv("my-mysql")

	// local
	// e.DBPassword = os.Getenv("123456")
    // e.DBHost = os.Getenv("my-mysql")

    dbPort, err := strconv.Atoi(os.Getenv("3306"))
    if err != nil {
        return err
    }
    e.DBPort = dbPort

    e.DBName = os.Getenv("logistics")

    e.RedisAddr = os.Getenv("redis-server:6379")
	e.RedisPassword = os.Getenv("my-myredis-password")

    redisDB, err := strconv.Atoi(os.Getenv("0"))
    if err != nil {
        return err
    }
    e.RedisDB = redisDB

    return nil
}
