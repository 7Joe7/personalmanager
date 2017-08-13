package resources

type Command struct {
	Action              string `json:"action,omitempty"`
	ID                  string `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	ProjectID           string `json:"project_id,omitempty"`
	GoalID              string `json:"goal_id,omitempty"`
	TaskID              string `json:"task_id,omitempty"`
	HabitID             string `json:"habit_id,omitempty"`
	Repetition          string `json:"repetition,omitempty"`
	Deadline            string `json:"deadline,omitempty"`
	Alarm               string `json:"alarm,omitempty"`
	Estimate            string `json:"estimate,omitempty"`
	Scheduled           string `json:"scheduled,omitempty"`
	TaskType            string `json:"task_type,omitempty"`
	Note                string `json:"note,omitempty"`
	NoneAllowed         bool   `json:"none_allowed,omitempty"`
	ActiveFlag          bool   `json:"active_flag,omitempty"`
	DoneFlag            bool   `json:"done_flag,omitempty"`
	DonePrevious        bool   `json:"done_previous,omitempty"`
	UndonePrevious      bool   `json:"undone_previous,omitempty"`
	NegativeFlag        bool   `json:"negative_flag,omitempty"`
	LearnedFlag         bool   `json:"learned_flag,omitempty"`
	BasePoints          int    `json:"base_points,omitempty"`
	HabitRepetitionGoal int    `json:"habit_repetition_goal,omitempty"`
}
