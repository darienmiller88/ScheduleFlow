package cronscheduler

import (
	"ScheduleFlow/Backend/utils"
	"ScheduleFlow/Backend/services"

	"github.com/go-co-op/gocron/v2"
)

func NewSchedulerForSendingEmails(
	scheduleService services.ScheduleService, 
	emailSendService services.EmailSendService,
) gocron.Scheduler{
	s, _ := gocron.NewScheduler()
	date := utils.GetSendDateForEmail()

	s.NewJob(
		gocron.WeeklyJob(1, gocron.NewWeekdays(date.Weekday()), gocron.NewAtTimes(gocron.NewAtTime(12,0, 0))),
		gocron.NewTask(
			func() {
			
			},
		),
	)


	return s
}

