package config

import (
	"github.com/spf13/viper"
	"log"
)

var (
	LogLevel           = "info"
	LogEncoding        = "console"
	LogFile            = false
	LogFileMaxSize     = 500
	LogFilePath        = "/tmp/"
	HttpPort           = "8080"
	DatabaseURL        = "postgres://postgres:postgres@localhost:5432/openapi?sslmode=disable"
	SessionKey         = ""
	RMQUrl             = "amqp://guest:guest@localhost:5672"
	LogFilePathsToSeek []string
)

func LoadConfig() error {
	viper.SetConfigType("json")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("log-aggregator")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.log-aggregator")
	viper.AddConfigPath("/etc/log-aggregator")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/app")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		return err
	}

	setDefaults()

	LogLevel = viper.GetString("LOG_LEVEL")
	LogEncoding = viper.GetString("LOG_ENCODING")
	LogFile = viper.GetBool("LOG_FILE")
	LogFileMaxSize = viper.GetInt("LOG_FILE_MAX_SIZE")
	LogFilePath = viper.GetString("LOG_FILE_PATH")
	HttpPort = viper.GetString("HTTP_PORT")
	DatabaseURL = viper.GetString("DATABASE_URL")
	SessionKey = viper.GetString("SESSION_KEY")
	RMQUrl = viper.GetString("RMQ_URL")
	LogFilePathsToSeek = viper.GetStringSlice("LOG_FILE_PATHS_TO_SEEK")

	return nil
}

func setDefaults() {
	viper.SetDefault("LOG_LEVEL", LogLevel)
	viper.SetDefault("LOG_ENCODING", LogEncoding)
	viper.SetDefault("LOG_FILE", LogFile)
	viper.SetDefault("LOG_FILE_MAX_SIZE", LogFileMaxSize)
	viper.SetDefault("LOG_FILE_PATH", LogFilePath)
	viper.SetDefault("HTTP_PORT", HttpPort)
	viper.SetDefault("DATABASE_URL", DatabaseURL)
	viper.SetDefault("SESSION_KEY", SessionKey)
	viper.SetDefault("RMQ_URL", RMQUrl)
	viper.SetDefault("LOG_FILE_PATHS_TO_SEEK", LogFilePathsToSeek)
}
