package db

import (
	"os"
	"testing"

	"github.com/7joe7/personalmanager/resources"
	"github.com/boltdb/bolt"
	"github.com/7joe7/personalmanager/test"
)

const (
	TEST_DB_PATH = "test-db.db"
)

func TestTransaction_InitializeBucket(t *testing.T) {
	testOpen(t)
	testBucketInitialization(t, resources.DB_DEFAULT_TASKS_BUCKET_NAME)
	removeDb(t)
}

func TestTransaction_GetValue(t *testing.T) {
	testProject1 := &resources.Project{Name: "testProject1", Note:"Note project"}
	testTask1 := &resources.Task{Name: "test1", Note: "note1", Project: testProject1}
	testOpen(t)
	testBucketInitialization(t, resources.DB_DEFAULT_TASKS_BUCKET_NAME)
	testAddEntity("0", testTask1, resources.DB_DEFAULT_TASKS_BUCKET_NAME, t)
	testGetValue("0", resources.DB_DEFAULT_TASKS_BUCKET_NAME, resources.DB_LAST_ID_KEY, t)
	removeDb(t)
}

func TestTransaction_SetValue(t *testing.T) {
	testProject1 := &resources.Project{Name: "testProject1", Note:"Note project"}
	testTask1 := &resources.Task{Name: "test1", Note: "note1", Project: testProject1}
	testOpen(t)
	testBucketInitialization(t, resources.DB_DEFAULT_TASKS_BUCKET_NAME)
	testAddEntity("0", testTask1, resources.DB_DEFAULT_TASKS_BUCKET_NAME, t)
	tr := newTransaction()
	tr.Add(func () error {
		return tr.SetValue(resources.DB_DEFAULT_TASKS_BUCKET_NAME, resources.DB_LAST_ID_KEY, []byte("7"))
	})
	test.ExpectSuccess(t, tr.execute())
	testGetValue("7", resources.DB_DEFAULT_TASKS_BUCKET_NAME, resources.DB_LAST_ID_KEY, t)
	removeDb(t)
}

func TestTransaction_EnsureEntity(t *testing.T) {
	testOpen(t)
	testBucketInitialization(t, resources.DB_DEFAULT_BASIC_BUCKET_NAME)
	status := &resources.Status{Score:49,Today:12}
	tr := newTransaction()
	tr.Add(func () error {
		return tr.EnsureEntity(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_STATUS_KEY, status)
	})
	test.ExpectSuccess(t, tr.execute())
	statusToVerify := &resources.Status{}
	test.ExpectSuccess(t, retrieveEntity(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_STATUS_KEY, statusToVerify))
	testRetrieveStatus(statusToVerify, status, t)
	status.Score = 24
	status.Today = 1
	test.ExpectSuccess(t, tr.execute())
	status2ToVerify := &resources.Status{}
	test.ExpectSuccess(t, retrieveEntity(resources.DB_DEFAULT_BASIC_BUCKET_NAME, resources.DB_ACTUAL_STATUS_KEY, status2ToVerify))
	testRetrieveStatus(status2ToVerify, statusToVerify, t)
	removeDb(t)
}

func TestTransaction_AddEntity(t *testing.T) {
	testProject1 := &resources.Project{Name: "testProject1", Note:"Note project"}
	testTask1 := &resources.Task{Name: "test1", Note: "note1", Project: testProject1}
	testOpen(t)
	testBucketInitialization(t, resources.DB_DEFAULT_TASKS_BUCKET_NAME)
	testBucketInitialization(t, resources.DB_DEFAULT_PROJECTS_BUCKET_NAME)
	testAddEntity("0", testProject1, resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, t)
	tr := newTransaction()
	tr.Add(func () error {
		return tr.AddEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, testTask1)
	})
	test.ExpectSuccess(t, tr.execute())
	taskToVerify := &resources.Task{}
	tr = newTransaction()
	tr.Add(func () error {
		return tr.RetrieveEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(testTask1.Id), taskToVerify)
	})
	test.ExpectSuccess(t, tr.execute())
	testRetrieveTask(testProject1, testTask1, taskToVerify, nil, t)
	removeDb(t)
}

func TestTransaction_RetrieveEntity(t *testing.T) {
	testProject1 := &resources.Project{Name: "testProject1", Note:"Note project"}
	testTask1 := &resources.Task{Name: "test1", Note: "note1", Project: testProject1}
	testOpen(t)
	testBucketInitialization(t, resources.DB_DEFAULT_TASKS_BUCKET_NAME)
	testBucketInitialization(t, resources.DB_DEFAULT_PROJECTS_BUCKET_NAME)
	testAddEntity("0", testProject1, resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, t)
	testAddEntity("0", testTask1, resources.DB_DEFAULT_TASKS_BUCKET_NAME, t)
	taskToVerify := &resources.Task{}
	tr := newTransaction()
	tr.Add(func () error {
		return tr.RetrieveEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(testTask1.Id), taskToVerify)
	})
	test.ExpectSuccess(t, tr.execute())
	testRetrieveTask(testProject1, testTask1, taskToVerify, nil, t)
	removeDb(t)
}

func TestTransaction_ModifyEntity(t *testing.T) {
	testProject1 := &resources.Project{Name: "testProject1", Note:"Note project"}
	testTask1 := &resources.Task{Name: "test1", Note: "note1", Project: testProject1}
	testOpen(t)
	testBucketInitialization(t, resources.DB_DEFAULT_TASKS_BUCKET_NAME)
	testBucketInitialization(t, resources.DB_DEFAULT_PROJECTS_BUCKET_NAME)
	testAddEntity("0", testProject1, resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, t)
	testAddEntity("0", testTask1, resources.DB_DEFAULT_TASKS_BUCKET_NAME, t)
	taskToModify := &resources.Task{}
	tr := newTransaction()
	tr.Add(func () error {
		return tr.ModifyEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(testTask1.Id), taskToModify, func () {
			taskToModify.Name = "name modified by transaction"
			taskToModify.Note = "note modified by transaction"
			taskToModify.Project.Name = "modified through task by transaction which is invalid"
		})
	})
	test.ExpectSuccess(t, tr.execute())
	taskToVerify := &resources.Task{}
	tr = newTransaction()
	tr.Add(func () error {
		return tr.RetrieveEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(testTask1.Id), taskToVerify)
	})
	test.ExpectSuccess(t, tr.execute())
	testRetrieveTask(testProject1, taskToModify, taskToVerify, nil, t)
	removeDb(t)
}

func TestTransaction_MapEntities(t *testing.T) {
	testProject1 := &resources.Project{Name: "testProject1", Note:"Note project"}
	testTask1 := &resources.Task{Name: "test1", Note: "note1", Project: testProject1}
	testTask2 := &resources.Task{Name: "test2", Note: "note2", Project: nil}
	testOpen(t)
	testBucketInitialization(t, resources.DB_DEFAULT_TASKS_BUCKET_NAME)
	testBucketInitialization(t, resources.DB_DEFAULT_PROJECTS_BUCKET_NAME)
	testAddEntity("0", testProject1, resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, t)
	testAddEntity("0", testTask1, resources.DB_DEFAULT_TASKS_BUCKET_NAME, t)
	testAddEntity("1", testTask2, resources.DB_DEFAULT_TASKS_BUCKET_NAME, t)
	task := &resources.Task{}
	tr := newTransaction()
	tr.Add(func () error {
		return tr.MapEntities(resources.DB_DEFAULT_TASKS_BUCKET_NAME, task, func () {
			task.Name += " name mapped by transaction"
			task.Note += " note mapped by transaction"
			task.Project = testProject1
		})
	})
	test.ExpectSuccess(t, tr.execute())
	taskToVerify := &resources.Task{}
	testRetrieveTask(testProject1, task, taskToVerify, retrieveEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(testTask2.Id), taskToVerify), t)
	removeDb(t)
}

func TestOpen(t *testing.T) {
	testOpen(t)
	removeDb(t)
}

func TestTransaction_FilterEntities(t *testing.T) {
	testProject1 := &resources.Project{Name: "testProject1", Note:"Note project"}
	testTask1 := &resources.Task{Name: "test1", Note: "note1", Project: testProject1}
	testTask2 := &resources.Task{Name: "test2", Note: "note2", Project: nil}
	testOpen(t)
	testBucketInitialization(t, resources.DB_DEFAULT_TASKS_BUCKET_NAME)
	testBucketInitialization(t, resources.DB_DEFAULT_PROJECTS_BUCKET_NAME)
	testAddEntity("0", testProject1, resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, t)
	testAddEntity("0", testTask1, resources.DB_DEFAULT_TASKS_BUCKET_NAME, t)
	testAddEntity("1", testTask2, resources.DB_DEFAULT_TASKS_BUCKET_NAME, t)
	task := &resources.Task{}
	tasks := map[string]*resources.Task{}
	copy := func () {
		copy := &resources.Task{}
		*copy = *task
		tasks[task.Id] = copy
	}
	test.ExpectSuccess(t, filterEntities(resources.DB_DEFAULT_TASKS_BUCKET_NAME, task, func () bool { return task.Project == nil }, copy))
	if len(tasks) != 1 {
		t.Errorf("Expected size of tasks to be 1, it is %d.", len(tasks))
	}
	testRetrieveTask(nil, testTask2, tasks[testTask2.Id], nil, t)
	removeDb(t)
}

func TestTransaction_DeleteEntity(t *testing.T) {
	testProject1 := &resources.Project{Name: "testProject1", Note:"Note project"}
	testTask1 := &resources.Task{Name: "test1", Note: "note1", Project: testProject1}
	testTask2 := &resources.Task{Name: "test2", Note: "note2", Project: nil}
	testOpen(t)
	testBucketInitialization(t, resources.DB_DEFAULT_TASKS_BUCKET_NAME)
	testBucketInitialization(t, resources.DB_DEFAULT_PROJECTS_BUCKET_NAME)
	testAddEntity("0", testProject1, resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, t)
	testAddEntity("0", testTask1, resources.DB_DEFAULT_TASKS_BUCKET_NAME, t)
	testAddEntity("1", testTask2, resources.DB_DEFAULT_TASKS_BUCKET_NAME, t)
	tr := newTransaction()
	tr.Add(func () error {
		return tr.DeleteEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(testTask1.Id))
	})
	test.ExpectSuccess(t, tr.execute())
	if err := retrieveEntity(resources.DB_DEFAULT_TASKS_BUCKET_NAME, []byte(testTask1.Id), &resources.Task{}); err == nil {
		t.Errorf("Expected error, got nil.")
	}
	removeDb(t)
}

func TestTransaction_RetrieveEntities(t *testing.T) {
	testProject1 := &resources.Project{Name: "testProject1", Note:"Note project"}
	testTask1 := &resources.Task{Name: "test1", Note: "note1", Project: testProject1}
	testTask2 := &resources.Task{Name: "test2", Note: "note2", Project: nil}
	testOpen(t)
	testBucketInitialization(t, resources.DB_DEFAULT_TASKS_BUCKET_NAME)
	testBucketInitialization(t, resources.DB_DEFAULT_PROJECTS_BUCKET_NAME)
	testAddEntity("0", testProject1, resources.DB_DEFAULT_PROJECTS_BUCKET_NAME, t)
	testAddEntity("0", testTask1, resources.DB_DEFAULT_TASKS_BUCKET_NAME, t)
	testAddEntity("1", testTask2, resources.DB_DEFAULT_TASKS_BUCKET_NAME, t)
	tasks := map[string]*resources.Task{}
	tr := newTransaction()
	tr.Add(func () error {
		return tr.RetrieveEntities(resources.DB_DEFAULT_TASKS_BUCKET_NAME, func (id string) resources.Entity {
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

func TestTransaction_View(t *testing.T) {
	testOpen(t)
	var addWasCalled bool
	tr := newTransaction()
	tr.Add(func () error {
		addWasCalled = true
		return nil
	})
	test.ExpectSuccess(t, tr.view())
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

func testAddEntity(expectedId string, entity resources.Entity, bucketName []byte, t *testing.T) {
	tr := newTransaction()
	tr.Add(func () error {
		return tr.AddEntity(bucketName, entity)
	})
	test.ExpectSuccess(t, tr.execute())
	if entity.GetId() != expectedId {
		t.Errorf("Expected first id to be '0', it is '%s'.", entity.GetId())
	}
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
		value := tr.GetValue(resources.DB_DEFAULT_TASKS_BUCKET_NAME, lastIdKey)
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
