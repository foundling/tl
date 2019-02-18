package task_test

import (
	"testing"
	"tl/task"
)

func TestTaskDeleteByIndex(t *testing.T) {

	tasks := make([]task.Task, 4)
	tasks[0] = task.Task{"task1", true}
	tasks[1] = task.Task{"task2", false}
	tasks[2] = task.Task{"task3", false}
	tasks[3] = task.Task{"task4", false}

	// check bounds
	tasks = task.DeleteTaskByIndex(tasks, 0)
	if len(tasks) != 4 {
		t.Fatalf("Bounds checked failed.")
	}

	tasks = task.DeleteTaskByIndex(tasks, 5)
	if len(tasks) != 4 {
		t.Fatalf("Bounds checked failed.")
	}

	// check single index deletes
	tasks = task.DeleteTaskByIndex(tasks, 1)
	if len(tasks) != 3 {
		t.Fatalf("Delete failed.")
	}

	tasks = task.DeleteTaskByIndex(tasks, 1)
	if len(tasks) != 2 {
		t.Fatalf("Delete failed.")
	}

	tasks = task.DeleteTaskByIndex(tasks, 1)
	if len(tasks) != 1 {
		t.Fatalf("Delete failed.")
	}

	tasks = task.DeleteTaskByIndex(tasks, 1)
	if len(tasks) != 0 {
		t.Fatalf("Delete failed.")
	}

}

func TestTasksDeleteByIndex(t *testing.T) {

	validIndexes := []int{1, 3, 4}
	invalidIndexes := []int{-3, 0, 100, 5}

	tasks := make([]task.Task, 4)
	tasks[0] = task.Task{"task1", true}
	tasks[1] = task.Task{"task2", false}
	tasks[2] = task.Task{"task3", false}
	tasks[3] = task.Task{"task4", false}

	tasks = task.DeleteTasksByIndex(tasks, invalidIndexes)
	if len(tasks) != 4 {
		t.Fatalf("Deleting invalid indexes failed.")
	}

	// delete all but task 2
	tasks = task.DeleteTasksByIndex(tasks, validIndexes)
	if len(tasks) != 1 {
		t.Fatalf("Delete by valid indexes failed.")
	}

	if tasks[0].Text != "task2" {
		t.Fatalf("Delete by valid indexes failed.")
	}

}

func TestTaskDeleteByRange(t *testing.T) {

	tasks := make([]task.Task, 4)
	tasks[0] = task.Task{"task1", true}
	tasks[1] = task.Task{"task2", false}
	tasks[2] = task.Task{"task3", false}
	tasks[3] = task.Task{"task4", false}

	tasks = task.DeleteTasksByRange(tasks, 1, 2)
	if len(tasks) != 2 {
		t.Fatalf("Delete by range failed.")
	}
	if tasks[0].Text != "task3" {
		t.Fatalf("Delete by range failed.")
	}
	if tasks[1].Text != "task4" {
		t.Fatalf("Delete by range failed.")
	}

}

func TestUpdateTask(t *testing.T) {

	tasks := make([]task.Task, 1)
	tasks[0] = task.Task{"update me", false}

	tasks = task.UpdateTask(tasks, "updated!", 1, true)

	if !tasks[0].Completed {
		t.Fatalf("TestUpdateTask toggle complete status failed")
	}
	if tasks[0].Text != "updated!" {
		t.Fatalf("TestUpdateTask update text failed")
	}

	tasks = task.UpdateTask(tasks, "updated again!", 1, true)

	if tasks[0].Completed {
		t.Fatalf("TestUpdateTask toggle complete status failed")
	}

	if tasks[0].Text != "updated again!" {
		t.Fatalf("TestUpdateTask toggle complete status failed")
	}

}
