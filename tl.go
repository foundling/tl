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
    -a <task text to append to list>
    -p <task text to prepend to list>
    -u <task #> [-t updated task text] [-c]
      update task by task #
      -t <task text>
      -c
        mark complete
    -d <task #>,...[<task #>]
      delete task by task #, comma-separated list allowed
    -d <task #>..<task #>
      delete task by task-range

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
	case "append":
		newTasklist = append(currentTasklist, cliAction.Task)
		task.WriteTasksToDisk(headers, newTasklist, cliAction.TaskFilepath)
	case "prepend":
		newTasklist = append([]task.Task{cliAction.Task}, currentTasklist...)
		task.WriteTasksToDisk(headers, newTasklist, cliAction.TaskFilepath)
	case "delete":
		newTasklist = task.DeleteTaskByIndex(currentTasklist, cliAction.DeleteIndex)
		task.WriteTasksToDisk(headers, newTasklist, cliAction.TaskFilepath)
	case "delete comma-delim":
		newTasklist = task.DeleteTasksByIndex(currentTasklist, cliAction.DeleteIndexes)
		task.WriteTasksToDisk(headers, newTasklist, cliAction.TaskFilepath)
	case "delete range":
		newTasklist = task.DeleteTasksByRange(currentTasklist, cliAction.DeleteRange[0], cliAction.DeleteRange[1])
		task.WriteTasksToDisk(headers, newTasklist, cliAction.TaskFilepath)
	case "update":
		newTasklist = task.UpdateTask(currentTasklist, cliAction.Task.Text, cliAction.UpdateIndex, cliAction.ToggleComplete)
		task.WriteTasksToDisk(headers, newTasklist, cliAction.TaskFilepath)
	}

}
