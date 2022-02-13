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
	dbUser := GetEnv("MYSQL_USER")
	dbPasswd := GetEnv("MYSQL_PASSWORD")
	dbHost := GetEnv("MYSQL_HOST")
	dbName := GetEnv("MYSQL_DB")

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPasswd, dbHost, dbName)

	return dsn
}

func GetTestDbDsn() string {
	dbUser := GetEnv("MYSQL_USER_TEST")
	dbPasswd := GetEnv("MYSQL_PASSWORD_TEST")
	dbHost := GetEnv("MYSQL_HOST_TEST")
	dbName := GetEnv("MYSQL_DB_TEST")

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPasswd, dbHost, dbName)

	return dsn
}

func GetRelativeDirPath() string {
	_, f, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(f))

	return dir
}
