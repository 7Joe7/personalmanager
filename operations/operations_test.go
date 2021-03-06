package operations

import (
	"fmt"
	"testing"

	"github.com/7joe7/personalmanager/resources"
	"github.com/7joe7/personalmanager/utils"
	"github.com/stretchr/testify/assert"
)

func TestInitializeBuckets(t *testing.T) {
	tm := &transactionMock{functionsCalled: []string{}}
	initializeBuckets(tm, resources.BUCKETS_TO_INTIALIZE)
	tm.Execute()
	verifyTransactionFlow(t, tm)

	for j := 2; j < len(resources.BUCKETS_TO_INTIALIZE)+2; j++ {
		expected := fmt.Sprintf(INITIALIZE_BUCKET_CALLED_FORMAT, string(resources.BUCKETS_TO_INTIALIZE[j-2]))
		assert.Equal(t, expected, tm.functionsCalled[j])
	}
}

func TestEnsureValues(t *testing.T) {
	tm := &transactionMock{functionsCalled: []string{}}
	ensureValues(tm)
	tm.Execute()
	verifyTransactionFlow(t, tm)

	expected := fmt.Sprintf(ENSURE_ENTITY_CALLED_FORMAT, string(resources.DB_DEFAULT_BASIC_BUCKET_NAME), string(resources.DB_REVIEW_SETTINGS_KEY), &resources.Review{Repetition: resources.HBT_REPETITION_WEEKLY, Deadline: utils.GetFirstSaturday()})
	assert.Equal(t, expected, tm.functionsCalled[2])
	expected = fmt.Sprintf(ENSURE_VALUE_CALLED_FORMAT, string(resources.DB_DEFAULT_BASIC_BUCKET_NAME), string(resources.DB_ANYBAR_ACTIVE_PORTS), []byte{})
	assert.Equal(t, expected, tm.functionsCalled[3])
	expected = fmt.Sprintf(ENSURE_ENTITY_CALLED_FORMAT, string(resources.DB_DEFAULT_BASIC_BUCKET_NAME), string(resources.DB_ACTUAL_STATUS_KEY), &resources.Status{})
	assert.Equal(t, expected, tm.functionsCalled[4])
}

func TestSynchronize(t *testing.T) {
	tm := &transactionMock{functionsCalled: []string{}}
	resources.Alr = &AlarmManagerMock{}
	synchronize(tm, false)
	tm.Execute()
	verifyTransactionFlow(t, tm)

	//expected := fmt.Sprintf(GET_VALUE_CALLED_FORMAT, string(resources.DB_DEFAULT_BASIC_BUCKET_NAME), string(resources.DB_LAST_SYNC_KEY))
	//test.ExpectString(expected, tm.functionsCalled[2], t)
	//expected = fmt.Sprintf(MAP_ENTITIES_CALLED_FORMAT, string(resources.DB_DEFAULT_HABITS_BUCKET_NAME), true, "getNewEntity", "mapFunc")
	//test.ExpectString(expected, tm.functionsCalled[3], t)
	//status := &resources.Status{}
	//expected = fmt.Sprintf(MODIFY_ENTITY_CALLED_FORMAT, string(resources.DB_DEFAULT_BASIC_BUCKET_NAME), true, resources.DB_ACTUAL_STATUS_KEY, status, "modifyFunc")
	//test.ExpectString(expected, tm.functionsCalled[8], t)
	//expected = fmt.Sprintf(SET_VALUE_CALLED_FORMAT, string(resources.DB_DEFAULT_BASIC_BUCKET_NAME), string(resources.DB_LAST_SYNC_KEY), time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006"))
	//test.ExpectString(expected, tm.functionsCalled[9], t)
}

func verifyTransactionFlow(t *testing.T, tm *transactionMock) {
	assert.Equal(t, "Add", tm.functionsCalled[0])
	assert.False(t, tm.functionsCalled[1] != "Execute" && tm.functionsCalled[1] != "View")
}
