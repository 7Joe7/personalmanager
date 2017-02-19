package utils

import (
	"github.com/7joe7/personalmanager/resources"
	"fmt"
	"os"
)

func removeTaskFromTasks(ts []*resources.Task, t *resources.Task) []*resources.Task {
	for i := 0; i < len(ts); i++ {
		if ts[i].Id == t.Id {
			return append(ts[:i], ts[i+1:]...)
		}
	}
	return ts
}

func removeHabitFromHabits(hs []*resources.Habit, h *resources.Habit) []*resources.Habit {
	for i := 0; i < len(hs); i++ {
		if hs[i].Id == h.Id {
			return append(hs[:i], hs[i+1:]...)
		}
	}
	return hs
}

func removeProjectFromProjects(ps []*resources.Project, p *resources.Project) []*resources.Project {
	for i := 0; i < len(ps); i++ {
		if ps[i].Id == p.Id {
			return append(ps[:i], ps[i+1:]...)
		}
	}
	return ps
}

func removeGoalFromGoals(gs []*resources.Goal, g *resources.Goal) []*resources.Goal {
	for i := 0; i < len(gs); i++ {
		if gs[i].Id == g.Id {
			return append(gs[:i], gs[i+1:]...)
		}
	}
	return gs
}

func removeTagFromTags(ts []*resources.Tag, t *resources.Tag) []*resources.Tag {
	for i := 0; i < len(ts); i++ {
		if ts[i].Id == t.Id {
			return append(ts[:i], ts[i+1:]...)
		}
	}
	return ts
}

func getAppSupportFolderPath() string {
	return fmt.Sprintf("%s/%s/%s.%s", os.Getenv("HOME"), resources.APP_SUPPORT_FOLDER_PATH, resources.VENDOR_ID, resources.APP_NAME)
}
