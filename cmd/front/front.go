package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"softcery/internal/app/front/config"
	"softcery/internal/app/front/controller"
	"softcery/internal/app/front/service"
	"softcery/internal/pkg"
	"softcery/internal/pkg/log"
	"softcery/internal/pkg/repository"
	"syscall"
)

func main() {
	log.InitLogger()
	config.InitConfig()

	rabbitRepository, err := repository.NewRabbitRepository()
	if err != nil {
		zap.S().Error(err)
		return
	}
	defer rabbitRepository.Close()

	fileStorageRepo := repository.NewLocalFileSystemRepository()

	imgManageService := service.NewImgManageService(rabbitRepository, fileStorageRepo)

	router := controller.InitControllers(imgManageService)

	srv := new(pkg.Server)
	go func() {
		zap.S().Info("Front is up!")

		if err := srv.Run(viper.GetString("server.port"), router); err != nil {
			panic(fmt.Sprintf("Error when starting the http server: %s", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		zap.S().Errorf("error occured on server shutting down: %s", err.Error())
	}
}
