package service

import (
	"fmt"
	"github.com/gofrs/uuid"
	"os"
	"softcery/internal/pkg/entity"
)

type QueueRepository interface {
	PushInImgQueue(body []byte) error
}

type FileStorageRepository interface {
	FindFileByPattern(fileName string) (*os.File, error)
}

type ImgManageService struct {
	queue           QueueRepository
	fileStorageRepo FileStorageRepository
}

func NewImgManageService(queue QueueRepository,
	fileStorageRepo FileStorageRepository) *ImgManageService {
	return &ImgManageService{
		queue:           queue,
		fileStorageRepo: fileStorageRepo,
	}
}

func (s *ImgManageService) GetImg(uuid uuid.UUID, quality int) (*entity.Image, error) {
	pattern := fmt.Sprintf("%s_%d.*", uuid.String(), quality)
	file, err := s.fileStorageRepo.FindFileByPattern(pattern)
	if err != nil {
		return nil, fmt.Errorf("Error during geting img by id and quality: %s", err)
	}

	return entity.NewImage(file, file.Name()), nil
}

func (s *ImgManageService) PushImgInQueue(img *entity.Image) error {
	return s.queue.PushInImgQueue(img.ToBytes())
}
