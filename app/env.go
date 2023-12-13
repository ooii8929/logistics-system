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
    e.DBUsername = os.Getenv("DB_USERNAME")
    e.DBPassword = os.Getenv("DB_PASSWORD")
    e.DBHost = os.Getenv("DB_HOST")

    dbPortStr := os.Getenv("DB_PORT")
    dbPort, err := strconv.Atoi(dbPortStr)
    if err != nil {
        return err
    }
    e.DBPort = dbPort

    e.DBName = os.Getenv("DB_NAME")

    e.RedisAddr = os.Getenv("REDIS_ADDR")
    e.RedisPassword = os.Getenv("REDIS_PASSWORD")

    // RedisDB remains commented out as per your code

    return nil
}
