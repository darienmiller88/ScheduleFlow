package services

import (
	"ScheduleFlow/Backend/models"
	"ScheduleFlow/Backend/utils"
	"errors"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/xuri/excelize/v2"
)

const uploadDir string = "uploads/schedules"

type UploadService interface {
	UploadSchedule(file multipart.File, specialistName string, specialistId string) models.Result[string]
}

type uploadService struct {
	// Add any dependencies or configurations needed for the service
}

func NewUploadService() UploadService {
	return &uploadService{}
}

func createFile(specialistName string, specialistId string) (*os.File, error) {
	// Create an uploads directory if it doesn’t exist
	if err := os.MkdirAll(filepath.Join(uploadDir, specialistId), 0755); err != nil {
		return nil, err
	}

	// Create an empty excel file in the uploads/schedules directory with a unique name.
	dst, err := os.Create(filepath.Join(uploadDir, specialistId, specialistId+".xlsx"))

	if err != nil {
		return nil, err
	}

	return dst, nil
}

// UploadSchedule implements [UploadService].
func (u *uploadService) UploadSchedule(file multipart.File, specialistName string, specialistId string) models.Result[string] {

	//Use excelize to open the uploaded file and check if it is a valid Excel file. If it is not, return an error.
	f, err := excelize.OpenReader(file)

	if err != nil {
		return utils.GetResult(errors.New("Error opening Excel file, or not a valid Excel file"), http.StatusBadRequest, "")
	}

	defer f.Close()

	// First, we need to create a new file in the uploads/schedules directory with a unique name.
	dst, err := createFile(specialistName, specialistId)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, "")
	}

	defer dst.Close()

	// The destination will receive the contents of the uploaded file. We can use ReadFrom to copy the
	// contents (bytes) from the uploaded file to the destination file, which will be
	// "uploads/schedules/<specialistId>/<specialistId>.xlsx"
	if err := f.Write(dst); err != nil {
		return utils.GetResult(errors.New("Error reading file"), http.StatusInternalServerError, "")
	}

	return utils.GetResult(nil, http.StatusOK, "")
}
