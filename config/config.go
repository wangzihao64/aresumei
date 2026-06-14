package config

import (
	"aresumei/cache"
	"aresumei/dao"
	"strings"

	"gopkg.in/ini.v1"
)

var (
	AppModel      string
	HttpPort      string
	DB            string
	DbHost        string
	DbPort        string
	DbUser        string
	DbPassword    string
	DbName        string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	RedisPoolSize int
	SmtpHost      string
	SmtpEmail     string
	SmtpPass      string
)

func LoadRedis(file *ini.File) {
	RedisHost = file.Section("redis").Key("RedisHost").String()
	RedisPort = file.Section("redis").Key("RedisPort").String()
	RedisPassword = file.Section("redis").Key("RedisPassword").String()
	RedisDB = file.Section("redis").Key("RedisDB").MustInt()
	RedisPoolSize = file.Section("redis").Key("RedisPoolSize").MustInt()
}
func LoadServer(file *ini.File) {
	AppModel = file.Section("service").Key("AppModel").String()
	HttpPort = file.Section("service").Key("HttpPort").String()
}
func LoadMysql(file *ini.File) {
	DB = file.Section("mysql").Key("DB").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassword = file.Section("mysql").Key("DbPassword").String()
	DbName = file.Section("mysql").Key("DbName").String()
}
func LoadEmail(file *ini.File) {
	SmtpHost = file.Section("email").Key("SmtpHost").String()
	SmtpEmail = file.Section("email").Key("SmtpEmail").String()
	SmtpPass = file.Section("email").Key("SmtpPass").String()
}
func Init() {
	file, err := ini.Load("./config/config.ini")
	if err != nil {
		panic(err)
	}
	LoadServer(file)
	LoadMysql(file)
	//mysql 主
	pathRead := strings.Join([]string{DbUser, ":", DbPassword, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8mb4&parseTime=true"}, "")
	pathWrite := strings.Join([]string{DbUser, ":", DbPassword, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8mb4&parseTime=true"}, "")
	//mysql 从
	//
	if err := dao.Database(pathRead, pathWrite); err != nil {
		panic(err)
	}
	LoadRedis(file)
	if err := cache.InitRedis(cache.RedisConfig{
		Host:     RedisHost,
		Port:     RedisPort,
		Password: RedisPassword,
		DB:       RedisDB,
		PoolSize: RedisPoolSize,
	}); err != nil {
		panic(err)
	}
	LoadEmail(file)
}
