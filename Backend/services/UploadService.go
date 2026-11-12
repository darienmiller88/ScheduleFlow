package services

import (
	"ScheduleFlow/Backend/models"
	"ScheduleFlow/Backend/utils"
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const uploadDir string = "uploads/schedules"

var (
	xlsxMagic = []byte{'P', 'K', 0x03, 0x04}                              // .xlsx is a zip archive
	xlsMagic  = []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1} // legacy .xls (OLE2)
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
	
	// 2. Validate the actual contents (the extension alone is trivially spoofed). Read the 
	// first 8 bytes for the magic number of the file to determine if it is a valid Excel file.
	head := make([]byte, 8)

	if _, err := io.ReadFull(file, head); err != nil {
		return utils.GetResult(errors.New("File is empty or unreadable"), http.StatusBadRequest, "")
	}

	//Check to see if the file is a valid Excel file based on its magic number prefix
	validXlsx := ext == ".xlsx" && bytes.HasPrefix(head, xlsxMagic)
	validXls := ext == ".xls" && bytes.HasPrefix(head, xlsMagic)

	if !validXlsx && !validXls {
		return utils.GetResult(errors.New("File does not appear to be a valid Excel file"), http.StatusBadRequest, "")
	}
	
	// rewind to the beginning after sniffing
	if _, err := file.Seek(0, io.SeekStart); err != nil { 
		return utils.GetResult(errors.New("Could not process file"), http.StatusInternalServerError, "")
	}


	// 3. Write to disk under a server-generated name
	if err := os.MkdirAll(uploadDir, 0o750); err != nil {
		return utils.GetResult(errors.New("Could not save file"), http.StatusInternalServerError, "")
	}
	
	destPath := filepath.Join(uploadDir, uuid.NewString() + ext)
	dst, err := os.OpenFile(destPath, os.O_WRONLY | os.O_CREATE | os.O_EXCL, 0o640)

	if err != nil {
		return utils.GetResult(errors.New("Could not save file"), http.StatusInternalServerError, "")
	}

	// Copy the file contents to the destination file
	if _, err := io.Copy(dst, file); err != nil {
		dst.Close()
		os.Remove(destPath)
		return utils.GetResult(errors.New("Could not save file"), http.StatusInternalServerError, "")
	}

	// Close can surface buffered write errors
	if err := dst.Close(); err != nil { 
		os.Remove(destPath)
		return utils.GetResult(errors.New("Could not save file"), http.StatusInternalServerError, "")
	}

	//return the destination path of the uploaded file for further processing (e.g., parsing the schedule)
	return utils.GetResult(nil, http.StatusOK, destPath)
}