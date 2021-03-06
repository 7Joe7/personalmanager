package validation

import (
	"fmt"
	"testing"
)

func TestVerifications(t *testing.T) {
	testCommonVerifications("Task", verifyTask, t)
	testCommonVerifications("Project", verifyProject, t)
	testCommonVerifications("Tag", verifyTag, t)
	testCommonVerifications("Goal", verifyGoal, t)
	testCommonVerifications("Habit", verifyHabit, t)
}

func testCommonVerifications(entityName string, verify func(string) error, t *testing.T) {
	expectedErrText := fmt.Sprintf("%s name is empty.", entityName)
	if err := verify(""); err == nil || err.Error() != expectedErrText {
		t.Errorf("Expected error with text '%s', got %v.", expectedErrText, err)
	}
	if err := verify("valid name"); err != nil {
		t.Errorf("Expected success, got error (%v).", err)
	}
}
