package resources

const (
	APP_SUPPORT_FOLDER_PATH = "Library/Application Support/Alfred 3/Workflow Data"
	VENDOR_ID               = "org.erneker"
	APP_NAME                = "personalmanager"

	DB_NAME = "personal-manager.db"

	MSG_CREATE_SUCCESS = "Successfully created %s."
	MSG_DELETE_SUCCESS = "Successfully deleted %s."
	MSG_MODIFY_SUCCESS = "Successfully modified %s."
	MSG_SET_SUCCESS    = "Successfully set %s to %s."
	MSG_EXPORT_SUCCESS = "Successfully exported %s."
	MSG_EXPORT_FAILURE = "Export failed %v."
	MSG_SYNC_SUCCESS   = "Successfully synchronized %s."

	LOG_FILE_NAME        = "personal-manager.log"
	LOG_DAEMON_FILE_NAME = "pema-daemon.log"

	HBT_REPETITION_DAILY   = "Daily"
	HBT_REPETITION_WEEKLY  = "Weekly"
	HBT_REPETITION_MONTHLY = "Monthly"
	HBT_REPETITION_YEARLY  = "Yearly"

	CFG_EXPORT_CONFIG_PATH = "export-config.json"

	ICO_BLACK_ALT    = "./icons/black_alt@2x.png"
	ICO_BLACK        = "./icons/black@2x.png"
	ICO_BLUE         = "./icons/blue@2x.png"
	ICO_CYAN         = "./icons/cyan@2x.png"
	ICO_EXCLAMATION  = "./icons/exclamation@2x.png"
	ICO_GREEN        = "./icons/green@2x.png"
	ICO_ORANGE       = "./icons/orange@2x.png"
	ICO_PURPLE       = "./icons/purple@2x.png"
	ICO_QUESTION     = "./icons/question@2x.png"
	ICO_QUESTION_ALT = "./icons/question_alt@2x.png"
	ICO_RED          = "./icons/red@2x.png"
	ICO_WHITE        = "./icons/white@2x.png"
	ICO_WHITE_ALT    = "./icons/white_alt@2x.png"
	ICO_YELLOW       = "./icons/yellow@2x.png"
	ICO_SPECIAL      = "./icons/special.png"
	ICO_HABIT        = "./icons/habit.jpeg"

	ANY_CMD_QUIT         = "quit"
	ANY_CMD_BLACK_ALT    = "black_alt"
	ANY_CMD_BLACK        = "black"
	ANY_CMD_BLUE         = "blue"
	ANY_CMD_CYAN         = "cyan"
	ANY_CMD_EXCLAMATION  = "exclamation"
	ANY_CMD_GREEN        = "green"
	ANY_CMD_ORANGE       = "orange"
	ANY_CMD_PURPLE       = "purple"
	ANY_CMD_QUESTION     = "question"
	ANY_CMD_QUESTION_ALT = "question_alt"
	ANY_CMD_RED          = "red"
	ANY_CMD_WHITE        = "white"
	ANY_CMD_WHITE_ALT    = "white_alt"
	ANY_CMD_YELLOW       = "yellow"
	ANY_CMD_BROWN        = "#8B4513"

	ANY_PORT_ACTIVE_TASK = 2800
	ANY_PORT_WEEKS_LEFT  = 2801
	ANY_PORTS_RANGE_BASE = 2100

	ANY_SLEEP_TIME = 400

	DATE_HOUR_MINUTE_FORMAT = "2.1.2006 15:04"
	DATE_FORMAT             = "2.1.2006"
	HOUR_MINUTE_FORMAT      = "15:04"

	SUB_FORMAT_ACTIVE_DAILY_HABIT            = "%d/%d, actual %d, points %d"
	SUB_FORMAT_ACTIVE_DAILY_HABIT_WITH_ALARM = "%d/%d, actual %d, alarm %v, points %d"
	SUB_FORMAT_ACTIVE_NOT_DAILY              = "%d/%d, actual %d, deadline %v, points %d"
	SUB_FORMAT_ACTIVE_NOT_DAILY_WITH_ALARM   = "%d/%d, actual %d, alarm %v, deadline %v, points %d"
	SUB_FORMAT_ACTIVE_BAD_HABIT              = "%d/%d, actual %d, today %d/%d, average %.2f, points %d"
	SUB_FORMAT_ACTIVE_BAD_HABIT_NOT_DAILY    = "%d/%d, actual %d, period %d/%d, average %.2f, till %v, points %d"
	SUB_FORMAT_NON_ACTIVE_HABIT              = "%d/%d"
	SUB_FORMAT_TASK                          = "%s %s"
	SUB_FORMAT_PROJECT                       = "%s"
	SUB_FORMAT_TAG                           = ""
	SUB_FORMAT_GOAL                          = "Priority %d, %d/%d"
	NAME_FORMAT_STATUS                       = "Total %d, today %d, yesterday %d."
	NAME_FORMAT_EMPTY                        = "There are no %ss."

	HBT_DONE_BASE_ORDER    = 4000
	HBT_BASE_ORDER_DAILY   = 1000
	HBT_BASE_ORDER_WEEKLY  = 1250
	HBT_BASE_ORDER_MONTHLY = 1500

	ACT_CREATE_TASK                      = "create-task"
	ACT_PRINT_DAY_PLAN                   = "print-day-plan"
	ACT_PRINT_TASKS                      = "print-tasks"
	ACT_PRINT_PERSONAL_TASKS             = "print-personal-tasks"
	ACT_PRINT_PERSONAL_NEXT_TASKS        = "print-next-tasks"
	ACT_PRINT_PERSONAL_UNSCHEDULED_TASKS = "print-unscheduled-tasks"
	ACT_PRINT_GOAL_TASKS                 = "print-goal-tasks"
	ACT_PRINT_PROJECT_GOALS              = "print-project-goals"
	ACT_PRINT_SHOPPING_TASKS             = "print-shopping-tasks"
	ACT_PRINT_WORK_NEXT_TASKS            = "print-work-next-tasks"
	ACT_PRINT_WORK_UNSCHEDULED_TASKS     = "print-work-unscheduled-tasks"
	ACT_PRINT_TASK_NOTE                  = "print-task-note"
	ACT_EXPORT_SHOPPING_TASKS            = "export-shopping-tasks"
	ACT_DELETE_TASK                      = "delete-task"
	ACT_MODIFY_TASK                      = "modify-task"
	ACT_CREATE_PROJECT                   = "create-project"
	ACT_PRINT_PROJECTS                   = "print-projects"
	ACT_PRINT_ACTIVE_PROJECTS            = "print-active-projects"
	ACT_PRINT_INACTIVE_PROJECTS          = "print-inactive-projects"
	ACT_DELETE_PROJECT                   = "delete-project"
	ACT_MODIFY_PROJECT                   = "modify-project"
	ACT_CREATE_TAG                       = "create-tag"
	ACT_PRINT_TAGS                       = "print-tags"
	ACT_DELETE_TAG                       = "delete-tag"
	ACT_MODIFY_TAG                       = "modify-tag"
	ACT_CREATE_GOAL                      = "create-goal"
	ACT_PRINT_GOALS                      = "print-goals"
	ACT_PRINT_ACTIVE_GOALS               = "print-active-goals"
	ACT_PRINT_NON_ACTIVE_GOALS           = "print-non-active-goals"
	ACT_PRINT_INCOMPLETE_GOALS           = "print-incomplete-goals"
	ACT_PRINT_DONE_GOALS                 = "print-done-goals"
	ACT_PRINT_POINTS_STATS               = "print-points-stats"
	ACT_DELETE_GOAL                      = "delete-goal"
	ACT_MODIFY_GOAL                      = "modify-goal"
	ACT_CREATE_HABIT                     = "create-habit"
	ACT_PRINT_HABITS                     = "print-habits"
	ACT_PRINT_HABIT_DESCRIPTION          = "print-habit-description"
	ACT_DELETE_HABIT                     = "delete-habit"
	ACT_MODIFY_HABIT                     = "modify-habit"
	ACT_PRINT_REVIEW                     = "print-review"
	ACT_MODIFY_REVIEW                    = "modify-review"
	ACT_SYNC_ANYBAR_PORTS                = "synchronize-anybar-ports"
	ACT_DEBUG_DATABASE                   = "debug-database"
	ACT_SYNC_WITH_JIRA                   = "sync-with-jira"
	ACT_BACKUP_DATABASE                  = "backup-database"
	ACT_SET_CONFIG_VALUE                 = "set-config-value"
	ACT_SET_EMAIL                        = "set-email"
	ACT_CUSTOM                           = "custom"

	TASK_SCHEDULED_NEXT = "NEXT"
	TASK_NOT_SCHEDULED  = "NONE"

	TASK_TYPE_PERSONAL = "PERSONAL"
	TASK_TYPE_WORK     = "WORK"
	TASK_TYPE_SHOPPING = "SHOPPING"

	STOP_CHARACTER = "\r\n\r\n"
	PORT           = 7007
)

var (
	DB_DEFAULT_BASIC_BUCKET_NAME    = []byte("default")
	DB_DEFAULT_TASKS_BUCKET_NAME    = []byte("default.tasks")
	DB_DEFAULT_TAGS_BUCKET_NAME     = []byte("default.tags")
	DB_DEFAULT_PROJECTS_BUCKET_NAME = []byte("default.projects")
	DB_DEFAULT_HABITS_BUCKET_NAME   = []byte("default.habits")
	DB_DEFAULT_GOALS_BUCKET_NAME    = []byte("default.goals")
	DB_DEFAULT_POINTS_BUCKET_NAME   = []byte("default.points")

	BUCKETS_TO_INTIALIZE = [][]byte{
		DB_DEFAULT_BASIC_BUCKET_NAME,
		DB_DEFAULT_PROJECTS_BUCKET_NAME,
		DB_DEFAULT_TAGS_BUCKET_NAME,
		DB_DEFAULT_TASKS_BUCKET_NAME,
		DB_DEFAULT_GOALS_BUCKET_NAME,
		DB_DEFAULT_HABITS_BUCKET_NAME,
		DB_DEFAULT_POINTS_BUCKET_NAME}

	DB_LAST_ID_KEY            = []byte("last.id")
	DB_LAST_SYNC_KEY          = []byte("last.sync")
	DB_ACTUAL_STATUS_KEY      = []byte("actual.status")
	DB_REVIEW_SETTINGS_KEY    = []byte("review.settings")
	DB_ACTUAL_ACTIVE_TASK_KEY = []byte("actual.active.task")
	DB_ANYBAR_ACTIVE_PORTS    = []byte("anybar.active.ports")
	DB_DEFAULT_EMAIL          = []byte("default.email")
	DB_WEEKS_LEFT             = []byte("weeks.left")

	ACTIONS = []string{
		ACT_CREATE_TASK, ACT_PRINT_PERSONAL_NEXT_TASKS, ACT_DELETE_TASK, ACT_MODIFY_TASK,
		ACT_CREATE_PROJECT, ACT_PRINT_PROJECTS, ACT_DELETE_PROJECT, ACT_MODIFY_PROJECT,
		ACT_CREATE_TAG, ACT_PRINT_TAGS, ACT_DELETE_TAG, ACT_MODIFY_TAG, ACT_CREATE_GOAL,
		ACT_PRINT_GOALS, ACT_DELETE_GOAL, ACT_MODIFY_GOAL, ACT_CREATE_HABIT, ACT_PRINT_HABITS,
		ACT_DELETE_HABIT, ACT_MODIFY_HABIT, ACT_PRINT_REVIEW, ACT_MODIFY_REVIEW}
)
