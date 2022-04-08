package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config-back")

	err := viper.ReadInConfig()
	panicOnErr(err)

	config := BackConfig{}
	err = viper.Unmarshal(&config)
	panicOnErr(err)

	validate := validator.New()
	if err = validate.Struct(&config); err != nil {
		panicOnErr(err)
	}
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

type BackConfig struct {
	Scale struct {
		Options []int `yaml:"options" validate:"required"`
	}
	Rabbit struct {
		Queue struct {
			Img struct {
				Name string `yaml:"name" validate:"required"`
			}
		}
	}
}
