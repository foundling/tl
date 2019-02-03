package main

import (
	"fmt"
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
	command := cliAction.ActionType
	records := cli.InitCli(cliAction.TaskFilepath)
	headers := records[0]
	currentTasklist := task.RecordsToTasks(records[1:])
	newTasklist := make([]task.Task, 0)

	switch command {
	case "help":
		fmt.Println(USAGE_TEXT)
	case "print":
		cli.PrintTasks(currentTasklist)
	case "printv":
		cli.PrintTasksVerbose(currentTasklist)
	case "add":
		newTasklist = append(currentTasklist, cliAction.Task)
		task.WriteTasksToDisk(headers, newTasklist, cliAction.TaskFilepath)
	case "delete":
		newTasklist = task.DeleteTask(currentTasklist, cliAction.TaskIndex-1)
		task.WriteTasksToDisk(headers, newTasklist, cliAction.TaskFilepath)
	case "update":
		newTasklist = task.UpdateTask(currentTasklist, cliAction.Task, cliAction.TaskIndex-1, cliAction.ToggleComplete)
		task.WriteTasksToDisk(headers, newTasklist, cliAction.TaskFilepath)
	}

}
