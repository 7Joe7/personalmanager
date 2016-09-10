package operations

import (
	"time"

	"github.com/7joe7/personalmanager/db"
	"github.com/7joe7/personalmanager/resources"
	"log"
)

func synchronize(t *db.Transaction) {
	t.Add(func () error {
		lastSync := string(t.GetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_LAST_SYNC_KEY))
		if lastSync == "" || isTimeForSync(lastSync) {
			var scoreChange *int
			habit := &resources.Habit{}
			if err := t.MapEntities(resources.DB_DEFAULT_HABITS_BUCKET_NAME, habit, GetSyncHabitFunc(habit, scoreChange)); err != nil {
				return err
			}
			status := &resources.Status{}
			if err := t.ModifyEntity(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_STATUS_KEY, status, GetSyncStatusFunc(status, *scoreChange)); err != nil {
				return err
			}
			if err := t.SetValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_LAST_SYNC_KEY, []byte(time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006"))); err != nil {
				return err
			}
		}
		return nil
	})
}

func initializeBuckets(t *db.Transaction, bucketsToInitialize [][]byte) {
	t.Add(func () error {
		for i := 0; i < len(bucketsToInitialize); i++ {
			if err := t.InitializeBucket(bucketsToInitialize[i]); err != nil {
				return err
			}
		}
		return nil
	})
}

func ensureValues(t *db.Transaction) {
	t.Add(func () error {
		return t.EnsureEntity(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_STATUS_KEY, &resources.Status{})
	})
}

func isTimeForSync(lastSync string) bool {
	t, err := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", lastSync)
	if err != nil {
		log.Fatalf("Unable to parse last sync time. %v", err)
	}
	return t.Truncate(24 * time.Hour).Before(time.Now().Truncate(24 * time.Hour))
}