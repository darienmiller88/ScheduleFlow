package cronscheduler

import (
	"ScheduleFlow/Backend/services"
	"ScheduleFlow/Backend/utils"
	"log"
	"github.com/google/uuid"

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
			func() error{
				schedulesResult := scheduleService.GetAllSchedules()

				if schedulesResult.Err != nil {
					return schedulesResult.Err
				}

				for _, schedule := range schedulesResult.ResultData {
					log.Printf("sche %v", schedule)
				}

				return nil
			},
		),
		gocron.WithName("send-weekly-schedules"),
		gocron.WithEventListeners(
			gocron.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
				log.Printf("job %q failed: %v", jobName, err)
			}),
			gocron.AfterJobRunsWithPanic(func(jobID uuid.UUID, jobName string, recoverData any) {
				log.Printf("job %q panicked: %v", jobName, recoverData)
			}),
		),
	)


	return s
}

