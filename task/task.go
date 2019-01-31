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

func writeTasksToFile(tasks []Task, filename string) {

	if err := os.Truncate(filename, 0); err != nil {
		log.Fatalln(CSV_FILE_TRUNCATE_FAILED)
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	check(err, CSV_FILE_APPEND_FAILED)

	w := csv.NewWriter(f)

	for _, task := range tasks {
		record := make([]string, 2)
		record[0] = task.Text
		record[1] = strconv.FormatBool(task.Completed)
		if err := w.Write(record); err != nil {
			log.Fatalln(CSV_FILE_WRITE_FAILED)
		}
	}

	w.Flush()

	return
}

func GetTasksFromFile(filename string) []Task {

	taskFileBytes, err := ioutil.ReadFile(TASKFILE_PATH)
	check(err)

	csvReader := csv.NewReader(strings.NewReader(string(taskFileBytes)))
	records, err := csvReader.ReadAll()
	check(err, CSV_PARSE_FAILED)

	tasks := recordsToTasks(records)

	return tasks

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

func UpdateTask(index int, update Task, filename string) {

	if index < 0 {
		return
	}

	tasks := GetTasksFromFile(filename)

	if index >= len(tasks) {
		return
	}

	if len(update.Text) > 0 {
		tasks[index].Text = update.Text
	}

	tasks[index].Completed = !tasks[index].Completed

	writeTasksToFile(tasks, filename)

	return

}

func DeleteTask(index int, filename string) {

	tasks := GetTasksFromFile(filename)
	if index >= len(tasks) {
		return
	}

	tasks = append(tasks[:index], tasks[index+1:]...)

	if err := os.Truncate(filename, 0); err != nil {
		log.Fatalln(CSV_FILE_TRUNCATE_FAILED)
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	check(err, CSV_FILE_APPEND_FAILED)

	w := csv.NewWriter(f)
	for _, task := range tasks {
		record := make([]string, 2)
		record[0] = task.Text
		record[1] = strconv.FormatBool(task.Completed)
		if err := w.Write(record); err != nil {
			log.Fatalln(CSV_FILE_WRITE_FAILED)
		}
	}

	w.Flush()

	return
}
