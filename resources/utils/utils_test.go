package utils

import (
	"testing"

	"github.com/7joe7/personalmanager/resources"
)

func TestRemoveTaskFromTasks(t *testing.T) {
	taskToRemove := &resources.Task{Id: "2"}
	tasks := []*resources.Task{{Id: "1"}, taskToRemove, {Id: "3"}}
	tasks = RemoveTaskFromTasks(tasks, taskToRemove)
	if len(tasks) != 2 {
		t.Fatalf("Task %v was not removed. %v", taskToRemove, tasks)
	}
	for i := 0; i < len(tasks); i++ {
		if tasks[i].Id == taskToRemove.Id {
			t.Fatalf("Wrong task was removed. %v", tasks)
		}
	}
}
