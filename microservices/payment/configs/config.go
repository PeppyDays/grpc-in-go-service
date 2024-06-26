package configs

import (
	"log"
	"os"
	"strconv"
)

func GetEnv() string {
	return getEnvironmentValue("ENV")
}

func GetDataSourceURL() string {
	return getEnvironmentValue("DATA_SOURCE_URL")
}

func GetApplicationPort() int {
	port, err := strconv.Atoi(getEnvironmentValue("APPLICATION_PORT"))
	if err != nil {
		log.Fatalf("port %s is invalid", getEnvironmentValue("APPLICATION_PORT"))
	}

	return port
}

func getEnvironmentValue(key string) string {
	if os.Getenv(key) == "" {
		log.Fatalf("%s environment variable is missing", key)
	}

	return os.Getenv(key)
}
