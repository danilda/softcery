package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/h2non/bimg"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"softcery/internal/pkg/entity"
	"softcery/internal/pkg/imgUtil"
	"sync"
)

type FileStoreRepository interface {
	StoreFile(fileName string, byteData []byte) error
}

type QueueRepository interface {
	ImgQueueConsumeChan(ctx context.Context) (<-chan *entity.Image, error)
}

type OptimizationService struct {
	fileStoreRepo   FileStoreRepository
	queueRepository QueueRepository
}

func NewOptimizationService(imgStoreRepo FileStoreRepository, queueRepository QueueRepository) *OptimizationService {
	return &OptimizationService{
		fileStoreRepo:   imgStoreRepo,
		queueRepository: queueRepository,
	}
}

func (s *OptimizationService) OptimizeImgFromQueue(cxt context.Context, wg *sync.WaitGroup) <-chan error {
	wg.Add(1)
	errCh := make(chan error)
	go func() {
		defer wg.Done()

		imgs, err := s.queueRepository.ImgQueueConsumeChan(cxt)
		if err != nil {
			errCh <- err
			return
		}

		workerPoolWg := new(sync.WaitGroup)
		for i := 0; i < viper.GetInt("img-optimization.workers"); i++ {
			workerPoolWg.Add(1)
			go func(imgs <-chan *entity.Image, workerWg *sync.WaitGroup) {
				defer workerWg.Done()
				for img := range imgs {
					s.HandleImgOptimization(img)
				}
			}(imgs, workerPoolWg)
		}
		workerPoolWg.Wait()

		errCh <- errors.New("Img channel closed")
	}()
	return errCh
}

func (s *OptimizationService) HandleImgOptimization(img *entity.Image) {
	originalImg := bimg.NewImage(img.Data)

	for _, scale := range imgUtil.ScalingOptions() {
		resized, err := optimizeImg(originalImg, scale)
		if err != nil {
			zap.S().Errorf("Error during handlePhotos: %v", err)
		}

		fileName := fmt.Sprintf("%s_%d.%s", img.Id, scale, img.Ext)
		if err := s.fileStoreRepo.StoreFile(fileName, resized.Image()); err != nil {
			zap.S().Errorf("")
		}
	}
}

func optimizeImg(image *bimg.Image, scale int) (*bimg.Image, error) {
	size, err := image.Size()
	if err != nil {
		return nil, err
	}

	resized, err := bimg.NewImage(image.Image()).Resize(
		(size.Width*scale)/100,
		(size.Height*scale)/100,
	)
	if err != nil {
		return nil, err
	}

	return bimg.NewImage(resized), nil
}
