package cli

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
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

var default_filepath string = path.Join(os.Getenv("HOME"), "tl.csv")

type Action struct {
	ActionType     string
	ToggleComplete bool
	TaskFilepath   string
	UpdateIndex    int
	DeleteIndexes  []int
	DeleteRange    [2]int
	Tasks          []task.Task
}

func Init(taskFilepath string) [][]string {

	// if create file if needed and write headers
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

	return records

}

func ParseAction() *Action {

	taskFilepath := flag.String("f", default_filepath, "alternate task data filepath to ~/tl.csv")
	taskTextAppend := flag.String("a", "", "task text to append")
	taskTextPrepend := flag.String("p", "", "task to prepend")
	updateIndex := flag.Int("u", -1, "task number to update")
	newTaskText := flag.String("t", "", "task update text")
	toggleComplete := flag.Bool("c", false, "toggle task complete status")
	deleteString := flag.String("d", "", "task number to delete")
	verbosePrint := flag.Bool("v", false, "print verbose information")
	usage := flag.Bool("h", false, "usage")

	flag.Parse()
	unparsedArgs := flag.Args()

	cliAction := Action{}
	cliAction.TaskFilepath = *taskFilepath

	if *usage {

		cliAction.ActionType = "help"

	} else if len(*taskTextAppend) > 0 {

		cliAction.ActionType = "append"

		newTaskStrings := append([]string{*taskTextAppend}, unparsedArgs...)
		for _, text := range newTaskStrings {
			task := task.Task{text, false}
			cliAction.Tasks = append(cliAction.Tasks, task)
		}

	} else if len(*taskTextPrepend) > 0 {

		cliAction.ActionType = "prepend"

		newTaskStrings := append([]string{*taskTextPrepend}, unparsedArgs...)
		for _, text := range newTaskStrings {
			task := task.Task{text, false}
			cliAction.Tasks = append(cliAction.Tasks, task)
		}

	} else if *updateIndex != -1 {

		if *updateIndex < -1 {
			log.Fatal("Invalid task #")
		}

		cliAction.ActionType = "update"
		cliAction.UpdateIndex = *updateIndex
		cliAction.Tasks = append(cliAction.Tasks, task.Task{*newTaskText, false})
		cliAction.ToggleComplete = *toggleComplete

	} else if len(*deleteString) > 0 {

		rangeRe := regexp.MustCompile("^[0-9]+\\.\\.[0-9]+$")   // 1..4
		commaDelimRe := regexp.MustCompile("^(\\d+(,\\d+)*)+$") // 1 or 1,2,3 etc

		isRange := rangeRe.MatchString(*deleteString)
		isCommaDelim := commaDelimRe.MatchString(*deleteString)

		if isRange {

			cliAction.ActionType = "delete range"
			deleteRange := strings.Split(*deleteString, "..")
			start, err := strconv.Atoi(deleteRange[0])

			if err != nil {
				log.Fatal(err)
			}

			end, err := strconv.Atoi(deleteRange[1])

			if err != nil {
				log.Fatal(err)
			}

			if end < start {
				cliAction.ActionType = "help"
			} else {

				cliAction.DeleteRange[0] = start
				cliAction.DeleteRange[1] = end

			}

		} else if isCommaDelim {

			cliAction.ActionType = "delete"
			deleteIndexes := strings.Split(*deleteString, ",")
			cliAction.DeleteIndexes = make([]int, len(deleteIndexes))

			for index, s := range deleteIndexes {
				parsedIndex, err := strconv.Atoi(s)
				if err != nil {
					log.Fatal(err)
				}
				cliAction.DeleteIndexes[index] = parsedIndex
			}

		} else {
			cliAction.ActionType = "help"
		}

	} else {

		if *verbosePrint {
			cliAction.ActionType = "printv"
		} else {
			cliAction.ActionType = "print"
		}

	}

	return &cliAction

}

func Run(records [][]string, cliAction *Action) {

	headers := records[0]
	currentTasklist := task.RecordsToTasks(records[1:])

	newTasklist := make([]task.Task, len(currentTasklist))
	switch cliAction.ActionType {

	case "help":
		fmt.Println(USAGE_TEXT)
		return

	case "print":
		PrintTasks(currentTasklist)
		return

	case "printv":
		PrintTasksVerbose(currentTasklist)
		return

	case "append":
		newTasklist = append(currentTasklist, cliAction.Tasks...)

	case "prepend":
		newTasklist = append(cliAction.Tasks, currentTasklist...)

	case "delete":
		newTasklist = task.DeleteTasksByIndex(currentTasklist, cliAction.DeleteIndexes)

	case "delete range":
		newTasklist = task.DeleteTasksByRange(currentTasklist, cliAction.DeleteRange[0], cliAction.DeleteRange[1])

	case "update":
		newTasklist = task.UpdateTask(currentTasklist, cliAction.Tasks[0].Text, cliAction.UpdateIndex, cliAction.ToggleComplete)

	}

	task.WriteTasksToDisk(headers, newTasklist, cliAction.TaskFilepath)

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
