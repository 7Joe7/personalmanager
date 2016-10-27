package operations

import (
	"encoding/json"

	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/utils"
	//"github.com/7joe7/personalmanager/anybar"
)

func initializeBuckets(t resources.Transaction, bucketsToInitialize [][]byte) {
	t.Add(func () error {
		for i := 0; i < len(bucketsToInitialize); i++ {
			if err := t.InitializeBucket(bucketsToInitialize[i]); err != nil {
				return err
			}
		}
		return nil
	})
}

func ensureValues(t resources.Transaction) {
	t.Add(func () error {
		err := t.EnsureEntity(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_REVIEW_SETTINGS_KEY, &resources.Review{Repetition:resources.HBT_REPETITION_WEEKLY, Deadline:utils.GetFirstSaturday()})
		if err != nil {
			return err
		}
		v, err := json.Marshal([]resources.ActivePort{})
		if err != nil {
			return err
		}
		err = t.EnsureValue(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ANYBAR_ACTIVE_PORTS, v)
		if err != nil {
			return err
		}
		return t.EnsureEntity(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_STATUS_KEY, &resources.Status{})
	})
}