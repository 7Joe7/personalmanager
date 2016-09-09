package resources

const (
	DB_PATH = "./personal-manager.db"

	MSG_CREATE_SUCCESS = "Successfully created %s with id '%s'."
	MSG_DELETE_SUCCESS = "Successfully deleted %s."
	MSG_MODIFY_SUCCESS = "Successfully modified %s."

	LOG_FILE_PATH = "./personal-manager.log"
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

	DB_LAST_ID_KEY   = []byte("last.id")
	DB_LAST_SYNC_KEY = []byte("last.sync")
)
