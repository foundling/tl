package main

import (
	"fmt"
	"strings"
	"tl/cli"
	"tl/task"
)

const USAGE_TEXT string = `tl usage:
  tl [-aduv]
    -v 
      print tasks in verbose format
    -a <text>
    -u <task #> [-t updated task text] [-c]
      update task by task #
      -t <task text>
      -c
        mark complete
    -d <task #>
      delete task by task #
`

func main() {

	cliAction := cli.ArgsToAction()
	records := cli.initCli(cliAction.TaskFilepath)
	headers := records[0]
	currentTasklist := records[1:]
	newTasklist := make([]task.Task, 0)

	switch cliAction.ActionType {
	case "print":
		cli.PrintTasks(currentTasklist)
	case "printv":
		cli.PrintTasksVerbose(currentTasklist)
	case "add":
		newTasklist := task.AppendTask(cliAction, currentTasklist)
	case "delete":
		newTasklist := task.DeleteTask(cliAction, currentTasklist)
	case "update":
		newTasklist := task.UpdateTask(cliAction, currentTasklist)
	case "help":
		fmt.Println(USAGE_TEXT)
	}

	if len(newTasklist) > 0 {
		WriteTasksToDisk(headers, newTasklist, cliAction.TaskFilepath)
	}

}
