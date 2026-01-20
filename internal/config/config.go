package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`

	Postgres struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
	} `mapstructure:"postgres"`

	Redis struct {
		Addr     string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`

	JWT struct {
		AccessSecret  string `mapstructure:"accesssecret"`
		RefreshSecret string `mapstructure:"refreshsecret"`
		AccessTTL     string `mapstructure:"accessttl"`
		RefreshTTL    string `mapstructure:"refreshttl"`
	} `mapstructure:"jwt"`

	Minio struct {
		EndPoint  string `mapstructure:"endpoint"`
		AccessKey string `mapstructure:"accesskey"`
		SecretKey string `mapstructure:"secretkey"`
		Bucket    string `mapstructure:"bucket"`
		UseSSL    bool   `mapstructure:"usessl"`
		Domain    string `mapstructure:"domain"`
	} `mapstructure:"minio"`
}

var App Config

func Init() {
	configDir := getConfigDir()
	log.Printf("Config directory: %s", configDir)

	envPath := filepath.Join(configDir, ".env")
	log.Printf("Looking for .env at: %s", envPath)

	if _, err := os.Stat(envPath); err != nil {
		log.Fatalf("no such file")
	}

	if err := godotenv.Load(envPath); err != nil {
		log.Println("No env file found, using environment variables only")
	}

	fmt.Printf(".env file : %s", envPath)
	v := viper.New()

	v.SetDefault("server.port", "8080")

	v.SetDefault("postgres.host", "localhost")
	v.SetDefault("postgres.port", "5432")
	v.SetDefault("postgres.user", "postgres")
	v.SetDefault("postgres.password", "password")
	v.SetDefault("postgres.dbname", "user_service")
	v.SetDefault("postgres.sslmode", "disable")

	v.SetDefault("redis.addr", "localhost:6379")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)

	v.SetDefault("jwt.accesssecret", "change_me")
	v.SetDefault("jwt.refreshsecret", "change_me")
	v.SetDefault("jwt.accessttl", "15m")
	v.SetDefault("jwt.refreshttl", "720h")

	v.SetDefault("minio.endpoint", "localhost:9000")
	v.SetDefault("minio.accesskey", "minioadmin")
	v.SetDefault("minio.secretkey", "minioadmin")
	v.SetDefault("minio.bucket", "user-photos")
	v.SetDefault("minio.usessl", false)
	v.SetDefault("minio.domain", "localhost:9000")

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.AutomaticEnv()

	if err := v.Unmarshal(&App); err != nil {
		log.Fatalf("failed to write config: %v", err)
	}
	log.Printf("successfully set config! ")
	log.Printf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		App.Postgres.Host,
		App.Postgres.Port,
		App.Postgres.User,
		App.Postgres.Password,
		App.Postgres.DBName,
		App.Postgres.SSLMode)
}

func getConfigDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)
}
