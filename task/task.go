package task

import (
	"encoding/csv"
	"fmt"
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

func ValidateRecords(records [][]string) {

	if len(records) < 1 {
		log.Fatal("task file has no headers.")
	}

	if records[0][0] != "Name" {
		log.Fatal(`task file not valid. First header field should be "Name".`)
	}

	if records[0][1] != "Completed" {
		log.Fatal(`task file not valid. Second header field should be "Complete".`)
	}

}

func RecordsToTasks(records [][]string) []Task {

	tasks := make([]Task, len(records))

	for index, record := range records {

		b, err := strconv.ParseBool(record[1])
		check(err)

		tasks[index] = Task{record[0], b}

	}

	return tasks

}

func WriteTasksToDisk(headers []string, tasks []Task, filepath string) {

	if err := os.Truncate(filepath, 0); err != nil {
		log.Fatalln(CSV_FILE_TRUNCATE_FAILED)
	}

	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err, CSV_FILE_APPEND_FAILED)

	records := make([][]string, len(headers)+len(tasks))

	records[0] = make([]string, len(headers))
	records[0][0] = "Name"
	records[0][1] = "Completed"

	for i, task := range tasks {
		records[i+1] = make([]string, len(headers))
		records[i+1][0] = task.Text
		records[i+1][1] = strconv.FormatBool(task.Completed)
	}

	w := csv.NewWriter(f)
	if err := w.WriteAll(records); err != nil {
		log.Fatalln(CSV_FILE_WRITE_FAILED)
	}

}

func AppendTask(task Task, tasks []Task) []Task {

	return append(tasks, task)

}

func DeleteTask(tasks []Task, index int) []Task {

	if index < 0 || index >= len(tasks) {
		return tasks
	}

	return append(tasks[:index], tasks[index+1:]...)

}

func UpdateTask(tasks []Task, task Task, index int, toggleComplete bool) []Task {

	if index < 0 || index >= len(tasks) {
		return tasks
	}

	if len(task.Text) > 0 {
		tasks[index].Text = task.Text
	}

	if toggleComplete {
		tasks[index].Completed = !tasks[index].Completed
	}

	return tasks

}

func ParseTaskfile(content string) [][]string {

	csvReader := csv.NewReader(strings.NewReader(content))
	records, err := csvReader.ReadAll()
	check(err, CSV_PARSE_FAILED)

	return records

}
