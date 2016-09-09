package db

import (
	"log"
	"time"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/7joe7/personalmanager/resources"
)

func synchronize() {
	err := db.Update(func (tx *bolt.Tx) error {
		b := tx.Bucket(resources.DB_DEFAULT_BASIC_BUCKET_NAME)
		lastSync := string(b.Get(resources.DB_LAST_SYNC_KEY))
		if lastSync == "" || isTimeForSync(lastSync) {
			if err := synchronizeHabits(tx); err != nil {
				return err
			}
			if err := b.Put(resources.DB_LAST_SYNC_KEY, []byte(time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006"))); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Unable to synchronize database. %v", err)
	}
}

func isTimeForSync(lastSync string) bool {
	t, err := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", lastSync)
	if err != nil {
		log.Fatalf("Unable to parse last sync time. %v", err)
	}
	return t.Truncate(24 * time.Hour).Before(time.Now().Truncate(24 * time.Hour))
}

func synchronizeHabits(tx *bolt.Tx) error {
	now := time.Now()
	b := tx.Bucket(resources.DB_DEFAULT_HABITS_BUCKET_NAME)
	return b.ForEach(func (k, v []byte) error {
		key := string(k)
		if key == string(resources.DB_LAST_ID_KEY) {
			return nil
		}
		habit := &resources.Habit{}
		if err := json.Unmarshal(v, habit); err != nil {
			return err
		}
		if !habit.Active {
			return nil
		}

		if habit.Deadline.Before(now) {
			if !habit.Done && habit.LastStreakEnd == nil {
				habit.LastStreakEnd = getTimePointer(*habit.Deadline)
				habit.LastStreak = habit.ActualStreak
				habit.ActualStreak = 0
			}
			switch habit.Repetition {
			case "Daily":
				habit.Deadline = getTimePointer(habit.Deadline.Add(24 * time.Hour))
			case "Weekly":
				habit.Deadline = getTimePointer(habit.Deadline.Add(7 * 24 * time.Hour))
			case "Monthly":
				habit.Deadline = getTimePointer(habit.Deadline.AddDate(0, 1, 0))
			}
			habit.Done = false
			habit.Tries += 1
		}

		var err error
		if v, err = json.Marshal(habit); err != nil {
			return err
		}
		return b.Put(k, v)
	})
}

func getTimePointer(t time.Time) *time.Time {
	return &t
}