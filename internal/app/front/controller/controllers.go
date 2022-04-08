package controller

import (
	"github.com/go-chi/chi"
	"net/http"
	"softcery/internal/app/front/controller/handler"
	"softcery/internal/app/front/service"
)

func InitControllers(imgManageService *service.ImgManageService) *chi.Mux {
	r := chi.NewRouter()

	InitImgApi(r, handler.NewImgWebApiHandler(imgManageService))

	return r
}

type ImgApiController interface {
	Upload(res http.ResponseWriter, req *http.Request)
	Download(res http.ResponseWriter, req *http.Request)
}

func InitImgApi(r *chi.Mux, ctr ImgApiController) {
	r.Route("/file/api", func(r chi.Router) {
		r.Post("/upload", ctr.Upload)
		r.Get("/download/{uuid}", ctr.Download)
	})
}
