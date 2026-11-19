package utils

import (
	"ScheduleFlow/Backend/models"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

func GetNextTuesday() time.Time {
	// Each day of the week is given a number (Sunday = 0, Monday = 1, ... )
	tuesdayConstant := 2
	
	//Subtract the day of the week from tuesday for the absolute distance from tuesday, and add 7
	//to find how far until the following tuesday.
	daysUntilTuesday := (tuesdayConstant - time.Now().Day()) + 7

	//create a date month/day/year that represents the following tuesday from today.
	dateForNextTuesday := time.Now().AddDate(0, 0, daysUntilTuesday)

	//Get the time zone for the east coast
	loc, _ := time.LoadLocation("America/New_York")

	//Finally, create the date for the next tuesday with the
	dateWithTime := time.Date(
		dateForNextTuesday.Year(),
		dateForNextTuesday.Month(),
		dateForNextTuesday.Day(),
		12,
		0,
		0,
		0,
		loc,
	)

	return dateWithTime
}

// WeekOfString returns a "Month D-D, YYYY" style range for the week containing t, where weeks start on Monday.
func WeekOfString(t time.Time) string {
	weekday := int(t.Weekday()) // Sunday = 0, Monday = 1, ..., Saturday = 6

	// Convert to Monday-start offset: how many days to subtract 
	// to reach this week's Monday.
	daysSinceMonday := (weekday + 6) % 7

	monday := t.AddDate(0, 0, -daysSinceMonday)
	sunday := monday.AddDate(0, 0, 6)

	if monday.Month() == sunday.Month() {
		// Same month: "Dec 16-22, 2024"
		return fmt.Sprintf("%s %d-%d, %d", monday.Format("Jan"), monday.Day(), sunday.Day(), monday.Year())
	}

	// Spans two months: "Dec 30, 2024 - Jan 5, 2025"
	return fmt.Sprintf("%s - %s", monday.Format("Jan 2, 2006"), sunday.Format("Jan 2, 2006"))
}

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
	fmt.Fprintf(res, `<p style="color: red; font-weight: bold; font-size: 1.2rem; margin-top: 0.8rem;">%s</p>`, template.HTMLEscapeString(errorText))
}