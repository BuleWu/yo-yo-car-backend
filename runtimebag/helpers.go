package runtimebag

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func GetEnvString(key, defaultValue string) string {
	fmt.Println("Key: ", key)
	fmt.Println("Val: ", os.Getenv(key))

	if x := os.Getenv(key); x != "" {
		return x
	}
	return defaultValue
}

func GetEnvInt(key string, defaultValue int64) int64 {
	if x := os.Getenv(key); x != "" {
		v, err := strconv.ParseInt(x, 10, 64)
		if err != nil {
			log.Panicf("Can not convert string to int64 %s", err.Error())
		}
		return int64(v)
	}
	return defaultValue
}

func GetEnvBool(key string, defaultValue bool) bool {
	if x := os.Getenv(key); x != "" {
		v, err := strconv.ParseBool(x)
		if err != nil {
			log.Panicf("Can not convert string to bool %s", err.Error())
		}
		return v
	}
	return defaultValue
}
