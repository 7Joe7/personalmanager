package resources

import (
	"encoding/json"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/7joe7/personalmanager/utils"
	"github.com/stretchr/testify/assert"
)

var (
	testId                = "test"
	testDeadlineTime      = utils.GetTimePointer(time.Now().Truncate(24 * time.Hour))
	testDeadlineFormatted = testDeadlineTime.Format(DATE_FORMAT)
	testTask              = &Task{Name: "testing task", Note: "testing task note", Project: &Project{Name: "testing project", Note: "testing project note"}}
	testProject           = &Project{Name: "testing project", Note: "testing project note"}
	testTag               = &Tag{Name: "testing tag"}
	testGoal              = &Goal{Name: "testing goal"}
	testStatus            = &Status{Score: 777, Today: 444}
	testHabitNonActive    = &Habit{Name: "testing non active habit", Active: false, Successes: 5, Tries: 15}
	testHabitActive       = &Habit{Name: "testing active habit", Active: true, Successes: 7, Tries: 18, Repetition: HBT_REPETITION_DAILY, ActualStreak: 3, LastStreak: 1, Deadline: testDeadlineTime, BasePoints: 8}
	testHabitDone         = &Habit{Name: "testing done habit", Active: true, Successes: 2, Tries: 4, Repetition: HBT_REPETITION_WEEKLY, ActualStreak: 1, LastStreak: 0, Deadline: testDeadlineTime, BasePoints: 4, Done: true}

	testTasks                 = Tasks{Tasks: map[string]*Task{"testTask": testTask}}
	expectedEmptyTasksJson    = `{"items":[{"title":"There are no tasks.","valid":false,"icon":{"path":"./icons/black@2x.png"},"mods":{"ctrl":{"valid":false,"subtitle":""},"alt":{"valid":false,"subtitle":""},"cmd":{"valid":false,"subtitle":""},"Fn":{"valid":false,"subtitle":""},"Shift":{"valid":false,"subtitle":""}}}]}`
	expectedTasksJson         = `{"items":[{"title":"testing task","arg":"testTask","subtitle":"testing project; 0; Spent: 0h0m/?","valid":true,"icon":{"path":"./icons/yellow@2x.png"}}]}`
	expectedNoneTasksJson     = `{"items":[{"title":"testing task","arg":"testTask","subtitle":"testing project; 0; Spent: 0h0m/?","valid":true,"icon":{"path":"./icons/yellow@2x.png"}},{"title":"No task","arg":"-","valid":true,"icon":{"path":"./icons/black@2x.png"},"mods":{"ctrl":{"valid":false,"subtitle":""},"alt":{"valid":false,"subtitle":""},"cmd":{"valid":false,"subtitle":""},"Fn":{"valid":false,"subtitle":""},"Shift":{"valid":false,"subtitle":""}}}]}`
	testProjects              = Projects{Projects: map[string]*Project{"testProject": testProject}}
	expectedEmptyProjectsJson = `{"items":[{"title":"There are no projects.","valid":false,"icon":{"path":"./icons/black@2x.png"},"mods":{"ctrl":{"valid":false,"subtitle":""},"alt":{"valid":false,"subtitle":""},"cmd":{"valid":false,"subtitle":""},"Fn":{"valid":false,"subtitle":""},"Shift":{"valid":false,"subtitle":""}}}]}`
	expectedProjectJson       = `{"items":[{"title":"testing project","arg":"testProject","subtitle":"0/0 tasks, 0/0 goals","valid":true,"icon":{"path":"./icons/black@2x.png"}}]}`
	expectedNoneProjectsJson  = `{"items":[{"title":"testing project","arg":"testProject","subtitle":"0/0 tasks, 0/0 goals","valid":true,"icon":{"path":"./icons/black@2x.png"}},{"title":"No project","arg":"-","valid":true,"icon":{"path":"./icons/black@2x.png"},"mods":{"ctrl":{"valid":false,"subtitle":""},"alt":{"valid":false,"subtitle":""},"cmd":{"valid":false,"subtitle":""},"Fn":{"valid":false,"subtitle":""},"Shift":{"valid":false,"subtitle":""}}}]}`
	testTags                  = Tags{Tags: map[string]*Tag{"testTag": testTag}}
	expectedEmptyTagsJson     = `{"items":[{"title":"There are no tags.","valid":false,"icon":{"path":"./icons/black@2x.png"},"mods":{"ctrl":{"valid":false,"subtitle":""},"alt":{"valid":false,"subtitle":""},"cmd":{"valid":false,"subtitle":""},"Fn":{"valid":false,"subtitle":""},"Shift":{"valid":false,"subtitle":""}}}]}`
	expectedTagsJson          = `{"items":[{"title":"testing tag","arg":"testTag","valid":true,"icon":{"path":""}}]}`
	expectedNoneTagsJson      = `{"items":[{"title":"testing tag","arg":"testTag","valid":true,"icon":{"path":""}},{"title":"No tag","arg":"-","valid":true,"icon":{"path":"./icons/black@2x.png"},"mods":{"ctrl":{"valid":false,"subtitle":""},"alt":{"valid":false,"subtitle":""},"cmd":{"valid":false,"subtitle":""},"Fn":{"valid":false,"subtitle":""},"Shift":{"valid":false,"subtitle":""}}}]}`
	testGoals                 = Goals{Goals: map[string]*Goal{"testGoal": testGoal}}
	expectedEmptyGoalsJson    = `{"items":[{"title":"There are no goals.","valid":false,"icon":{"path":"./icons/black@2x.png"},"mods":{"ctrl":{"valid":false,"subtitle":""},"alt":{"valid":false,"subtitle":""},"cmd":{"valid":false,"subtitle":""},"Fn":{"valid":false,"subtitle":""},"Shift":{"valid":false,"subtitle":""}}}]}`
	expectedGoalsJson         = `{"items":[{"title":"testing goal","arg":"testGoal","subtitle":"Priority 0, 0/0","valid":true,"icon":{"path":"./icons/black@2x.png"}}]}`
	expectedNoneGoalsJson     = `{"items":[{"title":"testing goal","arg":"testGoal","subtitle":"Priority 0, 0/0","valid":true,"icon":{"path":"./icons/black@2x.png"}},{"title":"No goal","arg":"-","valid":true,"icon":{"path":"./icons/black@2x.png"},"mods":{"ctrl":{"valid":false,"subtitle":""},"alt":{"valid":false,"subtitle":""},"cmd":{"valid":false,"subtitle":""},"Fn":{"valid":false,"subtitle":""},"Shift":{"valid":false,"subtitle":""}}}]}`

	testHabits              = Habits{Habits: map[string]*Habit{"testHabitActive": testHabitActive, "testHabitDone": testHabitDone, "testHabitNonActive": testHabitNonActive}}
	expectedEmptyHabitsJson = `{"items":[{"title":"There are no habits.","valid":false,"icon":{"path":"./icons/black@2x.png"},"mods":{"ctrl":{"valid":false,"subtitle":""},"alt":{"valid":false,"subtitle":""},"cmd":{"valid":false,"subtitle":""},"Fn":{"valid":false,"subtitle":""},"Shift":{"valid":false,"subtitle":""}}}]}`
	expectedHabitsJson      = fmt.Sprintf(`{"items":[{"title":"testing active habit","arg":"testHabitActive","subtitle":"7/18, actual 3, points 8","valid":true,"icon":{"path":"./icons/red@2x.png"}},{"title":"testing non active habit","arg":"testHabitNonActive","subtitle":"5/15","valid":true,"icon":{"path":"./icons/black@2x.png"}},{"title":"testing done habit","arg":"testHabitDone","subtitle":"2/4, actual 1, deadline %s, points 4","valid":true,"icon":{"path":"./icons/green@2x.png"}}]}`, testDeadlineFormatted)

	expectedNoneHabitsJson = fmt.Sprintf(`{"items":[{"title":"testing active habit","arg":"testHabitActive","subtitle":"7/18, actual 3, points 8","valid":true,"icon":{"path":"./icons/red@2x.png"}},{"title":"testing non active habit","arg":"testHabitNonActive","subtitle":"5/15","valid":true,"icon":{"path":"./icons/black@2x.png"}},{"title":"testing done habit","arg":"testHabitDone","subtitle":"2/4, actual 1, deadline %s, points 4","valid":true,"icon":{"path":"./icons/green@2x.png"}},{"title":"No habit","arg":"-","valid":true,"icon":{"path":"./icons/black@2x.png"},"mods":{"ctrl":{"valid":false,"subtitle":""},"alt":{"valid":false,"subtitle":""},"cmd":{"valid":false,"subtitle":""},"Fn":{"valid":false,"subtitle":""},"Shift":{"valid":false,"subtitle":""}}}]}`, testDeadlineFormatted)
	testHabitsOrdering     = []string{"testHabitActive", "testHabitNonActive", "testHabitDone"}

	testItems             = Items{[]*AlfredItem{(&Review{Repetition: HBT_REPETITION_WEEKLY, Deadline: utils.GetFirstSaturday()}).GetItem()}}
	expectedTestItemsJson = fmt.Sprintf(`{"items":[{"title":"Review repeated Weekly, next: %s.","valid":true,"icon":{"path":"./icons/black@2x.png"}}]}`, utils.GetFirstSaturday().Format("2.1.2006"))
)

func TestOrdering(t *testing.T) {
	tHabits := Habits{Habits: map[string]*Habit{"testHabitActive": testHabitActive, "testHabitDone": testHabitDone, "testHabitNonActive": testHabitNonActive}}
	items := alfredItems{}
	for id, habit := range tHabits.Habits {
		items = append(items, habit.GetAlfredItem(id))
	}
	sort.Sort(items)
	fmt.Println(items)
	for i := 0; i < len(items); i++ {
		if testHabitsOrdering[i] != items[i].Arg {
			t.Errorf("Expected order to be %s, got %s.", testHabitsOrdering[i], items[i].Arg)
		}
	}
}

func TestMarshalTasks(t *testing.T) {
	testMarshalling(Tasks{Tasks: map[string]*Task{}}, expectedEmptyTasksJson, t)
	testMarshalling(Projects{Projects: map[string]*Project{}}, expectedEmptyProjectsJson, t)
	testMarshalling(Tags{Tags: map[string]*Tag{}}, expectedEmptyTagsJson, t)
	testMarshalling(Goals{Goals: map[string]*Goal{}}, expectedEmptyGoalsJson, t)
	testMarshalling(Habits{Habits: map[string]*Habit{}}, expectedEmptyHabitsJson, t)

	testMarshalling(testTasks, expectedTasksJson, t)
	testMarshalling(testProjects, expectedProjectJson, t)
	testMarshalling(testTags, expectedTagsJson, t)
	testMarshalling(testGoals, expectedGoalsJson, t)
	testMarshalling(testHabits, expectedHabitsJson, t)

	testTasks.NoneAllowed = true
	testProjects.NoneAllowed = true
	testTags.NoneAllowed = true
	testGoals.NoneAllowed = true
	testHabits.NoneAllowed = true
	testMarshalling(testTasks, expectedNoneTasksJson, t)
	testMarshalling(testProjects, expectedNoneProjectsJson, t)
	testMarshalling(testTags, expectedNoneTagsJson, t)
	testMarshalling(testGoals, expectedNoneGoalsJson, t)
	testMarshalling(testHabits, expectedNoneHabitsJson, t)
	testMarshalling(testItems, expectedTestItemsJson, t)
}

func TestGetTaskItem(t *testing.T) {
	ai := testTask.GetAlfredItem(testId)
	testCommonAttr(ai, true, testId, testTask.Name, fmt.Sprintf("%s; 0; Spent: 0h0m/?", testTask.Project.Name), "./icons/yellow@2x.png", t)
}

func TestGetProjectItem(t *testing.T) {
	ai := testProject.GetAlfredItem(testId)
	testCommonAttr(ai, true, testId, testProject.Name, fmt.Sprintf(SUB_FORMAT_PROJECT, "0/0 tasks, 0/0 goals"), "./icons/black@2x.png", t)
}

func TestGetTagItem(t *testing.T) {
	ai := testTag.GetAlfredItem(testId)
	testCommonAttr(ai, true, testId, testTag.Name, fmt.Sprintf(SUB_FORMAT_TAG), "", t)
}

func TestGetGoalItem(t *testing.T) {
	ai := testGoal.GetAlfredItem(testId)
	var doneTasksNumber int
	for i := 0; i < len(testGoal.Tasks); i++ {
		if testGoal.Tasks[i].Done {
			doneTasksNumber++
		}
	}
	testCommonAttr(ai, true, testId, testGoal.Name, fmt.Sprintf(SUB_FORMAT_GOAL, testGoal.Priority, doneTasksNumber, len(testGoal.Tasks)), "./icons/black@2x.png", t)
}

func TestGetHabitNonActiveItem(t *testing.T) {
	ai := testHabitNonActive.GetAlfredItem(testId)
	testCommonAttr(ai, true, testId, testHabitNonActive.Name, fmt.Sprintf(SUB_FORMAT_NON_ACTIVE_HABIT,
		testHabitNonActive.Successes, testHabitNonActive.Tries), ICO_BLACK, t)
}

func TestGetHabitActiveItem(t *testing.T) {
	ai := testHabitActive.GetAlfredItem(testId)
	testCommonAttr(ai, true, testId, testHabitActive.Name, fmt.Sprintf(SUB_FORMAT_ACTIVE_DAILY_HABIT,
		testHabitActive.Successes, testHabitActive.Tries,
		testHabitActive.ActualStreak,
		testHabitActive.BasePoints), ICO_RED, t)
}

func TestGetHabitDoneItem(t *testing.T) {
	ai := testHabitDone.GetAlfredItem(testId)
	testCommonAttr(ai, true, testId, testHabitDone.Name, fmt.Sprintf(SUB_FORMAT_ACTIVE_NOT_DAILY,
		testHabitDone.Successes, testHabitDone.Tries, testHabitDone.ActualStreak,
		testHabitDone.Deadline.Format(DATE_FORMAT), testHabitDone.BasePoints), ICO_GREEN, t)
}

func TestGetStatusItem(t *testing.T) {
	ai := testStatus.GetAlfredItem()
	testCommonAttr(ai, false, "", fmt.Sprintf(NAME_FORMAT_STATUS, testStatus.Score, testStatus.Today, testStatus.Yesterday), "", ICO_HABIT, t)
}

func TestGetZeroItem(t *testing.T) {
	ai := getZeroItem(true, false, "entity")
	testCommonAttr(ai, true, "-", "No entity", "", ICO_BLACK, t)
	if ai = getZeroItem(false, false, "entity"); ai != nil {
		t.Errorf("Expected zero item to be nil, it was %v.", ai)
	}
	ai = getZeroItem(false, true, "entity")
	testCommonAttr(ai, false, "", fmt.Sprintf(NAME_FORMAT_EMPTY, "entity"), "", ICO_BLACK, t)
}

func testCommonAttr(ai *AlfredItem, valid bool, testId, expectedName, expectedSubtitle, expectedIconPath string, t *testing.T) {
	assert.Equal(t, testId, ai.Arg)
	assert.Equal(t, expectedName, ai.Name)
	assert.Equal(t, valid, ai.Valid)
	assert.Equal(t, expectedSubtitle, ai.Subtitle)
	assert.Equal(t, expectedIconPath, ai.Icon.Path)
}

func testMarshalling(entity interface{}, expectedJson string, t *testing.T) {
	bytes, err := json.Marshal(entity)
	assert.Nil(t, err)
	assert.Equal(t, expectedJson, string(bytes))
}
