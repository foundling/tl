package main

import (
	"fmt"
	"os"
	"path"
	"tl/cli"
	"tl/task"
)

var (
	HOMEDIR              = os.Getenv("HOME")
	TASKFILE_PATH string = path.Join(HOMEDIR, "tl.csv")
	CSV_PARSE_MSG string = "There was an error parsing your csv file."
	HELP_TEXT            = "tl [options]"
)

func main() {

	cliAction := cli.ArgsToAction()

	switch cliAction.ActionType {
	case "read":
		cli.PrintTasks(task.GetTasksFromFile(TASKFILE_PATH))
	case "add":
		task.AppendTask(cliAction.Task, TASKFILE_PATH)
	case "delete":
		task.DeleteTask(cliAction.TaskIndex-1, TASKFILE_PATH)
	case "update":
		task.UpdateTask(cliAction.TaskIndex-1, cliAction.Task, TASKFILE_PATH)
	case "help":
		fmt.Println(HELP_TEXT)
	}

}
