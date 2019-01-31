package main

import (
	"fmt"
	"os"
	"path"
	"tl/cli"
	"tl/task"
)

const USAGE_TEXT string = `tl usage:
    add:
      -a <task text>

    update:
      -u <task #> [-t updated task text] [-c]

    delete:
      -d <task #>

    sort:
      -s
`

var (
  HOMEDIR = os.Getenv("HOME")
	TASKFILE_PATH string = path.Join(HOMEDIR, "tl.csv")
)

func main() {

	var cliAction *cli.Action = cli.ArgsToAction()

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
		fmt.Println(USAGE_TEXT)
	}

}
