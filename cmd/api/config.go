package api

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

type config struct {
	Addr string `mapstructure:"ADDR"`
	Env  string `mapstructure:"ENV"`
	DB   struct {
		Driver      string `mapstructure:"DB_DRIVER"`
		DSN         string `mapstructure:"DB_DSN"`
		AutoMigrate bool   `mapstructure:"DB_AUTO_MIGRATE"`
	} `mapstructure:",squash"`
}

func newConfig(cfgFile string) config {
	var cfg config
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.Getwd()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}
	setDefaults()
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("failed to read configuration:", err.Error())
			os.Exit(1)
		}
	}
	for _, cfgName := range viper.AllKeys() {
		if err := viper.BindEnv(cfgName); err != nil {
			fmt.Println("bind env key:", cfgName, "has errors:", err)
		}
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Println("failed to unmarshal configuration:", err.Error())
		os.Exit(1)
	}
	return cfg
}

func setDefaults() {
	viper.SetDefault("ADDR", "localhost:3000")
	viper.SetDefault("ENV", "development")
	viper.SetDefault("DB_DRIVER", "postgres")
	viper.SetDefault("DB_DSN", "postgresql://postgres:password@localhost/pd?sslmode=disable")
	viper.SetDefault("DB_AUTO_MIGRATE", "true")
}
