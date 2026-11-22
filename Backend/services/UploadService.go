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

const UploadDir string = "uploads/schedules"

type UploadService interface {
	UploadSchedule(file multipart.File, specialistId string) models.Result[string]
	GetSchedule(specialistId string) (string, error)
}

type uploadService struct {
	// Add any dependencies or configurations needed for the service
}


func NewUploadService() UploadService {
	return &uploadService{}
}

// GetSchedule implements [UploadService].
func (u *uploadService) GetSchedule(specialistId string) (string, error) {
	pathName := filepath.Join(UploadDir, specialistId, specialistId+".xlsx")
	_, err := os.Open(pathName)

	if err != nil{
		if errors.Is(err, os.ErrNotExist){
			return "", errors.New("no schedule found")
		}

		return "", err
	}

	return pathName, nil
}

// UploadSchedule implements [UploadService].
func (u *uploadService) UploadSchedule(file multipart.File, specialistId string) models.Result[string] {

	//Use excelize to open the uploaded file and check if it is a valid Excel file. If it is not, return an error.
	f, err := excelize.OpenReader(file)

	if err != nil {
		return utils.GetResult(errors.New("Error opening Excel file, or not a valid Excel file"), http.StatusBadRequest, "")
	}

	defer f.Close()

	// First, we need to create a new file in the uploads/schedules directory with a unique name.
	dst, err := createFile(specialistId)

	if err != nil {
		return utils.GetResult(err, http.StatusInternalServerError, "")
	}

	defer dst.Close()

	// Use the excel file to write the contents of the uploaded file to the new file in the
	// "uploads/schedules/<specialistId>/<specialistId>.xlsx" directory.
	if err := f.Write(dst); err != nil {
		return utils.GetResult(errors.New("Error reading file"), http.StatusInternalServerError, "")
	}

	return utils.GetResult(nil, http.StatusOK, "")
}

// Helper function to create a new directory for a specialists schedule
func createFile(specialistId string) (*os.File, error) {
	// Create an uploads directory if it doesn’t exist
	if err := os.MkdirAll(filepath.Join(UploadDir, specialistId), 0755); err != nil {
		return nil, err
	}

	// Create an empty excel file in the uploads/schedules directory with a unique name.
	dst, err := os.Create(filepath.Join(UploadDir, specialistId, specialistId+".xlsx"))

	if err != nil {
		return nil, err
	}

	return dst, nil
}
