package operations

import (
	"testing"

	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/test"
	"github.com/7joe7/personalmanager/utils"
	"time"
)

var (
	tomorrowDeadline = utils.GetTimePointer(time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour))
	tomorrowDeadlineStr = tomorrowDeadline.Format(resources.DEADLINE_FORMAT)
)

func TestGetModifyHabitFunc(t *testing.T) {
	// change habit name
	h := getInactiveHabit("testHabit1", 2, 1)
	changeStatus := &resources.Status{}
	getModifyHabitFunc(h, "testHabit", "", "", false, false, false, -1, changeStatus)()
	verifyHabitState("testHabit", "", "", "testHabit1", false, false, 2, 1, 0, 0, 0, 0, 0, h, changeStatus, t)
	test.ExpectBool(true, h.LastStreakEnd == nil, t)

	// deactivate habit
	h = getActiveHabit("testHabit2", resources.HBT_REPETITION_DAILY, 3, 2, 1, 1, 7)
	changeStatus = &resources.Status{}
	getModifyHabitFunc(h, "", "", "", true, false, false, -1, changeStatus)()
	verifyHabitState("testHabit2", "", "", "testHabit2", false, false, 3, 2, 0, 0, 0, 0, 0, h, changeStatus, t)
	test.ExpectBool(true, h.LastStreakEnd == nil, t)

	// activate habit
	h = getInactiveHabit("testHabit3", 8, 5)
	changeStatus = &resources.Status{}
	getModifyHabitFunc(h, "", resources.HBT_REPETITION_DAILY, tomorrowDeadlineStr, true, false, false, 5, changeStatus)()
	verifyHabitState("testHabit3", resources.HBT_REPETITION_DAILY, tomorrowDeadlineStr, "testHabit3", true, false,
		9, 5, 0, 0, 5, 0, 0, h, changeStatus, t)

	// set habit done
	h = getActiveHabit("testHabit4", resources.HBT_REPETITION_WEEKLY, 32, 14, -2, 4, 9)
	previousDeadline := utils.GetTimePointer(h.Deadline.Add(-24 * 7 * time.Hour))
	h.LastStreakEnd = previousDeadline
	changeStatus = &resources.Status{}
	getModifyHabitFunc(h, "", "", "", false, true, false, -1, changeStatus)()
	verifyHabitState("testHabit4", resources.HBT_REPETITION_WEEKLY, tomorrowDeadlineStr, "testHabit4", true, true,
		32, 15, 1, 4, 9, 9, 9, h, changeStatus, t)
	test.ExpectBool(true, h.LastStreakEnd == previousDeadline, t)

	// fail done habit
	h = getActiveHabit("testHabit5", resources.HBT_REPETITION_MONTHLY, 2, 1, 1, 0, 12)
	h.Done = true
	changeStatus = &resources.Status{}
	getModifyHabitFunc(h, "", "", "", false, true, false, -1, changeStatus)()
	verifyHabitState("testHabit5", resources.HBT_REPETITION_MONTHLY, tomorrowDeadlineStr, "testHabit5", true, false,
		2, 1, -1, 1, 12, -12, -12, h, changeStatus, t)
	test.ExpectString(h.Deadline.Format(resources.DEADLINE_FORMAT), h.LastStreakEnd.Format(resources.DEADLINE_FORMAT), t)

	// set habit done previous period
	h = getActiveHabit("testHabit6", resources.HBT_REPETITION_WEEKLY, 8, 6, -1, 6, 3)
	h.LastStreakEnd = utils.GetTimePointer(h.Deadline.Add(-24 * 7 * time.Hour))
	changeStatus = &resources.Status{}
	getModifyHabitFunc(h, "", "", "", false, false, true, -1, changeStatus)()
	verifyHabitState("testHabit6", resources.HBT_REPETITION_WEEKLY, tomorrowDeadlineStr, "testHabit6", true, false,
		8, 7, 7, 6, 3, 24, 0, h, changeStatus, t)

	// set done habit done also for previous period
	h = getActiveHabit("testHabit7", resources.HBT_REPETITION_DAILY, 26, 20, 1, 5, 9)
	h.LastStreakEnd = utils.GetTimePointer(h.Deadline.Add(-24 * time.Hour))
	h.Done = true
	changeStatus = &resources.Status{}
	getModifyHabitFunc(h, "", "", "", false, false, true, -1, changeStatus)()
	verifyHabitState("testHabit7", resources.HBT_REPETITION_DAILY, tomorrowDeadlineStr, "testHabit7", true, true,
		26, 21, 7, 5, 9, 117, 54, h, changeStatus, t)
}

func TestGetSyncHabitFunc(t *testing.T) {
	// sync done habit
	h := getActiveHabit("testHabit8", resources.HBT_REPETITION_DAILY, 13, 7, 1, 3, 9)
	changeStatus := &resources.Status{}
	h.Done = true
	h.Deadline = utils.GetTimePointer(time.Now().Truncate(24 * time.Hour))
	getSyncHabitFunc(h, changeStatus)()
	verifyHabitState("testHabit8", resources.HBT_REPETITION_DAILY, tomorrowDeadlineStr, "testHabit8", true, false,
		14, 7, 1, 3, 9, 0, 0, h, changeStatus, t)

	// sync done habit
	h = getActiveHabit("testHabit9", resources.HBT_REPETITION_DAILY, 19, 12, 3, 2, 8)
	changeStatus = &resources.Status{}
	h.Deadline = utils.GetTimePointer(time.Now().Truncate(24 * time.Hour))
	getSyncHabitFunc(h, changeStatus)()
	verifyHabitState("testHabit9", resources.HBT_REPETITION_DAILY, tomorrowDeadlineStr, "testHabit9", true, false,
		20, 12, -1, 3, 8, -8, 0, h, changeStatus, t)
}

func verifyHabitState(expectedName, expectedRepetition, expectedDeadline, expectedId string, expectedActive, expectedDone bool,
		expectedTries, expectedSuccesses, expectedActualStreak, expectedLastStreak, expectedBasePoints, expectedScore,
		expectedTodayScore int, h *resources.Habit, changeStatus *resources.Status, t *testing.T) {
	test.ExpectString(expectedName, h.Name, t)
	test.ExpectString(expectedRepetition, h.Repetition, t)
	test.ExpectBool(expectedActive, h.Active, t)
	test.ExpectBool(expectedDone, h.Done, t)
	if expectedDeadline == "" {
		test.ExpectBool(true, h.Deadline == nil, t)
	} else {
		test.ExpectString(expectedDeadline, h.Deadline.Format(resources.DEADLINE_FORMAT), t)
	}
	test.ExpectInt(expectedTries, h.Tries, t)
	test.ExpectInt(expectedSuccesses, h.Successes, t)
	test.ExpectInt(expectedActualStreak, h.ActualStreak, t)
	test.ExpectInt(expectedLastStreak, h.LastStreak, t)
	//test.ExpectBool(expectedLastStreakEnd, h.LastStreakEnd == nil, t)
	test.ExpectInt(expectedBasePoints, h.BasePoints, t)
	test.ExpectString(expectedId, h.Id, t)
	test.ExpectInt(expectedScore, changeStatus.Score, t)
	test.ExpectInt(expectedTodayScore, changeStatus.Today, t)
}

func getInactiveHabit(name string, tries, successes int) *resources.Habit {
	return &resources.Habit{Name:name, Active:false, Done:false, Deadline: nil, Tries: tries, Successes: successes,
		ActualStreak: 0, LastStreak: 0, LastStreakEnd: nil, Repetition: "", BasePoints: 0, Id: name}
}

func getActiveHabit(name, repetition string, tries, successes, actualStreak, lastStreak, basePoints int) *resources.Habit {
	return &resources.Habit{Name:name, Active:true, Done:false,Deadline:utils.GetTimePointer(time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour)),
		Tries:tries,Successes:successes,ActualStreak:actualStreak,LastStreak:lastStreak,LastStreakEnd:nil,Repetition:repetition, BasePoints:basePoints,Id:name}
}
