package main

import (
    "fmt"
    "log"
    "time"

    "github.com/7joe7/personalmanager/operations"
    "github.com/7joe7/personalmanager/resources"
    "github.com/everdev/mack"
    "runtime/debug"
    "os"
    "os/signal"
    "syscall"
    "github.com/7joe7/personalmanager/db"
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
    log.Println("habits:", habits)
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
            log.Printf("alarm time: %v - id: %v", at, id)
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
        if nearestAlarm != nil {
            log.Printf("waiting for timer, reset or quit of alarm worker - %v\n", nearestAlarm.time)
        } else {
            log.Printf("waiting for timer, reset or quit of alarm worker - %v\n", nil)
        }
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

func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Panicked. %v %s\n", r, string(debug.Stack()))
            log.Fatalf("Panicked. %v %s", r, string(debug.Stack()))
            os.Exit(3)
        }
    }()

    // initialize catching the interrupt and terminate signal
    c := make(chan os.Signal)
    signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

    resources.Alr = NewAlarmManager()

    go resources.Alr.Run()
    log.Println(1)

    db.Open()
    log.Println(2)
    t := db.NewTransaction()
    log.Println(3)
    operations.InitializeBuckets(t)
    log.Println(4)
    operations.EnsureValues(t)
    log.Println(5)
    //operations.Synchronize(t)
    //log.Println(6)
    t.Execute()
    log.Println(7)

    resources.Alr.Sync()
    log.Println(8)

    // waiting for the signals
    <-c
    signal.Stop(c)
    close(c)
}
