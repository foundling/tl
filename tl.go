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

func validateTaskfile(fileContent string) {
  r := csv.NewReader(strings.NewReader(fileContent))
  if _, err := r.ReadAll(); err != nil {
    log.Fatal(err)
  }
}

func initCli(taskFilepath string) {

  // if file doesn't exist, create it and write headers
	if _, err := os.Stat(taskFilepath); os.IsNotExist(err) {

		if f, err := os.OpenFile(taskFilepath, os.O_RDWR|os.O_CREATE, 0755); err != nil {
			log.Fatal(err)
		} else {
			f.WriteString("Name,Completed\n")
		}
	}

  taskfileBytes, err := ioutil.ReadFile(taskFilepath)
  if err != nil {
    log.Fatal(err)
  }
  validateTaskfile(string(taskfileBytes))


}

func main() {

	var cliAction *cli.Action = cli.ArgsToAction()
	initCli(cliAction.TaskFilepath)

	switch cliAction.ActionType {
	case "print":
		cli.PrintTasks(task.GetTasksFromFile(cliAction.TaskFilepath))
  case "printv":
		cli.PrintTasksVerbose(task.GetTasksFromFile(cliAction.TaskFilepath))
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
