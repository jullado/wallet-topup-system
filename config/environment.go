package config

import (
	"log"
	"reflect"
	"time"

	"github.com/spf13/viper"
)

// Set default environments
var Env = struct {
	// app
	Env          string        `mapstructure:"ENV"`
	Port         string        `mapstructure:"PORT"`
	Cors         string        `mapstructure:"CORS"`
	ReadTimeout  time.Duration `mapstructure:"READ_TIMEOUT"`
	WriteTimeout time.Duration `mapstructure:"WRITE_TIMEOUT"`
	TimeoutTopUp time.Duration `mapstructure:"TIMEOUT_TOP_UP"`
	APIKey       string        `mapstructure:"API_KEY"`

	// third party
	DBURI    string `mapstructure:"DB_URI"`
	CacheURI string `mapstructure:"CACHE_URI"`
}{
	Env:          "production",
	Port:         "3000",
	Cors:         "*",
	ReadTimeout:  1 * time.Minute,
	WriteTimeout: 1 * time.Minute,

	TimeoutTopUp: 1 * time.Minute,
}

func NewAppInitEnvironment() {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Environment variables not used from .env")
	}

	// try load settings from env vars
	r := reflect.TypeOf(Env)
	for i := range r.NumField() {
		f := r.Field(i).Tag.Get("mapstructure")
		viper.BindEnv(f)
	}

	if err := viper.Unmarshal(&Env); err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}
}
