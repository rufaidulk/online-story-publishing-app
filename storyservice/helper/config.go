package helper

import (
	"fmt"
	"log"
	"path"
	"runtime"

	"github.com/spf13/viper"
)

func GetEnv(key string) string {
	viper.SetConfigFile(fmt.Sprintf("%s/../.env", GetRelativeDirPath()))
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	value, ok := viper.Get(key).(string)
	if !ok {
		log.Printf("ENV key: %s\n", key)
		log.Fatal("invalid type assertion")
	}

	return value
}

func GetDsn() string {
	dbUser := GetEnv("MONGO_USER")
	dbPasswd := GetEnv("MONGO_PASSWORD")
	dbHost := GetEnv("MONGO_HOST")

	uri := fmt.Sprintf("mongodb://%s:%s@%s:27017/admin", dbUser, dbPasswd, dbHost)

	return uri
}

func GetRelativeDirPath() string {
	_, f, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(f))

	return dir
}
