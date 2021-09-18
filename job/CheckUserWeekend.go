package job

import (
	"fmt"
	"polkovnik/adapter/notifyChannel"
	"polkovnik/domain"
	"time"
)

type jobCheckUserWeekendEntry struct {
	User *domain.User
}

type jobCheckUserWeekendData struct {
	HolidayStart []jobCheckUserWeekendEntry
	HolidayEnd   []jobCheckUserWeekendEntry
	Date         string
	CountStart   int
	CountEnd     int
}

func (p Processor) CheckUserWeekend(team *domain.Team, task *domain.Task, story *domain.History, channel notifyChannel.Interface, dateNow time.Time) error {
	settings := task.GetTaskCheckUserWeekendSettingsDto()

	var dateChek = dateNow
	duration, err := time.ParseDuration(settings.GetDateModifyDuration())
	if err != nil {
		return err
	}
	dateChek = dateChek.Add(duration)

	story.AddLine(fmt.Sprintf("Date now is %s", dateNow.String()))
	story.AddLine(fmt.Sprintf("Date check is %s", dateChek.String()))

	// Находим тех, кто собирается в отпуск и возвращается из него
	var holidayStart []jobCheckUserWeekendEntry
	var holidayEnd []jobCheckUserWeekendEntry
	for _, user := range team.Users {
		if user.Active == false {
			story.AddLine(fmt.Sprintf("User %s is not active", user.Name))
			continue
		}

		if !user.Weekend.IsHoliday(dateNow) && user.Weekend.IsHoliday(dateChek) {
			holidayStart = append(holidayStart, jobCheckUserWeekendEntry{User: user})
			story.AddLine(fmt.Sprintf("User %s go holiday", user.Name))
		}

		if user.Weekend.IsHoliday(dateNow) && !user.Weekend.IsHoliday(dateChek) {
			holidayEnd = append(holidayEnd, jobCheckUserWeekendEntry{User: user})
			story.AddLine(fmt.Sprintf("User %s return from holiday", user.Name))
		}
	}

	if len(holidayStart) == 0 && len(holidayEnd) == 0 {
		story.AddLine("All users is worked")
		return nil
	}

	message, err := channel.CreateMessageFromTemplate(
		"CheckUserWeekend.html",
		jobCheckUserWeekendData{
			HolidayStart: holidayStart,
			HolidayEnd:   holidayEnd,
			Date:         dateChek.Format("02.01.2006"),
			CountStart:   len(holidayStart),
			CountEnd:     len(holidayEnd),
		},
	)
	if err != nil {
		return err
	}

	story.AddLine("Send message:")
	story.AddLine(message.Text)
	_, err = channel.SendTeamMessage(message)
	if err != nil {
		return err
	}

	story.AddLine("Message sent")

	return nil
}
