package resources

import (
	"testing"
	"time"
	"fmt"
	"sort"
	"encoding/json"

	"github.com/7joe7/personalmanager/utils"
)

var (
	testId = "test"
	testDeadlineTime = utils.GetTimePointer(time.Now().Truncate(24 * time.Hour))
	testDeadlineFormatted = testDeadlineTime.Format(DEADLINE_FORMAT)
	testTask = &Task{Name:"testing task", Note:"testing task note", Project: &Project{Name:"testing project", Note:"testing project note"}}
	testProject = &Project{Name:"testing project", Note:"testing project note"}
	testTag = &Tag{Name:"testing tag"}
	testGoal = &Goal{Name:"testing goal"}
	testStatus = &Status{Score:777,Today:444}
	testHabitNonActive = &Habit{Name:"testing non active habit", Active: false, Successes: 5, Tries: 15}
	testHabitActive = &Habit{Name:"testing active habit", Active: true, Successes: 7, Tries: 18, Repetition: HBT_REPETITION_DAILY, ActualStreak: 3, LastStreak: 1, Deadline: testDeadlineTime, BasePoints: 8}
	testHabitDone = &Habit{Name:"testing done habit", Active: true, Successes: 2, Tries: 4, Repetition: HBT_REPETITION_WEEKLY, ActualStreak: 1, LastStreak: 0, Deadline: testDeadlineTime, BasePoints: 4, Done: true}

	testTasks = Tasks{Tasks:map[string]*Task{"testTask":testTask}}
	expectedEmptyTasksJson = `{"items":[{"title":"There are no tasks.","valid":false,"icon":{"path":"./icons/black@2x.png"},"mods":{"ctrl":{"valid":false,"subtitle":""},"alt":{"valid":false,"subtitle":""},"cmd":{"valid":false,"subtitle":""},"Fn":{"valid":false,"subtitle":""},"Shift":{"valid":false,"subtitle":""}}}]}`
	expectedTasksJson = `{"items":[{"title":"testing task","arg":"testTask","subtitle":"testing project testing task note","valid":true,"icon":{"path":""}}]}`
	expectedNoneTasksJson = `{"items":[{"title":"testing task","arg":"testTask","subtitle":"testing project testing task note","valid":true,"icon":{"path":""}},{"title":"None","arg":"-1","valid":true,"icon":{"path":"./icons/black@2x.png"}}]}`
	testProjects = Projects{Projects:map[string]*Project{"testProject":testProject}}
	expectedEmptyProjectsJson = `{"items":[{"title":"There are no projects.","valid":false,"icon":{"path":"./icons/black@2x.png"},"mods":{"ctrl":{"valid":false,"subtitle":""},"alt":{"valid":false,"subtitle":""},"cmd":{"valid":false,"subtitle":""},"Fn":{"valid":false,"subtitle":""},"Shift":{"valid":false,"subtitle":""}}}]}`
	expectedProjectJson = `{"items":[{"title":"testing project","arg":"testProject","subtitle":"testing project note","valid":true,"icon":{"path":""}}]}`
	expectedNoneProjectsJson = `{"items":[{"title":"testing project","arg":"testProject","subtitle":"testing project note","valid":true,"icon":{"path":""}},{"title":"None","arg":"-1","valid":true,"icon":{"path":"./icons/black@2x.png"}}]}`
	testTags = Tags{Tags:map[string]*Tag{"testTag":testTag}}
	expectedEmptyTagsJson = `{"items":[{"title":"There are no tags.","valid":false,"icon":{"path":"./icons/black@2x.png"},"mods":{"ctrl":{"valid":false,"subtitle":""},"alt":{"valid":false,"subtitle":""},"cmd":{"valid":false,"subtitle":""},"Fn":{"valid":false,"subtitle":""},"Shift":{"valid":false,"subtitle":""}}}]}`
	expectedTagsJson = `{"items":[{"title":"testing tag","arg":"testTag","valid":true,"icon":{"path":""}}]}`
	expectedNoneTagsJson = `{"items":[{"title":"testing tag","arg":"testTag","valid":true,"icon":{"path":""}},{"title":"None","arg":"-1","valid":true,"icon":{"path":"./icons/black@2x.png"}}]}`
	testGoals = Goals{Goals:map[string]*Goal{"testGoal":testGoal}}
	expectedEmptyGoalsJson = `{"items":[{"title":"There are no goals.","valid":false,"icon":{"path":"./icons/black@2x.png"},"mods":{"ctrl":{"valid":false,"subtitle":""},"alt":{"valid":false,"subtitle":""},"cmd":{"valid":false,"subtitle":""},"Fn":{"valid":false,"subtitle":""},"Shift":{"valid":false,"subtitle":""}}}]}`
	expectedGoalsJson = `{"items":[{"title":"testing goal","arg":"testGoal","valid":true,"icon":{"path":""}}]}`
	expectedNoneGoalsJson = `{"items":[{"title":"testing goal","arg":"testGoal","valid":true,"icon":{"path":""}},{"title":"None","arg":"-1","valid":true,"icon":{"path":"./icons/black@2x.png"}}]}`

	testHabits = Habits{Habits:map[string]*Habit{"testHabitActive":testHabitActive, "testHabitDone":testHabitDone, "testHabitNonActive":testHabitNonActive}}
	expectedEmptyHabitsJson = `{"items":[{"title":"There are no habits.","valid":false,"icon":{"path":"./icons/black@2x.png"},"mods":{"ctrl":{"valid":false,"subtitle":""},"alt":{"valid":false,"subtitle":""},"cmd":{"valid":false,"subtitle":""},"Fn":{"valid":false,"subtitle":""},"Shift":{"valid":false,"subtitle":""}}}]}`
	expectedHabitsJson = fmt.Sprintf(`{"items":[{"title":"testing active habit","arg":"testHabitActive","subtitle":"Daily, 7/18, actual 3, %s, base points 8","valid":true,"icon":{"path":"./icons/red@2x.png"}},{"title":"testing non active habit","arg":"testHabitNonActive","subtitle":"5/15","valid":true,"icon":{"path":"./icons/black@2x.png"}},{"title":"testing done habit","arg":"testHabitDone","subtitle":"Weekly, 2/4, actual 1, %s, base points 4","valid":true,"icon":{"path":"./icons/green@2x.png"}}]}`, testDeadlineFormatted, testDeadlineFormatted)
	expectedNoneHabitsJson = fmt.Sprintf(`{"items":[{"title":"testing active habit","arg":"testHabitActive","subtitle":"Daily, 7/18, actual 3, %s, base points 8","valid":true,"icon":{"path":"./icons/red@2x.png"}},{"title":"testing non active habit","arg":"testHabitNonActive","subtitle":"5/15","valid":true,"icon":{"path":"./icons/black@2x.png"}},{"title":"testing done habit","arg":"testHabitDone","subtitle":"Weekly, 2/4, actual 1, %s, base points 4","valid":true,"icon":{"path":"./icons/green@2x.png"}},{"title":"None","arg":"-1","valid":true,"icon":{"path":"./icons/black@2x.png"}}]}`, testDeadlineFormatted, testDeadlineFormatted)
	testHabitsOrdering = map[int]string{0:"testHabitActive",1:"testHabitNonActive", 2:"testHabitDone"}
)

func TestOrdering(t *testing.T) {
	tHabits := Habits{Habits:map[string]*Habit{"testHabitActive":testHabitActive, "testHabitDone":testHabitDone, "testHabitNonActive":testHabitNonActive}}
	items := items{}
	for id, habit := range tHabits.Habits {
		items = append(items, habit.getItem(id))
	}
	sort.Sort(items)
	for i := 0; i < len(items); i++ {
		if testHabitsOrdering[i] != items[i].Arg {
			t.Errorf("Expected order to be %s, got %s.", testHabitsOrdering[i], items[i].Arg)
		}
	}
}

func TestMarshalTasks(t *testing.T) {
	testMarshalling(Tasks{Tasks:map[string]*Task{}}, expectedEmptyTasksJson, t)
	testMarshalling(Projects{Projects:map[string]*Project{}}, expectedEmptyProjectsJson, t)
	testMarshalling(Tags{Tags:map[string]*Tag{}}, expectedEmptyTagsJson, t)
	testMarshalling(Goals{Goals:map[string]*Goal{}}, expectedEmptyGoalsJson, t)
	testMarshalling(Habits{Habits:map[string]*Habit{}}, expectedEmptyHabitsJson, t)

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
}

func TestGetTaskItem(t *testing.T) {
	ai := testTask.getItem(testId)
	testCommonAttr(ai, true, testId, testTask.Name, fmt.Sprintf(SUB_FORMAT_TASK, testTask.Project.Name, testTask.Note), "", t)
}

func TestGetProjectItem(t *testing.T) {
	ai := testProject.getItem(testId)
	testCommonAttr(ai, true, testId, testProject.Name, fmt.Sprintf(SUB_FORMAT_PROJECT, testProject.Note), "", t)
}

func TestGetTagItem(t *testing.T) {
	ai := testTag.getItem(testId)
	testCommonAttr(ai, true, testId, testTag.Name, fmt.Sprintf(SUB_FORMAT_TAG), "", t)
}

func TestGetGoalItem(t *testing.T) {
	ai := testGoal.getItem(testId)
	testCommonAttr(ai, true, testId, testGoal.Name, fmt.Sprintf(SUB_FORMAT_GOAL), "", t)
}

func TestGetHabitNonActiveItem(t *testing.T) {
	ai := testHabitNonActive.getItem(testId)
	testCommonAttr(ai, true, testId, testHabitNonActive.Name, fmt.Sprintf(SUB_FORMAT_NON_ACTIVE_HABIT,
		testHabitNonActive.Successes, testHabitNonActive.Tries), ICO_BLACK, t)
	expectedOrder := HBT_BASE_ORDER
	if ai.order != expectedOrder {
		t.Errorf("Expected order to be %d, it was %d.", expectedOrder, ai.order)
	}
}

func TestGetHabitActiveItem(t *testing.T) {
	ai := testHabitActive.getItem(testId)
	testCommonAttr(ai, true, testId, testHabitActive.Name, fmt.Sprintf(SUB_FORMAT_ACTIVE_HABIT,
		testHabitActive.Repetition, testHabitActive.Successes, testHabitActive.Tries,
		testHabitActive.ActualStreak, testHabitActive.Deadline.Format(DEADLINE_FORMAT),
		testHabitActive.BasePoints), ICO_RED, t)
	expectedOrder := HBT_BASE_ORDER - testHabitActive.BasePoints
	if ai.order != expectedOrder {
		t.Errorf("Expected order to be %d, it was %d.", expectedOrder, ai.order)
	}
}

func TestGetHabitDoneItem(t *testing.T) {
	ai := testHabitDone.getItem(testId)
	testCommonAttr(ai, true, testId, testHabitDone.Name, fmt.Sprintf(SUB_FORMAT_ACTIVE_HABIT,
		testHabitDone.Repetition, testHabitDone.Successes, testHabitDone.Tries, testHabitDone.ActualStreak,
		testHabitDone.Deadline.Format(DEADLINE_FORMAT), testHabitDone.BasePoints), ICO_GREEN, t)
	expectedOrder := HBT_DONE_BASE_ORDER - testHabitDone.BasePoints
	if ai.order != expectedOrder {
		t.Errorf("Expected order to be %d, it was %d.", expectedOrder, ai.order)
	}
}

func TestGetStatusItem(t *testing.T) {
	ai := testStatus.getItem()
	testCommonAttr(ai, false, "", fmt.Sprintf(NAME_FORMAT_STATUS, testStatus.Score, testStatus.Today), "", ICO_BLACK, t)
}

func TestGetZeroItem(t *testing.T) {
	ai := getZeroItem(true, false, "entity")
	testCommonAttr(ai, true, "-1", "None", "", ICO_BLACK, t)
	if ai = getZeroItem(false, false, "entity"); ai != nil {
		t.Errorf("Expected zero item to be nil, it was %v.", ai)
	}
	ai = getZeroItem(false, true, "entity")
	testCommonAttr(ai, false, "", fmt.Sprintf(NAME_FORMAT_EMPTY, "entity"), "", ICO_BLACK, t)
}

func testCommonAttr(ai *AlfredItem, valid bool, testId, expectedName, expectedSubtitle, expectedIconPath string, t *testing.T) {
	if ai.Arg != testId {
		t.Errorf("Expected Arg attribute to be %s, it was %s.", testId, ai.Arg)
	}
	if ai.Name != expectedName {
		t.Errorf("Expected Name attribute to be %s, it was %s.", expectedName, ai.Name)
	}
	if ai.Valid != valid {
		t.Errorf("Expected Valid attribute to be %v, it was %v.", valid, ai.Valid)
	}
	if ai.Subtitle != expectedSubtitle {
		t.Errorf("Expected Subtitle attribute to be %s, it was %s.", expectedSubtitle, ai.Subtitle)
	}
	if ai.Icon.Path != expectedIconPath {
		t.Errorf("Expected Icon attribute to be %s, it was %s.", expectedIconPath, ai.Icon.Path)
	}
}

func testMarshalling(entity interface{}, expectedJson string, t *testing.T) {
	bytes, err := json.Marshal(entity)
	if err != nil {
		t.Errorf("Expected success, got error %v.", err)
	}
	actualJson := string(bytes)
	if actualJson != expectedJson {
		t.Errorf("Expected '%s', got '%s'.", expectedJson, actualJson)
	}
}