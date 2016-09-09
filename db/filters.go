package db

import "github.com/7joe7/personalmanager/resources"

func GetActiveHabits() map[string]*resources.Habit {
	return FilterHabits(func (h *resources.Habit) bool {
		return h.Active
	})
}
