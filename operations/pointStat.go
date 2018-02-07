package operations

import (
	"fmt"

	"strconv"
	"time"

	"github.com/7joe7/personalmanager/db"
	"github.com/7joe7/personalmanager/resources"
)

func getPointStats() map[string]*resources.PointStat {
	pointStats := map[string]*resources.PointStat{}
	var values map[string]string
	tr := db.NewTransaction()
	tr.Add(
		func() error {
			var err error
			values, err = tr.GetValues(resources.DB_DEFAULT_POINTS_BUCKET_NAME)
			return err
		})
	tr.Execute()
	for id, value := range values {
		idTime, err := time.Parse("2006-01-02", id)
		if err != nil {
			panic(err)
		}
		valueInt, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
		pointStats[fmt.Sprintf("%sPS", id)] = &resources.PointStat{Value: valueInt, Id: idTime}
	}
	return pointStats
}
