package operations

import (
	"testing"
	"time"

	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/test"
	"github.com/7joe7/personalmanager/utils"
	"github.com/stretchr/testify/assert"
)

func TestCountScoreChange(t *testing.T) {
	timeSpent := utils.GetDurationPointer(time.Duration(int64(86400000000000)))
	timeEstimate := utils.GetDurationPointer(time.Duration(int64(43200000000000)))
	testTask1 := &resources.Task{Name: "test1", Note: "note1", BasePoints: 13, TimeSpent: timeSpent, TimeEstimate: timeEstimate}
	assert.Equal(t, 9490, testTask1.CountScoreChange())

	timeSpent = utils.GetDurationPointer(time.Duration(int64(43200000000000)))
	timeEstimate = utils.GetDurationPointer(time.Duration(int64(86400000000000)))
	testTask2 := &resources.Task{Name: "test2", Note: "note2", Project: nil, BasePoints: 2, TimeSpent: timeSpent, TimeEstimate: timeEstimate}
	assert.Equal(t, 2900, testTask2.CountScoreChange())
}

func TestStopProgress(t *testing.T) {
	resources.Abr = test.NewAnybarManagerMock()
	timeSpent := utils.GetDurationPointer(time.Duration(int64(86400000000000)))
	testTask1 := &resources.Task{
		Name:            "test1",
		Note:            "note1",
		BasePoints:      13,
		TimeSpent:       timeSpent,
		InProgress:      true,
		InProgressSince: utils.GetTimePointer(time.Now().Add(time.Duration(int64(-43200000000000))))}
	stopProgress(testTask1)
	assert.False(t, testTask1.InProgress)
	assert.Equal(t, 130, testTask1.CountScoreChange())

	testTask2 := &resources.Task{
		Name:            "test2",
		Note:            "note2",
		BasePoints:      7,
		InProgress:      true,
		InProgressSince: utils.GetTimePointer(time.Now().Add(time.Duration(int64(-43200000000000)))),
		TimeEstimate:    utils.GetDurationPointer(time.Duration(int64(86400000000000)))}
	stopProgress(testTask2)
	assert.False(t, testTask2.InProgress)
	assert.Equal(t, 10150, testTask2.CountScoreChange())
}
