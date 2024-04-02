package gsb

import (
	"os"
	"reflect"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func getTypeName(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func getStrEnv(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defaultValue
	}
	return val
}

func getIntEnv(key string, defaultValue int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	ret, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("Expected environment variable, %s, to be an integer. Got, %s", key, val)
	}
	return ret
}
