package alarm

import (
	"fmt"
	"log"
	"time"

	"github.com/7joe7/personalmanager/operations"
	"github.com/7joe7/personalmanager/resources"
	"github.com/everdev/mack"
)

type alTime struct {
	time  *time.Time
	title string
}

type AlManager struct {
	times      map[string]*alTime
	resetTimer chan struct{}
	Quit       chan struct{}
}

func NewAlarmManager() *AlManager {
	return &AlManager{times: map[string]*alTime{}, resetTimer: make(chan struct{}), Quit: make(chan struct{})}
}

func (am *AlManager) Sync() {
	habits := operations.FilterHabits(func(h *resources.Habit) bool {
		return h.Active && !h.Done && h.AlarmTime != nil
	})
	am.times = map[string]*alTime{}
	for _, h := range habits {
		am.times[fmt.Sprintf("%s-habit", h.AlarmTime.Format(resources.DATE_HOUR_MINUTE_FORMAT))] = &alTime{time: h.AlarmTime, title: "Time to grow"}
	}
	log.Println("syncing alarms")
	am.resetTimer <- struct{}{}
}

func (am *AlManager) Run() {
	timer := time.NewTimer(time.Hour * 24)
	for {
		var nearestAlarm *alTime
		now := time.Now()
		for id, at := range am.times {
			if at.time.Before(now) {
				delete(am.times, id)
				continue
			}
			if nearestAlarm == nil {
				nearestAlarm = at
				continue
			}
			if at.time.Before(*nearestAlarm.time) {
				nearestAlarm = at
			}
		}
		if nearestAlarm != nil {
			timer.Reset(nearestAlarm.time.Sub(time.Now()))
		}
		log.Println("waiting for timer, reset or quit of alarm worker")
		select {
		case <-timer.C:
			if nearestAlarm != nil {
				err := mack.Notify(nearestAlarm.title)
				if err != nil {
					log.Printf("unable to notify: %v\n", err)
				}
			}
		case <-am.resetTimer:
		case <-am.Quit:
			return
		}
	}
}
