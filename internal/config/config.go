package config

import (
	"flag"
	"fmt"
	"miniorder-order-service/internal/domain"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func GetServiceProperties() (svcProperties domain.ServiceProperties) {
	isLocalDev := flag.Bool("local", false, "=(true/false)")
	flag.Parse()

	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	if *isLocalDev {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("error reading transaction file: %w", err))
		}
	} else {
		viper.SetConfigFile(".env.default")
		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("error reading transaction file: %w", err))
		}
		viper.SetConfigFile(".env.svc")
		if err := viper.MergeInConfig(); err != nil {
			panic(fmt.Errorf("error merging the .env.svc file: %w", err))
		}
	}

	svcProperties = GetEnvServiceProperties()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		svcProperties = GetEnvServiceProperties()
	})
	return svcProperties
}

func GetEnvServiceProperties() (svcProperties domain.ServiceProperties) {
	fmt.Println("Starting Load Env " + time.Now().Format("2006-01-02 15:04:05"))

	svcProperties = domain.ServiceProperties{
		DebugMode:      viper.GetBool("DEBUG_MODE"),
		ServiceName:    viper.GetString("SERVICE_NAME"),
		Timeout:        viper.GetDuration("TIMEOUT"),
		PoolConnection: viper.GetInt("POOL_CONNECTION"),
		Redis: domain.RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Password: viper.GetString("REDIS_PASSWORD"),
			Database: viper.GetInt("REDIS_DATABASE"),
		},
		Database: domain.DatabaseConfig{
			IP:          viper.GetString("DB_IP"),
			Port:        viper.GetInt("DB_PORT"),
			User:        viper.GetString("DB_USER"),
			Password:    viper.GetString("DB_PASSWORD"),
			Name:        viper.GetString("DB_NAME"),
			DialTimeout: viper.GetDuration("DB_DIAL_TIMEOUT"),
		},
	}

	if err := validator.New().Struct(svcProperties); err != nil {
		panic(err)
	}

	fmt.Println("Finish Load Env " + time.Now().Format("2006-01-02 15:04:05"))
	return svcProperties
}
