package operations

import (
	"fmt"
	"testing"
	"time"

	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/test"
	"github.com/7joe7/personalmanager/utils"
)

const (
	GET_VALUE_CALLED_FORMAT         = "GetValue%s%s"
	SET_VALUE_CALLED_FORMAT         = "SetValue%s%s%s"
	ENSURE_VALUE_CALLED_FORMAT      = "EnsureValue%s%s%v"
	MODIFY_VALUE_CALLED_FORMAT      = "ModifyValue%s%s%v"
	ENSURE_ENTITY_CALLED_FORMAT     = "EnsureEntity%s%s%v"
	ADD_ENTITY_CALLED_FORMAT        = "AddEntity%s%v"
	DELETE_ENTITY_CALLED_FORMAT     = "DeleteEntity%s%s"
	RETRIEVE_ENTITY_CALLED_FORMAT   = "RetrieveEntity%s%s%v"
	RETRIEVE_ENTITIES_CALLED_FORMAT = "RetrieveEntities%s%v"
	MODIFY_ENTITY_CALLED_FORMAT     = "ModifyEntity%s%s%v%v"
	MAP_ENTITIES_CALLED_FORMAT      = "MapEntities%s%v%v"
	INITIALIZE_BUCKET_CALLED_FORMAT = "InitializeBucket%s"
	FILTER_ENTITIES_CALLED_FORMAT   = "FilterEntities%s%v%v%v"
	EXECUTE_CALLED_FORMAT           = "Execute"
	VIEW_CALLED_FORMAT              = "View"
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

func (tm *transactionMock) EnsureValue(bucketName, key, defaultValue []byte) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(ENSURE_VALUE_CALLED_FORMAT, string(bucketName), string(key), string(defaultValue)))
	return nil
}

func (tm *transactionMock) ModifyValue(bucketName, key []byte, modify func ([]byte) []byte) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(MODIFY_VALUE_CALLED_FORMAT, string(bucketName), string(key), modify))
	return nil
}

func (tm *transactionMock) EnsureEntity(bucketName, key []byte, entity resources.Entity) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(ENSURE_ENTITY_CALLED_FORMAT, string(bucketName), string(key), entity))
	return nil
}

func (tm *transactionMock) AddEntity(bucketName []byte, entity resources.Entity) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(ADD_ENTITY_CALLED_FORMAT, string(bucketName), entity))
	return nil
}

func (tm *transactionMock) DeleteEntity(bucketName, id []byte) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(DELETE_ENTITY_CALLED_FORMAT, string(bucketName), string(id)))
	return nil
}

func (tm *transactionMock) RetrieveEntity(bucketName, id []byte, entity resources.Entity) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(RETRIEVE_ENTITY_CALLED_FORMAT, string(bucketName), string(id), entity))
	return nil
}

func (tm *transactionMock) ModifyEntity(bucketName, key []byte, entity resources.Entity, modifyFunc func()) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(MODIFY_ENTITY_CALLED_FORMAT, string(bucketName), string(key), entity, modifyFunc))
	return nil
}

func (tm *transactionMock) MapEntities(bucketName []byte, entity resources.Entity, mapFunc func()) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(MAP_ENTITIES_CALLED_FORMAT, string(bucketName), entity, mapFunc))
	return nil
}

func (tm *transactionMock) InitializeBucket(bucketName []byte) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(INITIALIZE_BUCKET_CALLED_FORMAT, string(bucketName)))
	return nil
}

func (tm *transactionMock) RetrieveEntities(bucketName []byte, getObject func(string) resources.Entity) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(RETRIEVE_ENTITIES_CALLED_FORMAT, string(bucketName), getObject))
	return nil
}

func (tm *transactionMock) FilterEntities(bucketName []byte, addEntity func (), getNewEntity func () resources.Entity, filterFunc func () bool) error {
	tm.functionsCalled = append(tm.functionsCalled, fmt.Sprintf(FILTER_ENTITIES_CALLED_FORMAT, string(bucketName), addEntity, getNewEntity, filterFunc))
	return nil
}

func (tm *transactionMock) Execute() {
	tm.functionsCalled = append(tm.functionsCalled, EXECUTE_CALLED_FORMAT)
	for i := 0; i < len(tm.execs); i++ { tm.execs[i]() }
}

func (tm *transactionMock) View() {
	tm.functionsCalled = append(tm.functionsCalled, VIEW_CALLED_FORMAT)
	for i := 0; i < len(tm.execs); i++ { tm.execs[i]() }
}

func (tm *transactionMock) Add(exec func() error) {
	tm.functionsCalled = append(tm.functionsCalled, ADD_CALLED_FORMAT)
	tm.execs = append(tm.execs, exec)
}

func TestInitializeBuckets(t *testing.T) {
	tm := &transactionMock{functionsCalled: []string{}}
	initializeBuckets(tm, resources.BUCKETS_TO_INTIALIZE)
	tm.Execute()
	verifyTransactionFlow(t, tm)

	for j := 2; j < len(resources.BUCKETS_TO_INTIALIZE)+2; j++ {
		expected := fmt.Sprintf(INITIALIZE_BUCKET_CALLED_FORMAT, string(resources.BUCKETS_TO_INTIALIZE[j-2]))
		test.ExpectString(expected, tm.functionsCalled[j], t)
	}
}

func TestEnsureValues(t *testing.T) {
	tm := &transactionMock{functionsCalled: []string{}}
	ensureValues(tm)
	tm.Execute()
	verifyTransactionFlow(t, tm)

	expected := fmt.Sprintf(ENSURE_ENTITY_CALLED_FORMAT, string(resources.DB_DEFAULT_BASIC_BUCKET_NAME), string(resources.DB_REVIEW_SETTINGS_KEY), &resources.Review{Repetition:resources.HBT_REPETITION_WEEKLY, Deadline:utils.GetFirstSaturday()})
	test.ExpectString(expected, tm.functionsCalled[2], t)
	expected = fmt.Sprintf(ENSURE_VALUE_CALLED_FORMAT, string(resources.DB_DEFAULT_BASIC_BUCKET_NAME), string(resources.DB_ANYBAR_ACTIVE_PORTS), []byte{})
	test.ExpectString(expected, tm.functionsCalled[3], t)
	expected = fmt.Sprintf(ENSURE_ENTITY_CALLED_FORMAT, string(resources.DB_DEFAULT_BASIC_BUCKET_NAME), string(resources.DB_ACTUAL_STATUS_KEY), &resources.Status{})
	test.ExpectString(expected, tm.functionsCalled[4], t)
}

func TestSynchronize(t *testing.T) {
	tm := &transactionMock{functionsCalled: []string{}}
	synchronize(tm)
	tm.Execute()
	verifyTransactionFlow(t, tm)

	expected := fmt.Sprintf(GET_VALUE_CALLED_FORMAT, string(resources.DB_DEFAULT_BASIC_BUCKET_NAME), string(resources.DB_LAST_SYNC_KEY))
	test.ExpectString(expected, tm.functionsCalled[2], t)
	habit := &resources.Habit{}
	changeStatus := &resources.Status{}
	expected = fmt.Sprintf(MAP_ENTITIES_CALLED_FORMAT, string(resources.DB_DEFAULT_HABITS_BUCKET_NAME), habit, getSyncHabitFunc(habit, changeStatus, &transactionMock{}))
	test.ExpectString(expected, tm.functionsCalled[3], t)
	status := &resources.Status{}
	expected = fmt.Sprintf(MODIFY_ENTITY_CALLED_FORMAT, string(resources.DB_DEFAULT_BASIC_BUCKET_NAME), resources.DB_ACTUAL_STATUS_KEY, status, getSyncStatusFunc(status, changeStatus))
	test.ExpectString(expected, tm.functionsCalled[4], t)
	expected = fmt.Sprintf(SET_VALUE_CALLED_FORMAT, string(resources.DB_DEFAULT_BASIC_BUCKET_NAME), string(resources.DB_LAST_SYNC_KEY), time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006"))
	test.ExpectString(expected, tm.functionsCalled[5], t)
}

func verifyTransactionFlow(t *testing.T, tm *transactionMock) {
	test.ExpectString("Add", tm.functionsCalled[0], t)
	test.ExpectBool(false,  tm.functionsCalled[1] != "Execute" && tm.functionsCalled[1] != "View", t)
}
