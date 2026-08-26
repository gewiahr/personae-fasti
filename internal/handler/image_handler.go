package handler

import (
	"errors"
	"net/http"

	"personae-fasti/internal/dto"
	e "personae-fasti/internal/pkg/errorutils"
	"personae-fasti/internal/pkg/httputils"
	"personae-fasti/internal/service"
)

type ImageHandler struct {
	svc *service.ImageService
}

func NewImageHandler(svc *service.ImageService) *ImageHandler {
	return &ImageHandler{svc: svc}
}

func (h *ImageHandler) List(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	images, err := h.svc.List(req.Context, req.Player, req.Request.PathValue("type"), req.Request.PathValue("ext"))
	if err != nil {
		return e.ErrToApiError(err)
	}
	return httputils.Response{Status: http.StatusOK, Body: dto.ImageList{Images: images}}
}

func (h *ImageHandler) CreateExternal(req httputils.RequestData[dto.ExternalImageCreate]) httputils.Responder {
	image, err := h.svc.CreateExternal(req.Context, req.Player, req.Request.PathValue("type"), req.Request.PathValue("ext"), req.Body.URL)
	if err != nil {
		return e.ErrToApiError(err)
	}
	return httputils.Response{Status: http.StatusCreated, Body: image}
}

func (h *ImageHandler) Upload(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	if err := req.Request.ParseMultipartForm(1 << 20); err != nil {
		if req.Request.MultipartForm != nil {
			_ = req.Request.MultipartForm.RemoveAll()
		}
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			return e.NewApiError(http.StatusRequestEntityTooLarge, "Файл слишком большой")
		}
		return e.NewApiError(http.StatusBadRequest, "Некорректная форма загрузки")
	}
	if req.Request.MultipartForm != nil {
		defer req.Request.MultipartForm.RemoveAll()
	}
	file, header, err := req.Request.FormFile("file")
	if err != nil {
		return e.NewApiError(http.StatusBadRequest, "Выберите файл изображения")
	}
	defer file.Close()
	image, err := h.svc.Upload(req.Context, req.Player, req.Request.PathValue("type"), req.Request.PathValue("ext"), file, header.Size)
	if err != nil {
		return e.ErrToApiError(err)
	}
	return httputils.Response{Status: http.StatusCreated, Body: image}
}

func (h *ImageHandler) SetMain(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	image, err := h.svc.SetMain(req.Context, req.Player, req.Request.PathValue("imageExt"))
	if err != nil {
		return e.ErrToApiError(err)
	}
	return httputils.Response{Status: http.StatusOK, Body: image}
}

func (h *ImageHandler) Delete(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	if err := h.svc.Delete(req.Context, req.Player, req.Request.PathValue("imageExt")); err != nil {
		return e.ErrToApiError(err)
	}
	return httputils.Response{Status: http.StatusOK, Body: nil}
}

func (h *ImageHandler) GetQuota(req httputils.RequestData[dto.NoBody]) httputils.Responder {
	quota, err := h.svc.GetQuota(req.Context, req.Player)
	if err != nil {
		return e.ErrToApiError(err)
	}
	return httputils.Response{Status: http.StatusOK, Body: quota}
}
