package resources

const (
	DB_PATH = "./personal-manager.db"

	MSG_CREATE_SUCCESS = "Successfully created %s with id '%s'."
	MSG_DELETE_SUCCESS = "Successfully deleted %s."
	MSG_MODIFY_SUCCESS = "Successfully modified %s."

	LOG_FILE_PATH = "./personal-manager.log"

	HBT_REPETITION_DAILY   = "Daily"
	HBT_REPETITION_WEEKLY  = "Weekly"
	HBT_REPETITION_MONTHLY = "Monthly"
)

var (
	DB_DEFAULT_BASIC_BUCKET_NAME    = []byte("default")
	DB_DEFAULT_TASKS_BUCKET_NAME    = []byte("default.tasks")
	DB_DEFAULT_TAGS_BUCKET_NAME     = []byte("default.tags")
	DB_DEFAULT_PROJECTS_BUCKET_NAME = []byte("default.projects")
	DB_DEFAULT_HABITS_BUCKET_NAME   = []byte("default.habits")
	DB_DEFAULT_GOALS_BUCKET_NAME    = []byte("default.goals")

	BUCKETS_TO_INTIALIZE = [][]byte{
		DB_DEFAULT_BASIC_BUCKET_NAME,
		DB_DEFAULT_PROJECTS_BUCKET_NAME,
		DB_DEFAULT_TAGS_BUCKET_NAME,
		DB_DEFAULT_TASKS_BUCKET_NAME,
		DB_DEFAULT_GOALS_BUCKET_NAME,
		DB_DEFAULT_HABITS_BUCKET_NAME}

	DB_LAST_ID_KEY       = []byte("last.id")
	DB_LAST_SYNC_KEY     = []byte("last.sync")
	DB_ACTUAL_STATUS_KEY = []byte("actual.status")

	ACT_CREATE_TASK    = "create-task"
	ACT_PRINT_TASKS    = "print-tasks"
	ACT_DELETE_TASK    = "delete-task"
	ACT_MODIFY_TASK    = "modify-task"
	ACT_CREATE_PROJECT = "create-project"
	ACT_PRINT_PROJECTS = "print-projects"
	ACT_DELETE_PROJECT = "delete-project"
	ACT_MODIFY_PROJECT = "modify-project"
	ACT_CREATE_TAG     = "create-tag"
	ACT_PRINT_TAGS     = "print-tags"
	ACT_DELETE_TAG     = "delete-tag"
	ACT_MODIFY_TAG     = "modify-tag"
	ACT_CREATE_GOAL    = "create-goal"
	ACT_PRINT_GOALS    = "print-goals"
	ACT_DELETE_GOAL    = "delete-goal"
	ACT_MODIFY_GOAL    = "modify-goal"
	ACT_CREATE_HABIT   = "create-habit"
	ACT_PRINT_HABITS   = "print-habits"
	ACT_DELETE_HABIT   = "delete-habit"
	ACT_MODIFY_HABIT   = "modify-habit"

	ACTIONS = []string{
		ACT_CREATE_TASK, ACT_PRINT_TASKS, ACT_DELETE_TASK, ACT_MODIFY_TASK,
		ACT_CREATE_PROJECT, ACT_PRINT_PROJECTS, ACT_DELETE_PROJECT, ACT_MODIFY_PROJECT,
		ACT_CREATE_TAG, ACT_PRINT_TAGS, ACT_DELETE_TAG, ACT_MODIFY_TAG, ACT_CREATE_GOAL,
		ACT_PRINT_GOALS, ACT_DELETE_GOAL, ACT_MODIFY_GOAL, ACT_CREATE_HABIT, ACT_PRINT_HABITS,
		ACT_DELETE_HABIT, ACT_MODIFY_HABIT}
)
