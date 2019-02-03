package cli

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"tl/task"
)

var DEFAULT_FILEPATH string = path.Join(os.Getenv("HOME"), "tl.csv")

type Action struct {
	ActionType     string
	ToggleComplete bool
	TaskFilepath   string
	TaskIndex      int
	Task           task.Task
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

	records := task.ParseTaskfile(string(taskfileBytes))
	task.validateRecords(records)

	return &records

}

func ArgsToAction() *Action {

	taskFilepath := flag.String("f", DEFAULT_FILEPATH, "alternate task data filepath to ~/tl.csv")

	taskText := flag.String("a", "", "task text to add")

	updateIndex := flag.Int("u", -1, "task number to update")
	newTaskText := flag.String("t", "", "task update text")
	toggleComplete := flag.Bool("c", false, "toggle task complete status")

	deleteIndex := flag.Int("d", -1, "task number to delete")
	verbosePrint := flag.Bool("v", false, "print verbose information")

	usage := flag.Bool("h", false, "usage")

	flag.Parse()

	cliAction := Action{}
	cliAction.TaskFilepath = *taskFilepath

	if *usage {
		// usage
		cliAction.ActionType = "help"
	} else if len(*taskText) > 0 {

		// add task
		cliAction.ActionType = "add"
		cliAction.Task.Text = *taskText
		if *toggleComplete {
			cliAction.Task.Completed = true
		}

	} else if *updateIndex != -1 {

		// update task
		if *updateIndex < -1 {
			log.Fatal("Invalid task #")
		}
		cliAction.ActionType = "update"
		cliAction.TaskIndex = *updateIndex
		cliAction.Task.Text = *newTaskText
		cliAction.ToggleComplete = *toggleComplete

	} else if *deleteIndex != -1 {

		// delete task
		if *deleteIndex < -1 {
			log.Fatal("Invalid task #")
		}

		cliAction.ActionType = "delete"
		cliAction.TaskIndex = *deleteIndex

	} else {

		// print tasks
		if *verbosePrint {
			cliAction.ActionType = "printv"
		} else {
			cliAction.ActionType = "print"
		}

	}

	return &cliAction

}

func PrintTasks(tasks []task.Task) {

	fmt.Println("total", len(tasks))
	for i, task := range tasks {
		var icon string
		if task.Completed {
			icon = "\033[12;32m✓\033[0m"
		} else {
			icon = "\033[12;31m✖\033[0m"
		}
		fmt.Printf("%d: %s ", i+1, task.Text)
		fmt.Print(icon)
		fmt.Print("\n")
	}

}

func PrintTasksVerbose(tasks []task.Task) {

	fmt.Println("total", len(tasks))

	if len(tasks) == 0 {
		return
	}

	fmt.Println("\nto do")
	fmt.Println("~~~~~~")
	for i, task := range tasks {
		if !task.Completed {
			var icon string
			if task.Completed {
				icon = "\033[12;32m✓\033[0m"
			} else {
				icon = "\033[12;31m✖\033[0m"
			}
			fmt.Printf("%d: %s ", i+1, task.Text)
			fmt.Print(icon)
			fmt.Println("")
		}
	}

	fmt.Println("")
	fmt.Println("\ncomplete")
	fmt.Println("~~~~~~~~~")
	for i, task := range tasks {
		if task.Completed {
			var icon string
			if task.Completed {
				icon = "\033[12;32m✓\033[0m"
			} else {
				icon = "\033[12;31m✖\033[0m"
			}
			fmt.Printf("%d: %s ", i+1, task.Text)
			fmt.Print(icon)
			fmt.Print("\n")
		}
	}
	fmt.Println("")

}
