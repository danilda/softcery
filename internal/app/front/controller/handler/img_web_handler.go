package handler

import (
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mime/multipart"
	"net/http"
	"softcery/internal/app/front/service"
	"softcery/internal/pkg/entity"
	"softcery/internal/pkg/imgUtil"
	"strconv"
)

type ImgWebApiHandler struct {
	imgManageService *service.ImgManageService
}

func NewImgWebApiHandler(imgManageService *service.ImgManageService) *ImgWebApiHandler {
	return &ImgWebApiHandler{imgManageService: imgManageService}
}

func (h *ImgWebApiHandler) Upload(res http.ResponseWriter, req *http.Request) {
	defer closeBody(req, "img upload")

	maxSizeImg := viper.GetInt64("scale.maxImgSize")
	if err := req.ParseMultipartForm(maxSizeImg << 20); err != nil {
		http.Error(res, "Too large img", http.StatusBadRequest)
		return
	}

	file, handler, err := req.FormFile("file")
	if err != nil {
		zap.S().Errorf("Error during reading file from upload request %v", err)
		http.Error(res, "Unable to read file", http.StatusBadRequest)
		return
	}
	defer closeFile(file)

	img := entity.NewImage(file, handler.Filename)
	err = h.imgManageService.PushImgInQueue(img)
	if err != nil {
		zap.S().Errorf("Error during pushing img in queue: %v", err)
		http.Error(res, "error during uploading", http.StatusInternalServerError)
		return
	}

	if _, err = res.Write([]byte(img.Id)); err != nil {
		zap.S().Errorf("Error during responding for img upload reqeust: %v", err)
	}
}

func (h *ImgWebApiHandler) Download(res http.ResponseWriter, req *http.Request) {
	imgId, err := uuid.FromString(chi.URLParam(req, "uuid"))
	if err != nil {
		http.Error(res, "Invalid img id", http.StatusBadRequest)
		return
	}
	zap.S().Debugf("Download uuid: %v", imgId)

	qualities, ok := req.URL.Query()["quality"]
	if !ok || len(qualities[0]) < 1 {
		http.Error(res, "Parameter 'qualities' is missing", http.StatusBadRequest)
		return
	}

	qualityStr := qualities[0]
	quality, err := strconv.Atoi(qualityStr)
	if err != nil {
		http.Error(res, "Invalid parameter 'qualities'", http.StatusBadRequest)
		return
	} else {
		if !imgUtil.IsValidScalingOption(quality) {
			http.Error(res, "Unexpected parameter 'qualities'", http.StatusBadRequest)
			return
		}
	}

	img, err := h.imgManageService.GetImg(imgId, quality)
	if err != nil {
		http.Error(res, "Unable to resolve img with specified id", http.StatusNotFound)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/octet-stream")

	if _, err = res.Write(img.Data); err != nil {
		zap.S().Errorf("Error during responding for img upload reqeust: %v", err)
	}
}

func closeBody(req *http.Request, methodName string) {
	err := req.Body.Close()
	if err != nil {
		zap.S().Errorf("Error during closing body of %s: %v ", methodName, req)
	}
}

func closeFile(file multipart.File) {
	err := file.Close()
	if err != nil {
		zap.S().Errorf("Error during closing file: %v", err)
	}
}
