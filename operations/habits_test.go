package operations

import (
	"testing"
	"time"

	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/utils"
	"github.com/stretchr/testify/assert"
)

var (
	tomorrowDeadline    = utils.GetTimePointer(time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour))
	tomorrowDeadlineStr = tomorrowDeadline.Format(resources.DATE_FORMAT)
)

func TestGetNumberOfMissedDeadlines(t *testing.T) {
	h := getActiveHabit("habit1", resources.HBT_REPETITION_DAILY, 7, 3, 3, 0, 7)
	h.Deadline = utils.GetTimePointer(time.Now().Truncate(24 * time.Hour))
	assert.Equal(t, 1, getNumberOfMissedDeadlines(h))

	h = getActiveHabit("habit2", resources.HBT_REPETITION_DAILY, 7, 3, 3, 0, 7)
	h.Deadline = utils.GetTimePointer(time.Now().Add(-24 * time.Hour).Truncate(24 * time.Hour))
	assert.Equal(t, 2, getNumberOfMissedDeadlines(h))

	h = getActiveHabit("habit3", resources.HBT_REPETITION_WEEKLY, 7, 3, 3, 0, 7)
	h.Deadline = utils.GetTimePointer(time.Now().Add(-24 * time.Hour).Truncate(24 * time.Hour))
	assert.Equal(t, 1, getNumberOfMissedDeadlines(h))

	h = getActiveHabit("habit3", resources.HBT_REPETITION_WEEKLY, 7, 3, 3, 0, 7)
	h.Deadline = utils.GetTimePointer(time.Now().Add(-24 * 7 * time.Hour).Truncate(24 * time.Hour))
	assert.Equal(t, 2, getNumberOfMissedDeadlines(h))

	h = getActiveHabit("habit3", resources.HBT_REPETITION_MONTHLY, 7, 3, 3, 0, 7)
	h.Deadline = utils.GetTimePointer(time.Now().Add(-24 * time.Hour).Truncate(24 * time.Hour))
	assert.Equal(t, 1, getNumberOfMissedDeadlines(h))

	h = getActiveHabit("habit3", resources.HBT_REPETITION_MONTHLY, 7, 3, 3, 0, 7)
	h.Deadline = utils.GetTimePointer(time.Now().Add(-24 * 30 * time.Hour).Truncate(24 * time.Hour))
	assert.Equal(t, 2, getNumberOfMissedDeadlines(h))
}

func TestGetModifyHabitFunc(t *testing.T) {
	// change habit name
	h := getInactiveHabit("testHabit1", 2, 1)
	changeStatus := &resources.Status{}
	tr := &transactionMock{}
	tr.Add(func() error {
		getModifyHabitFunc(h, &resources.Command{Name: "testHabit", BasePoints: -1, HabitRepetitionGoal: -1}, changeStatus)()
		return nil
	})
	tr.Execute()

	verifyHabitState("testHabit", "", "", "testHabit1", false, false, 2, 1, 0, 0, 0, 0, 0, h, changeStatus, t)
	assert.Nil(t, h.LastStreakEnd)

	// deactivate habit
	h = getActiveHabit("testHabit2", resources.HBT_REPETITION_DAILY, 3, 2, 1, 1, 7)
	changeStatus = &resources.Status{}
	tr = &transactionMock{}
	tr.Add(func() error {
		getModifyHabitFunc(h, &resources.Command{ActiveFlag: true, BasePoints: -1, HabitRepetitionGoal: -1}, changeStatus)()
		return nil
	})
	tr.Execute()
	verifyHabitState("testHabit2", "Daily", "", "testHabit2", false, false, 3, 2, 1, 1, 7, 0, 0, h, changeStatus, t)
	assert.Nil(t, h.LastStreakEnd)

	// activate habit
	h = getInactiveHabit("testHabit3", 8, 5)
	changeStatus = &resources.Status{}
	tr = &transactionMock{}
	tr.Add(func() error {
		getModifyHabitFunc(h, &resources.Command{Repetition: resources.HBT_REPETITION_DAILY, Deadline: tomorrowDeadlineStr, ActiveFlag: true, BasePoints: 5, HabitRepetitionGoal: -1}, changeStatus)()
		return nil
	})
	tr.Execute()
	verifyHabitState("testHabit3", resources.HBT_REPETITION_DAILY, tomorrowDeadlineStr, "testHabit3", true, false,
		9, 5, 0, 0, 5, 0, 0, h, changeStatus, t)

	// set habit done
	h = getActiveHabit("testHabit4", resources.HBT_REPETITION_WEEKLY, 32, 14, -2, 4, 9)
	previousDeadline := utils.GetTimePointer(h.Deadline.Add(-24 * 7 * time.Hour))
	h.LastStreakEnd = previousDeadline
	changeStatus = &resources.Status{}
	tr = &transactionMock{}
	tr.Add(func() error {
		getModifyHabitFunc(h, &resources.Command{DoneFlag: true, BasePoints: -1, HabitRepetitionGoal: -1}, changeStatus)()
		return nil
	})
	tr.Execute()
	verifyHabitState("testHabit4", resources.HBT_REPETITION_WEEKLY, tomorrowDeadlineStr, "testHabit4", true, true,
		32, 15, 1, 4, 9, 18, 18, h, changeStatus, t)
	assert.True(t, h.LastStreakEnd == previousDeadline)

	// fail done habit
	h = getActiveHabit("testHabit5", resources.HBT_REPETITION_MONTHLY, 2, 1, 1, 0, 12)
	h.Done = true
	changeStatus = &resources.Status{}
	tr = &transactionMock{}
	tr.Add(func() error {
		getModifyHabitFunc(h, &resources.Command{DoneFlag: true, BasePoints: -1, HabitRepetitionGoal: -1}, changeStatus)()
		return nil
	})
	tr.Execute()
	verifyHabitState("testHabit5", resources.HBT_REPETITION_MONTHLY, tomorrowDeadlineStr, "testHabit5", true, false,
		2, 0, -1, 1, 12, -36, -36, h, changeStatus, t)
	assert.Equal(t, h.Deadline.Format(resources.DATE_FORMAT), h.LastStreakEnd.Format(resources.DATE_FORMAT))

	// set habit done previous period
	h = getActiveHabit("testHabit6", resources.HBT_REPETITION_WEEKLY, 8, 6, -1, 6, 3)
	h.LastStreakEnd = utils.GetTimePointer(h.Deadline.Add(-24 * 7 * time.Hour))
	changeStatus = &resources.Status{}
	tr = &transactionMock{}
	tr.Add(func() error {
		getModifyHabitFunc(h, &resources.Command{DonePrevious: true, BasePoints: -1, HabitRepetitionGoal: -1}, changeStatus)()
		return nil
	})
	tr.Execute()
	verifyHabitState("testHabit6", resources.HBT_REPETITION_WEEKLY, tomorrowDeadlineStr, "testHabit6", true, false,
		8, 7, 7, 6, 3, 195, 0, h, changeStatus, t)

	// set done habit done also for previous period
	h = getActiveHabit("testHabit7", resources.HBT_REPETITION_DAILY, 26, 20, 1, 5, 9)
	h.LastStreakEnd = utils.GetTimePointer(h.Deadline.Add(-24 * time.Hour))
	h.Done = true
	changeStatus = &resources.Status{}
	tr = &transactionMock{}
	tr.Add(func() error {
		getModifyHabitFunc(h, &resources.Command{DonePrevious: true, BasePoints: -1, HabitRepetitionGoal: -1}, changeStatus)()
		return nil
	})
	tr.Execute()
	verifyHabitState("testHabit7", resources.HBT_REPETITION_DAILY, tomorrowDeadlineStr, "testHabit7", true, true,
		26, 21, 7, 5, 9, 774, 432, h, changeStatus, t)
}

func TestGetSyncHabitFunc(t *testing.T) {
	// sync done habit
	h := getActiveHabit("testHabit8", resources.HBT_REPETITION_DAILY, 13, 7, 1, 3, 9)
	changeStatus := &resources.Status{}
	h.Done = true
	h.Deadline = utils.GetTimePointer(time.Now().Truncate(24 * time.Hour))
	tr := &transactionMock{}
	tr.Add(func() error {
		getSyncHabitFunc(changeStatus)(h)()
		return nil
	})
	tr.Execute()
	verifyHabitState("testHabit8", resources.HBT_REPETITION_DAILY, tomorrowDeadlineStr, "testHabit8", true, false,
		14, 7, 1, 3, 9, 0, 0, h, changeStatus, t)

	// sync undone habit
	h = getActiveHabit("testHabit9", resources.HBT_REPETITION_DAILY, 19, 12, 3, 2, 8)
	changeStatus = &resources.Status{}
	todayDeadline := utils.GetTimePointer(time.Now().Truncate(24 * time.Hour))
	h.Deadline = todayDeadline
	tr = &transactionMock{}
	tr.Add(func() error {
		getSyncHabitFunc(changeStatus)(h)()
		return nil
	})
	tr.Execute()
	verifyHabitState("testHabit9", resources.HBT_REPETITION_DAILY, tomorrowDeadlineStr, "testHabit9", true, false,
		20, 12, -1, 3, 8, -8, 0, h, changeStatus, t)
	assert.Equal(t, todayDeadline.Format(resources.DATE_FORMAT), h.LastStreakEnd.Format(resources.DATE_FORMAT))

	h = getActiveHabit("testHabit10", resources.HBT_REPETITION_DAILY, 21, 13, 2, 4, 6)
	changeStatus = &resources.Status{}
	h.Deadline = utils.GetTimePointer(time.Now().Add(time.Duration(-1000000000 * 86400)).Truncate(24 * time.Hour))
	tr = &transactionMock{}
	tr.Add(func() error {
		getSyncHabitFunc(changeStatus)(h)()
		return nil
	})
	tr.Execute()
	verifyHabitState("testHabit10", resources.HBT_REPETITION_DAILY, tomorrowDeadlineStr, "testHabit10", true, false,
		23, 13, -2, -1, 6, -30, 0, h, changeStatus, t)

	h = getActiveHabit("testHabit11", resources.HBT_REPETITION_DAILY, 21, 13, 2, 4, 6)
	changeStatus = &resources.Status{}
	h.Deadline = utils.GetTimePointer(time.Now().Add(time.Duration(-1000000000 * 86400 * 7)).Truncate(24 * time.Hour))
	tr = &transactionMock{}
	tr.Add(func() error {
		getSyncHabitFunc(changeStatus)(h)()
		return nil
	})
	tr.Execute()
	verifyHabitState("testHabit11", resources.HBT_REPETITION_DAILY, tomorrowDeadlineStr, "testHabit11", true, false,
		29, 13, -8, -7, 6, -1224, 0, h, changeStatus, t)
}

func TestGetHabits(t *testing.T) {
	tm := &transactionMock{functionsCalled: []string{}}
	tm.Add(func() error {
		return tm.RetrieveEntities(resources.DB_DEFAULT_HABITS_BUCKET_NAME, false, func(id string) resources.Entity {
			return &resources.Habit{}
		})
	})
	tm.Execute()
	verifyTransactionFlow(t, tm)
}

func TestGetHabit(t *testing.T) {
	tm := &transactionMock{functionsCalled: []string{}}
	tm.Add(func() error {
		return tm.RetrieveEntity(resources.DB_DEFAULT_HABITS_BUCKET_NAME, []byte("id"), &resources.Habit{}, false)
	})
	tm.Execute()
	verifyTransactionFlow(t, tm)
}

func verifyHabitState(expectedName, expectedRepetition, expectedDeadline, expectedId string, expectedActive, expectedDone bool,
	expectedTries, expectedSuccesses, expectedActualStreak, expectedLastStreak, expectedBasePoints, expectedScore,
	expectedTodayScore int, h *resources.Habit, changeStatus *resources.Status, t *testing.T) {
	assert.Equal(t, expectedName, h.Name)
	assert.Equal(t, expectedRepetition, h.Repetition)
	assert.Equal(t, expectedActive, h.Active)
	assert.Equal(t, expectedDone, h.Done)
	if expectedDeadline == "" {
		assert.Nil(t, h.Deadline)
	} else {
		assert.Equal(t, expectedDeadline, h.Deadline.Format(resources.DATE_FORMAT))
	}
	assert.Equal(t, expectedTries, h.Tries)
	assert.Equal(t, expectedSuccesses, h.Successes)
	assert.Equal(t, expectedActualStreak, h.ActualStreak)
	assert.Equal(t, expectedLastStreak, h.LastStreak)
	//assert.Equal(t, expectedLastStreakEnd, h.LastStreakEnd == nil)
	assert.Equal(t, expectedBasePoints, h.BasePoints)
	assert.Equal(t, expectedId, h.Id)
	assert.Equal(t, expectedScore, changeStatus.Score)
	assert.Equal(t, expectedTodayScore, changeStatus.Today)
}

func getInactiveHabit(name string, tries, successes int) *resources.Habit {
	return &resources.Habit{Name: name, Active: false, Done: false, Deadline: nil, Tries: tries, Successes: successes,
		ActualStreak: 0, LastStreak: 0, LastStreakEnd: nil, Repetition: "", BasePoints: 0, Id: name}
}

func getActiveHabit(name, repetition string, tries, successes, actualStreak, lastStreak, basePoints int) *resources.Habit {
	return &resources.Habit{Name: name, Active: true, Done: false, Deadline: utils.GetTimePointer(time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour)),
		Tries: tries, Successes: successes, ActualStreak: actualStreak, LastStreak: lastStreak, LastStreakEnd: nil, Repetition: repetition, BasePoints: basePoints, Id: name}
}
