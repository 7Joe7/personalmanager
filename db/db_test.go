package db

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/7joe7/personalmanager/resources"
	"github.com/boltdb/bolt"
	"github.com/7joe7/personalmanager/test"
)

const (
	TEST_DB_PATH = "test-db.db"
)

var (
	testBasicBucketName = []byte("testBasicBucketName")
	testTasksBucketName = []byte("TestTasksBucket")
	testProjectsBucketName = []byte("TestProjectsBucket")
	testProject1   = &resources.Project{Name: "testProject1", Note:"Note project"}
	testTask1      = &resources.Task{Name: "test1", Note: "note1", Project: testProject1}
	testTask2      = &resources.Task{Name: "test2", Note: "note2", Project: nil}
)

func TestTransaction_InitializeBucket(t *testing.T) {
	testOpen(t)
	testBucketInitialization(t, testTasksBucketName)
	removeDb(t)
}

func TestTransaction_GetValue(t *testing.T) {
	testOpen(t)
	testBucketInitialization(t, testTasksBucketName)
	testAddEntity("0", testTask1, testTasksBucketName, t)
	testGetValue("0", testTasksBucketName, resources.DB_LAST_ID_KEY, t)
	removeDb(t)
}

func TestTransaction_SetValue(t *testing.T) {
	testOpen(t)
	testBucketInitialization(t, testTasksBucketName)
	testAddEntity("0", testTask1, testTasksBucketName, t)
	tr := newTransaction()
	tr.Add(func () error {
		return tr.SetValue(testTasksBucketName, resources.DB_LAST_ID_KEY, []byte("7"))
	})
	test.ExpectSuccess(t, tr.execute())
	testGetValue("7", testTasksBucketName, resources.DB_LAST_ID_KEY, t)
	removeDb(t)
}

func TestTransaction_EnsureEntity(t *testing.T) {
	testOpen(t)
	testBucketInitialization(t, testBasicBucketName)
	status := &resources.Status{Score:49,Today:12}
	tr := newTransaction()
	tr.Add(func () error {
		return tr.EnsureEntity(testBasicBucketName, resources.DB_ACTUAL_STATUS_KEY, status)
	})
	test.ExpectSuccess(t, tr.execute())
	statusToVerify := &resources.Status{}
	test.ExpectSuccess(t, retrieveEntity(testBasicBucketName, resources.DB_ACTUAL_STATUS_KEY, statusToVerify))
	testRetrieveStatus(statusToVerify, status, t)
	status.Score = 24
	status.Today = 1
	test.ExpectSuccess(t, tr.execute())
	status2ToVerify := &resources.Status{}
	test.ExpectSuccess(t, retrieveEntity(testBasicBucketName, resources.DB_ACTUAL_STATUS_KEY, status2ToVerify))
	testRetrieveStatus(status2ToVerify, statusToVerify, t)
	removeDb(t)
}

func TestTransaction_AddEntity(t *testing.T) {
	testOpen(t)
	testBucketInitialization(t, testTasksBucketName)
	testBucketInitialization(t, testProjectsBucketName)
	testAddEntity("0", testProject1, testProjectsBucketName, t)
	tr := newTransaction()
	tr.Add(func () error {
		_, err := tr.AddEntity(testTasksBucketName, testTask1)
		return err
	})
	test.ExpectSuccess(t, tr.execute())
	taskToVerify := &resources.Task{}
	testRetrieveTask(testProject1, testTask1, taskToVerify, retrieveEntity(testTasksBucketName, []byte(testTask1.Id), taskToVerify), t)
	removeDb(t)
}

func TestTransaction_RetrieveEntity(t *testing.T) {
	testOpen(t)
	testBucketInitialization(t, testTasksBucketName)
	testBucketInitialization(t, testProjectsBucketName)
	testAddEntity("0", testTask1, testTasksBucketName, t)
	taskToVerify := &resources.Task{}
	tr := newTransaction()
	tr.Add(func () error {
		return tr.RetrieveEntity(testTasksBucketName, []byte(testTask1.Id), taskToVerify)
	})
	test.ExpectSuccess(t, tr.execute())
	testRetrieveTask(testProject1, testTask1, taskToVerify, nil, t)
	removeDb(t)
}

func TestTransaction_ModifyEntity(t *testing.T) {
	testOpen(t)
	testBucketInitialization(t, testTasksBucketName)
	testBucketInitialization(t, testProjectsBucketName)
	testAddEntity("0", testProject1, testProjectsBucketName, t)
	testAddEntity("0", testTask1, testTasksBucketName, t)
	taskToModify := &resources.Task{}
	tr := newTransaction()
	tr.Add(func () error {
		return tr.ModifyEntity(testTasksBucketName, []byte(testTask1.Id), taskToModify, func () {
			taskToModify.Name = "name modified by transaction"
			taskToModify.Note = "note modified by transaction"
			taskToModify.Project.Name = "modified through task by transaction which is invalid"
		})
	})
	test.ExpectSuccess(t, tr.execute())
	taskToVerify := &resources.Task{}
	testRetrieveTask(testProject1, taskToModify, taskToVerify, retrieveEntity(testTasksBucketName, []byte(testTask1.Id), taskToVerify), t)
	removeDb(t)
}

func TestTransaction_MapEntities(t *testing.T) {
	testOpen(t)
	testBucketInitialization(t, testTasksBucketName)
	testBucketInitialization(t, testProjectsBucketName)
	testAddEntity("0", testProject1, testProjectsBucketName, t)
	testAddEntity("0", testTask1, testTasksBucketName, t)
	testAddEntity("1", testTask2, testTasksBucketName, t)
	task := &resources.Task{}
	tr := newTransaction()
	tr.Add(func () error {
		return tr.MapEntities(testTasksBucketName, task, func () {
			task.Name += " name mapped by transaction"
			task.Note += " note mapped by transaction"
			task.Project = testProject1
		})
	})
	test.ExpectSuccess(t, tr.execute())
	taskToVerify := &resources.Task{}
	testRetrieveTask(testProject1, task, taskToVerify, retrieveEntity(testTasksBucketName, []byte(testTask2.Id), taskToVerify), t)
	removeDb(t)
}

func TestOpen(t *testing.T) {
	testOpen(t)
	removeDb(t)
}

func TestAddEntity(t *testing.T) {
	testOpen(t)
	testBucketInitialization(t, testTasksBucketName)
	testBucketInitialization(t, testProjectsBucketName)
	testAddEntity("0", testProject1, testProjectsBucketName, t)
	testAddEntity("0", testTask1, testTasksBucketName, t)
	task := &resources.Task{}
	err := db.View(func(tx *bolt.Tx) error {
		return json.Unmarshal(tx.Bucket(testTasksBucketName).Get([]byte(testTask1.Id)), task)
	})
	testRetrieveTask(testProject1, testTask1, task, err, t)
	removeDb(t)
}

func TestDeleteEntity(t *testing.T) {
	testOpen(t)
	testBucketInitialization(t, testTasksBucketName)
	testBucketInitialization(t, testProjectsBucketName)
	testAddEntity("0", testTask1, testTasksBucketName, t)
	test.ExpectSuccess(t, deleteEntity([]byte(testTask1.Id), testTasksBucketName))
	err := db.View(func (tx *bolt.Tx) error {
		value := tx.Bucket(testTasksBucketName).Get([]byte(testTask1.Id))
		if value != nil {
			t.Errorf("Expected nil, got '%s'.", string(value))
		}
		return nil
	})
	test.ExpectSuccess(t, err)
	removeDb(t)
}

func TestRetrieveEntity(t *testing.T) {
	testOpen(t)
	testBucketInitialization(t, testTasksBucketName)
	testBucketInitialization(t, testProjectsBucketName)
	testAddEntity("0", testProject1, testProjectsBucketName,  t)
	testAddEntity("0", testTask1, testTasksBucketName, t)
	task := &resources.Task{}
	testRetrieveTask(testProject1, testTask1, task, retrieveEntity(testTasksBucketName, []byte(testTask1.Id), task), t)
	removeDb(t)
}

func TestModifyEntity(t *testing.T) {
	testOpen(t)
	testBucketInitialization(t, testTasksBucketName)
	testBucketInitialization(t, testProjectsBucketName)
	testAddEntity("0", testProject1, testProjectsBucketName, t)
	testAddEntity("0", testTask1, testTasksBucketName, t)
	task := &resources.Task{}
	testRetrieveTask(testProject1, testTask1, task, retrieveEntity(testTasksBucketName, []byte(testProject1.Id), task), t)
	taskToBeModified := &resources.Task{}
	err := modifyEntity(testTasksBucketName, []byte(testTask1.Id), taskToBeModified, func () {
		taskToBeModified.Name = "Completely new name"
		taskToBeModified.Note = "Completely new note"
		taskToBeModified.Project.Name = "Modified project name through task which is invalid"
		taskToBeModified.Project.Note = "Modified project note through task which is invalid"
	})
	test.ExpectSuccess(t, err)
	taskToVerify := &resources.Task{}
	testRetrieveTask(testProject1, taskToBeModified, taskToVerify, retrieveEntity(testTasksBucketName, []byte(testTask1.Id), taskToVerify), t)
	removeDb(t)
}

func TestModifyEntity2(t *testing.T) {
	testOpen(t)
	testBucketInitialization(t, testTasksBucketName)
	testBucketInitialization(t, testProjectsBucketName)
	testAddEntity("0", testProject1, testProjectsBucketName, t)
	testAddEntity("0", testTask1, testTasksBucketName, t)
	task := &resources.Task{}
	testRetrieveTask(testProject1, testTask1, task, retrieveEntity(testTasksBucketName, []byte(testTask1.Id), task), t)
	projectToBeModified := &resources.Project{}
	err := modifyEntity(testProjectsBucketName, []byte(testProject1.Id), projectToBeModified, func () {
		projectToBeModified.Name = "New project name"
		projectToBeModified.Note = "New project note"
	})
	test.ExpectSuccess(t, err)
	taskToVerify := &resources.Task{}
	testRetrieveTask(projectToBeModified, testTask1, taskToVerify, retrieveEntity(testTasksBucketName, []byte(testTask1.Id), taskToVerify), t)
	removeDb(t)
}

func TestRetrieveEntities(t *testing.T) {
	testOpen(t)
	testBucketInitialization(t, testTasksBucketName)
	testBucketInitialization(t, testProjectsBucketName)
	testAddEntity("0", testProject1, testProjectsBucketName, t)
	testAddEntity("0", testTask1, testTasksBucketName, t)
	tasks := map[string]*resources.Task{}
	err := retrieveEntities(testTasksBucketName, func (id string) interface{} {
		tasks[id] = &resources.Task{}
		return tasks[id]
	})
	testRetrieveTask(testProject1, testTask1, tasks[testTask1.Id], err, t)
	removeDb(t)
}

func TestMapEntities(t *testing.T) {
	testOpen(t)
	testBucketInitialization(t, testTasksBucketName)
	testBucketInitialization(t, testProjectsBucketName)
	testAddEntity("0", testProject1, testProjectsBucketName, t)
	testAddEntity("0", testTask1, testTasksBucketName, t)
	task := &resources.Task{}
	err := mapEntities(task, testTasksBucketName, func () {
		task.Name += " task mapped"
		task.Note += " task mapped"
		task.Project.Name = " project modified through task which is invalid"
	})
	test.ExpectSuccess(t, err)
	taskToVerify := &resources.Task{}
	test.ExpectSuccess(t, retrieveEntity(testTasksBucketName, []byte(testTask1.Id), taskToVerify))
	testRetrieveTask(testProject1, task, taskToVerify, nil, t)
	removeDb(t)
}

func TestFilterEntities(t *testing.T) {
	testOpen(t)
	testBucketInitialization(t, testTasksBucketName)
	testBucketInitialization(t, testProjectsBucketName)
	testAddEntity("0", testProject1, testProjectsBucketName, t)
	testAddEntity("0", testTask1, testTasksBucketName, t)
	testAddEntity("1", testTask2, testTasksBucketName, t)
	task := &resources.Task{}
	tasks := map[string]*resources.Task{}
	copy := func () {
		copy := &resources.Task{}
		*copy = *task
		tasks[task.Id] = copy
	}
	test.ExpectSuccess(t, filterEntities(testTasksBucketName, task, func () bool { return task.Project == nil }, copy))
	if len(tasks) != 1 {
		t.Errorf("Expected size of tasks to be 1, it is %d.", len(tasks))
	}
	testRetrieveTask(nil, testTask2, tasks[testTask2.Id], nil, t)
	removeDb(t)
}

func TestTransaction_DeleteEntity(t *testing.T) {
	testOpen(t)
	testBucketInitialization(t, testTasksBucketName)
	testBucketInitialization(t, testProjectsBucketName)
	testAddEntity("0", testProject1, testProjectsBucketName, t)
	testAddEntity("0", testTask1, testTasksBucketName, t)
	testAddEntity("1", testTask2, testTasksBucketName, t)
	tr := newTransaction()
	tr.Add(func () error {
		return tr.DeleteEntity(testTasksBucketName, []byte(testTask1.Id))
	})
	test.ExpectSuccess(t, tr.execute())
	if err := retrieveEntity(testTasksBucketName, []byte(testTask1.Id), &resources.Task{}); err == nil {
		t.Errorf("Expected error, got nil.")
	}
	removeDb(t)
}

func TestTransaction_RetrieveEntities(t *testing.T) {
	testOpen(t)
	testBucketInitialization(t, testTasksBucketName)
	testBucketInitialization(t, testProjectsBucketName)
	testAddEntity("0", testProject1, testProjectsBucketName, t)
	testAddEntity("0", testTask1, testTasksBucketName, t)
	testAddEntity("1", testTask2, testTasksBucketName, t)
	tasks := map[string]*resources.Task{}
	tr := newTransaction()
	tr.Add(func () error {
		return tr.RetrieveEntities(testTasksBucketName, func (id string) interface{} {
			tasks[id] = &resources.Task{}
			return tasks[id]
		})
	})
	test.ExpectSuccess(t, tr.execute())
	testRetrieveTask(testProject1, testTask1, tasks[testTask1.Id], nil, t)
	testRetrieveTask(nil, testTask2, tasks[testTask2.Id], nil, t)
	removeDb(t)
}

func TestTransaction_Add(t *testing.T) {
	tr := newTransaction()
	tr.Add(func () error {
		return nil
	})
	if len(tr.execs) != 1 {
		t.Errorf("Expected 1, got %d.", len(tr.execs))
	}
}

func TestTransaction_Execute(t *testing.T) {
	testOpen(t)
	var addWasCalled bool
	tr := newTransaction()
	tr.Add(func () error {
		addWasCalled = true
		return nil
	})
	test.ExpectSuccess(t, tr.execute())
	if !addWasCalled {
		t.Errorf("Expected true, got %v.", addWasCalled)
	}
	removeDb(t)
}



func testRetrieveTask(expectedProject *resources.Project, expectedEntity, gotEntity *resources.Task, err error, t *testing.T) {
	test.ExpectSuccess(t, err)
	if gotEntity.Name != expectedEntity.Name {
		t.Errorf("Expected saved name to be '%s'. It is '%s'.", expectedEntity.Name, gotEntity.Name)
	}
	if gotEntity.Note != expectedEntity.Note {
		t.Errorf("Expected saved note to be '%s'. It is '%s'.", expectedEntity.Note, gotEntity.Note)
	}
	if gotEntity.Id != expectedEntity.Id {
		t.Errorf("Expected ids to be equal. Expected '%s', got '%s'.", expectedEntity.Id, gotEntity.Id)
	}
	if expectedProject == nil {
		if gotEntity.Project != nil {
			t.Errorf("Expected referenced project to be '%v'. It is '%v'.", expectedProject, gotEntity.Project)
		}
	} else {
		if gotEntity.Project.Id != expectedProject.Id {
			t.Errorf("Expected referenced project to be present. Expected '%s', got '%s'.", expectedProject.Id, gotEntity.Project.Id)
		}
		if gotEntity.Project.Name != expectedProject.Name {
			t.Errorf("Expected referenced project name to be equal. Expected '%s', got '%s'.", expectedProject.Name, gotEntity.Project.Name)
		}
		if gotEntity.Project.Note != expectedProject.Note {
			t.Errorf("Expected referenced project note to be equal. Expected '%s', got '%s'.", expectedProject.Note, gotEntity.Project.Note)
		}
	}
}

func testRetrieveStatus(toVerify, expected *resources.Status, t *testing.T) {
	if toVerify.Score != expected.Score {
		t.Errorf("Expected %d, got %d.", expected.Score, toVerify.Score)
	}
	if toVerify.Today != expected.Today {
		t.Errorf("Expected %d, got %d.", expected.Today, toVerify.Today)
	}
}

func testAddEntity(expectedId string, entity resources.Entity, bucketName []byte, t *testing.T) string {
	id, err := addEntity(entity, bucketName)
	test.ExpectSuccess(t, err)
	if id != expectedId {
		t.Errorf("Expected first id to be '0', it is '%s'.", id)
	}
	return id
}

func testOpen(t *testing.T) {
	test.ExpectSuccess(t, open(TEST_DB_PATH))
	fi, err := os.Stat(TEST_DB_PATH)
	test.ExpectSuccess(t, err)
	if fi.Size() == 0 {
		t.Errorf("Verify database existence - database has 0 size.")
	}
}

func testBucketInitialization(t *testing.T, bucketName []byte) {
	tr := newTransaction()
	tr.Add(func() error {
		return tr.InitializeBucket(bucketName)
	})
	test.ExpectSuccess(t, tr.execute())
	err := db.View(func(tx *bolt.Tx) error {
		if tx.Bucket(bucketName) == nil {
			t.Errorf("Expected created bucket, got nil.")
		}
		return nil
	})
	test.ExpectSuccess(t, err)
}

func testGetValue(expectedValue string, bucketName, lastIdKey []byte, t *testing.T) {
	tr := newTransaction()
	tr.Add(func () error {
		value := tr.GetValue(testTasksBucketName, lastIdKey)
		if string(value) != expectedValue {
			t.Errorf("Expected last id key equal to '%s', got '%s'.", expectedValue, string(value))
		}
		return nil
	})
	test.ExpectSuccess(t, tr.execute())
}

func removeDb(t *testing.T) {
	test.ExpectSuccess(t, os.Remove(TEST_DB_PATH))
}
