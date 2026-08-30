package utils

import (
	"ScheduleFlow/Backend/models"
	"fmt"
	"html/template"
	"net/http"
)

// Helper function to allow repos to send result payloads with less text.
func GetResult[T any](err error, statusCode int, payload T) models.Result[T] {
	return models.Result[T]{
		StatusCode: statusCode,
		Err:        err,
		ResultData: payload,
	}
}

func SendHtmlError(res http.ResponseWriter, statusCode int, errorText string) {
    res.Header().Set("Content-Type", "text/html") 
    res.WriteHeader(statusCode)
    
	fmt.Fprintf(res, `<p class="error-text">%s</p>`, template.HTMLEscapeString(errorText))
}