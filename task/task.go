package task

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

var (
	CSV_PARSE_FAILED                = "Failed to parse your CSV file."
	CSV_FILE_NOT_FOUND              = "Failed to parse your CSV file."
	CSV_FILE_APPEND_FAILED          = "Failed to append to your CSV file."
	CSV_FILE_TRUNCATE_FAILED        = "Failed to overwrite your CSV file."
	CSV_FILE_WRITE_FAILED           = "Failed to write to your CSV file."
	TASKFILE_PATH            string = path.Join(os.Getenv("HOME"), "tl.csv")
  HEADER_LINE              string = "Name,Completed\n"
)

type Task struct {
	Text      string
	Completed bool
}

func check(e error, msg ...string) {
	if e != nil {
		if len(msg) > 0 {
			fmt.Println(msg)
		}
		log.Fatalln(e)
	}
}

func recordsToTasks(records [][]string) []Task {

	tasks := make([]Task, len(records))

	for index, record := range records {

		b, err := strconv.ParseBool(record[1])
		check(err)

		tasks[index] = Task{record[0], b}

	}

	return tasks

}

func writeOutTaskfile(tasks []Task, filename string) {

	if err := os.Truncate(filename, 0); err != nil {
		log.Fatalln(CSV_FILE_TRUNCATE_FAILED)
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err, CSV_FILE_APPEND_FAILED)

	records := make([][]string, len(tasks)+1)

	records[0] = make([]string, 2)
	records[0][0] = "Name"
	records[0][1] = "Complete"

	for i, task := range tasks {
		records[i+1] = make([]string, 2)
		records[i+1][0] = task.Text
		records[i+1][1] = strconv.FormatBool(task.Completed)
	}

	w := csv.NewWriter(f)
	if err := w.WriteAll(records); err != nil {
		log.Fatalln(CSV_FILE_WRITE_FAILED)
	}

}

func ParseTaskfile(content *string) []string {

	csvReader := csv.NewReader(strings.NewReader(content))
	records, err := csvReader.ReadAll()
	check(err, CSV_PARSE_FAILED)

	return records

}

func getTasksFromFile(filename string) []Task {

	taskfileBytes, err := ioutil.ReadFile(TASKFILE_PATH)
	check(err)

  taskfileContent := string(taskfileBytes)
  records := ParseTaskfile(taskfileContent)

  return records[1:]

}


func AppendTask(task Task, filename string) {

	record := make([]string, 2)
	record[0] = task.Text
	record[1] = strconv.FormatBool(task.Completed)

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err, CSV_FILE_APPEND_FAILED)

	w := csv.NewWriter(f)
	if err := w.Write(record); err != nil {
		log.Fatalln(CSV_FILE_WRITE_FAILED)
	}

	w.Flush()

	return
}

func UpdateTask(index int, newText string, toggleComplete bool, filename string) {

	if index < 0 {
		return
	}

	tasks := GetTasksFromFile(filename)

	if index >= len(tasks) {
		return
	}

	if len(newText) > 0 {
		tasks[index].Text = newText
	}

	if toggleComplete {
		tasks[index].Completed = !tasks[index].Completed
	}

	writeOutTaskfile(tasks, filename)

	return

}

func DeleteTask(index int, filename string) {

	if index < 0 {
		return
	}

	tasks := GetTasksFromFile(filename)

	if index >= len(tasks) {
		return
	}

	tasks = append(tasks[:index], tasks[index+1:]...)

	writeOutTaskfile(tasks, filename)

}
