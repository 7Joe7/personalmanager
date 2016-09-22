package operations

import (
	"testing"
	"time"

	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/utils"
	"github.com/7joe7/personalmanager/test"
)

func TestCountScoreChange(t *testing.T) {
	timeSpent := utils.GetDurationPointer(time.Duration(int64(86400000000000)))
	testTask1 := &resources.Task{Name: "test1", Note: "note1", BasePoints: 13, TimeSpent: timeSpent }
	test.ExpectInt(18850, countScoreChange(testTask1), t)

	timeSpent = utils.GetDurationPointer(time.Duration(int64(43200000000000)))
	testTask2 := &resources.Task{Name: "test2", Note: "note2", Project: nil, BasePoints: 2, TimeSpent: timeSpent }
	test.ExpectInt(1460, countScoreChange(testTask2), t)
}

func TestStopProgress(t *testing.T) {
	timeSpent := utils.GetDurationPointer(time.Duration(int64(86400000000000)))
	testTask1 := &resources.Task{
		Name: "test1",
		Note: "note1",
		BasePoints: 13,
		TimeSpent: timeSpent,
		InProgress: true,
		InProgressSince: utils.GetTimePointer(time.Now().Add(time.Duration(int64(-43200000000000))))}
	stopProgress(testTask1)
	test.ExpectBool(false, testTask1.InProgress, t)
	test.ExpectInt(28210, countScoreChange(testTask1), t)

	testTask2 := &resources.Task{
		Name: "test2",
		Note: "note2",
		BasePoints: 7,
		InProgress: true,
		InProgressSince: utils.GetTimePointer(time.Now().Add(time.Duration(int64(-43200000000000))))}
	stopProgress(testTask2)
	test.ExpectBool(false, testTask2.InProgress, t)
	test.ExpectInt(5110, countScoreChange(testTask2), t)
}

