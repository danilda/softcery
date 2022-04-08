package main

import (
	"context"
	"go.uber.org/zap"
	"os/signal"
	"softcery/internal/app/back/config"
	"softcery/internal/app/back/service"
	"softcery/internal/pkg/log"
	"softcery/internal/pkg/repository"
	"sync"
	"syscall"
)

func main() {
	log.InitLogger()
	config.InitConfig()

	rabbitRepo, err := repository.NewRabbitRepository()
	if err != nil {
		zap.S().Error(err)
		return
	}
	defer rabbitRepo.Close()

	fileStorageRepo := repository.NewLocalFileSystemRepository()

	optimizationService := service.NewOptimizationService(fileStorageRepo, rabbitRepo)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	wg := new(sync.WaitGroup)

	optServiceErrCh := optimizationService.OptimizeImgFromQueue(ctx, wg)
	zap.S().Info("Back stated!")

	select {
	case <-ctx.Done():
		stop()
	case err := <-optServiceErrCh:
		zap.S().Error(err)
	}

	wg.Wait()
	zap.S().Info("Back stopped!")
}
