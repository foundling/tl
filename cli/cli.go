package cli

import (
	"flag"
	"fmt"
	"log"
	"tl/task"
)

type Action struct {
	ActionType string
	TaskIndex  int
	Task       task.Task
}

func ArgsToAction() *Action {

	cliAction := Action{}

	// append a new task
	taskText := flag.String("a", "", "task text to add")

	// update existing task
	updateIndex := flag.Int("u", -1, "task number to update")
	newTaskText := flag.String("t", "", "task update text")
	toggleComplete := flag.Bool("c", false, "toggle task complete status")

	// delete existing task
	deleteIndex := flag.Int("d", -1, "task number to delete")

	flag.Parse()

	if len(*taskText) > 0 {

		cliAction.ActionType = "add"
		cliAction.Task.Text = *taskText
		cliAction.Task.Completed = *toggleComplete // not working

	} else if *updateIndex != -1 {

		if *updateIndex < -1 {
			log.Fatal("Invalid task #")
		}
		cliAction.ActionType = "update"
		cliAction.TaskIndex = *updateIndex
		cliAction.Task.Text = *newTaskText
		cliAction.Task.Completed = *toggleComplete

	} else if *deleteIndex != -1 {

		if *deleteIndex < -1 {
			log.Fatal("Invalid task #")
		}

		cliAction.ActionType = "delete"
		cliAction.TaskIndex = *deleteIndex

	} else {

		cliAction.ActionType = "read"

	}

	return &cliAction

}

func PrintTasks(tasks []task.Task) {
	fmt.Printf("total tasks: %d\n", len(tasks))
	for i, task := range tasks {
		var icon string
		if task.Completed {
			icon = "✓"
		} else {
			icon = "✖"
		}
		fmt.Printf("(%d) %s ", i+1, task.Text)
		fmt.Print(icon)
		fmt.Print("\n")
	}
}
