package operations

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/7joe7/personalmanager/db"
	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/test"
)

const (
	TEST_DB_PATH = "./test-db.db"

	GET_VALUE_CALLED_FORMAT         = "GetValue%s%s"
	SET_VALUE_CALLED_FORMAT         = "SetValue%s%s%s"
	ENSURE_ENTITY_CALLED_FORMAT     = "EnsureEntity%s%s%v"
	ADD_ENTITY_CALLED_FORMAT        = "AddEntity%s%v"
	DELETE_ENTITY_CALLED_FORMAT     = "DeleteEntity%s%s"
	RETRIEVE_ENTITY_CALLED_FORMAT   = "RetrieveEntity%s%s%v"
	RETRIEVE_ENTITIES_CALLED_FORMAT = "RetrieveEntities%s%v"
	MODIFY_ENTITY_CALLED_FORMAT     = "ModifyEntity%s%s%v%v"
	MAP_ENTITIES_CALLED_FORMAT      = "MapEntities%s%v%v"
	INITIALIZE_BUCKET_CALLED_FORMAT = "InitializeBucket%s"
	EXECUTE_CALLED_FORMAT           = "Execute"
	ADD_CALLED_FORMAT               = "Add"
)

type transactionMock struct {
	functionsCalled []string
	execs           []func() error
}

func (tm *transactionMock) GetValue(bucketName, key []byte) []byte {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(GET_VALUE_CALLED_FORMAT, string(bucketName), string(key)))
	return nil
}

func (tm *transactionMock) SetValue(bucketName, key, value []byte) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(SET_VALUE_CALLED_FORMAT, string(bucketName), string(key), string(value)))
	return nil
}

func (tm *transactionMock) EnsureEntity(bucketName, key []byte, entity interface{}) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(ENSURE_ENTITY_CALLED_FORMAT, string(bucketName), string(key), entity))
	return nil
}

func (tm *transactionMock) AddEntity(bucketName []byte, entity resources.Entity) (string, error) {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(ADD_ENTITY_CALLED_FORMAT, string(bucketName), entity))
	return "", nil
}

func (tm *transactionMock) DeleteEntity(bucketName, id []byte) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(DELETE_ENTITY_CALLED_FORMAT, string(bucketName), string(id)))
	return nil
}

func (tm *transactionMock) RetrieveEntity(bucketName, id []byte, entity interface{}) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(RETRIEVE_ENTITY_CALLED_FORMAT, string(bucketName), string(id), entity))
	return nil
}

func (tm *transactionMock) ModifyEntity(bucketName, key []byte, entity interface{}, modifyFunc func()) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(MODIFY_ENTITY_CALLED_FORMAT, string(bucketName), string(key), entity, modifyFunc))
	return nil
}

func (tm *transactionMock) MapEntities(bucketName []byte, entity interface{}, mapFunc func()) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(MAP_ENTITIES_CALLED_FORMAT, string(bucketName), entity, mapFunc))
	return nil
}

func (tm *transactionMock) InitializeBucket(bucketName []byte) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(INITIALIZE_BUCKET_CALLED_FORMAT, string(bucketName)))
	return nil
}

func (tm *transactionMock) RetrieveEntities(bucketName []byte, getObject func(string) interface{}) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(RETRIEVE_ENTITIES_CALLED_FORMAT, string(bucketName), getObject))
	return nil
}

func (tm *transactionMock) Execute() {
	tm.functionsCalled = append(tm.functionsCalled, EXECUTE_CALLED_FORMAT)
	for i := 0; i < len(tm.execs); i++ {
		tm.execs[i]()
	}
}

func (tm *transactionMock) Add(exec func() error) {
	tm.functionsCalled = append(tm.functionsCalled, ADD_CALLED_FORMAT)
	tm.execs = append(tm.execs, exec)
}

func TestInitializeBuckets(t *testing.T) {
	db.Open(TEST_DB_PATH)

	tm := &transactionMock{functionsCalled: []string{}}
	initializeBuckets(tm, resources.BUCKETS_TO_INTIALIZE)
	tm.Execute()
	verifyTransactionFlow(t, tm)

	for j := 2; j < len(resources.BUCKETS_TO_INTIALIZE)+2; j++ {
		expected := fmt.Sprintf(INITIALIZE_BUCKET_CALLED_FORMAT, string(resources.BUCKETS_TO_INTIALIZE[j-2]))
		test.ExpectString(expected, tm.functionsCalled[j], t)
	}

	removeDb(t)
}

func TestEnsureValues(t *testing.T) {
	db.Open(TEST_DB_PATH)

	tm := &transactionMock{functionsCalled: []string{}}
	ensureValues(tm)
	tm.Execute()
	verifyTransactionFlow(t, tm)

	expected := fmt.Sprintf(ENSURE_ENTITY_CALLED_FORMAT, string(resources.DB_DEFAULT_BASIC_BUCKET_NAME), string(resources.DB_ACTUAL_STATUS_KEY), &resources.Status{})
	test.ExpectString(expected, tm.functionsCalled[2], t)

	removeDb(t)
}

func TestSynchronize(t *testing.T) {
	db.Open(TEST_DB_PATH)

	tm := &transactionMock{functionsCalled: []string{}}
	synchronize(tm)
	tm.Execute()
	verifyTransactionFlow(t, tm)

	expected := fmt.Sprintf(GET_VALUE_CALLED_FORMAT, string(resources.DB_DEFAULT_BASIC_BUCKET_NAME), string(resources.DB_LAST_SYNC_KEY))
	test.ExpectString(expected, tm.functionsCalled[2], t)
	habit := &resources.Habit{}
	changeStatus := &resources.Status{}
	expected = fmt.Sprintf(MAP_ENTITIES_CALLED_FORMAT, string(resources.DB_DEFAULT_HABITS_BUCKET_NAME), habit, getSyncHabitFunc(habit, changeStatus))
	test.ExpectString(expected, tm.functionsCalled[3], t)
	status := &resources.Status{}
	expected = fmt.Sprintf(MODIFY_ENTITY_CALLED_FORMAT, string(resources.DB_DEFAULT_BASIC_BUCKET_NAME), resources.DB_ACTUAL_STATUS_KEY, status, getSyncStatusFunc(status, changeStatus))
	test.ExpectString(expected, tm.functionsCalled[4], t)
	expected = fmt.Sprintf(SET_VALUE_CALLED_FORMAT, string(resources.DB_DEFAULT_BASIC_BUCKET_NAME), string(resources.DB_LAST_SYNC_KEY), time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006"))
	test.ExpectString(expected, tm.functionsCalled[5], t)

	removeDb(t)
}

func verifyTransactionFlow(t *testing.T, tm *transactionMock) {
	if tm.functionsCalled[0] != "Add" {
		t.Errorf("Initialize buckets failed. No add called.")
	}

	if tm.functionsCalled[1] != "Execute" {
		t.Errorf("Initialize buckets failed. No execute called.")
	}
}

func removeDb(t *testing.T) {
	test.ExpectSuccess(t, os.Remove(TEST_DB_PATH))
}
