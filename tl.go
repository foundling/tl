package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

func validateTaskfile(records []string) {

  records := task.ParseTaskfile(fileContent)

  if len(records) < 1 {
    log.Fatal("task file has no records.")
  }

  if records[0][0] != "Name" {
    log.Fatal(`task file not valid. First header field should be "Name".`)
  }

  if records[0][1] != "Completed" {
    log.Fatal(`task file not valid. Second header field should be "Complete".`)
  }

}

func initCli(taskFilepath string) [][]string {

  // if file doesn't exist, create it, write headers
	if _, err := os.Stat(taskFilepath); os.IsNotExist(err) {

		if f, err := os.OpenFile(taskFilepath, os.O_RDWR|os.O_CREATE, 0755); err != nil {
			log.Fatal(err)
		} else {
			f.WriteString(task.HEADER_LINE)
		}
	}

  // filepath to string
	taskfileBytes, err := ioutil.ReadFile(taskFilepath)
	if err != nil {
		log.Fatal(err)
	}

  records := ParseTaskfile(string(taskfileBytes))
  validateTaskfile(records)

  return records

}

func main() {

	var cliAction *cli.Action = cli.ArgsToAction()
  tasks := initCli(cliAction.TaskFilepath)

	switch cliAction.ActionType {
	case "print":
    cli.PrintTasks(tasks)
	case "printv":
		cli.PrintTasksVerbose(tasks)
	case "add":
		task.AppendTask(cliAction.Task, cliAction.TaskFilepath)
	case "delete":
		task.DeleteTask(cliAction.TaskIndex-1, cliAction.TaskFilepath)
	case "update":
		task.UpdateTask(cliAction.TaskIndex-1, cliAction.Task.Text, cliAction.ToggleComplete, cliAction.TaskFilepath)
	case "help":
		fmt.Println(USAGE_TEXT)
	}

}
