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

	config := FrontConfig{}
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

type FrontConfig struct {
	Server struct {
		Host string `yaml:"host" validate:"required"`
		Port int    `yaml:"port" validate:"required"`
	}
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
