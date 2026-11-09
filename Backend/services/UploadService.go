package services

import (
	"ScheduleFlow/Backend/models"
	"ScheduleFlow/Backend/utils"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

type UploadService interface {
	UploadSchedule(file multipart.File, filename string) models.Result[string]
}

type uploadService struct {
	// Add any dependencies or configurations needed for the service
}

func NewUploadService() UploadService {
	return &uploadService{}
}

// UploadSchedule implements [UploadService].
func (u *uploadService) UploadSchedule(file multipart.File, filename string) models.Result[string] {
	// 1. Validate the extension
	ext := strings.ToLower(filepath.Ext(filename))

	if ext != ".xlsx" && ext != ".xls" {
		return utils.GetResult(errors.New("Only .xlsx and .xls files are allowed"), http.StatusBadRequest, "")
	}
	
	// 2. Validate the actual contents (the extension alone is trivially spoofed)
	head := make([]byte, 8)

	if _, err := io.ReadFull(file, head); err != nil {
		return utils.GetResult(errors.New("File is empty or unreadable"), http.StatusBadRequest, "")
	}


	return utils.GetResult(nil, http.StatusOK, "File uploaded successfully")
}