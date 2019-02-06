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
	TaskIndex      int // change to dynamic array of indexes to handle delete range case 
	Task           task.Task
}

func InitCli(taskFilepath string) [][]string {

	// if file doesn't exist, create it, write headers
	if _, err := os.Stat(taskFilepath); os.IsNotExist(err) {

		if f, err := os.OpenFile(taskFilepath, os.O_RDWR|os.O_CREATE, 0755); err != nil {
			log.Fatal(err)
		} else {
			f.WriteString(task.HEADER_LINE)
		}
	}

	taskfileBytes, err := ioutil.ReadFile(taskFilepath)
	if err != nil {
		log.Fatal(err)
	}

	records := task.ParseTaskfile(string(taskfileBytes))
	task.ValidateRecords(records)

	return records[:]

}

func ArgsToAction() *Action {

	taskFilepath := flag.String("f", DEFAULT_FILEPATH, "alternate task data filepath to ~/tl.csv")

	taskTextAppend := flag.String("a", "", "task text to append")
	taskTextPrepend := flag.String("p", "", "task to prepend")

	updateIndex := flag.Int("u", -1, "task number to update")
	newTaskText := flag.String("t", "", "task update text")
	toggleComplete := flag.Bool("c", false, "toggle task complete status")

  //flag becomes string, parse for ints
	deleteString := flag.String("d", "", "task number to delete")

	verbosePrint := flag.Bool("v", false, "print verbose information")

	usage := flag.Bool("h", false, "usage")

	flag.Parse()

	cliAction := Action{}
	cliAction.TaskFilepath = *taskFilepath

	if *usage {

		// usage
		cliAction.ActionType = "help"
	} else if len(*taskTextAppend) > 0 {

		// add to back
		cliAction.ActionType = "append"
		cliAction.Task.Text = *taskTextAppend
		if *toggleComplete {
			cliAction.Task.Completed = true
		}

	} else if len(*taskTextPrepend) > 0 {

		// add to front
		cliAction.ActionType = "prepend"
		cliAction.Task.Text = *taskTextPrepend
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

	} else if len(*deleteString) > 0 {

    rangeRe := regexp.MustCompile("^[0-9]+\\.\\.[0-9]+")
    commaDelimRe := regexp.MustCompile("^([0-9],)+[0-9]$")

    isRange := re.MatchString(*deleteString)
    isCommaDelim := re.MatchString(*deleteString)

    if isRange {
      // 1..4
      cliAction.ActionType = "delete"

    } else if isComaDelim {
      // 1,4,7 -- no trailing comma allowed
      cliAction.ActionType = "delete"

    } else {
      // no match, print help
      cliAction.ActionType = "help"
    }

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
