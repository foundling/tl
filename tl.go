package main

import (
  "encoding/csv"
  "flag"
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "strconv"
  "strings"
)

const TASK_FILENAME string = "tl.csv"
const CSV_PARSE_MSG string = "There was an error parsing your csv file."

type task struct {
  text string
  completed bool
}

type action struct {
  actionType string
  taskIndex int
  task task
}

func getArgs(cliAction *action) {

  // add new task
  taskText := flag.String("a", "", "task text to add")

  // update existing task
  updateIndex := flag.Int("u", -1, "task number to update")
  newTaskText := flag.String("t", "", "task update text")
  toggleComplete := flag.Bool("c", false, "toggle task complete status")

  // delete existing task
  deleteIndex := flag.Int("d", -1, "task number to delete")

  flag.Parse()

  if len(*taskText) > 0 {

    cliAction.actionType = "add"
    cliAction.task.text = *taskText
    return

  } else if *updateIndex != -1 {

    if *updateIndex < -1 {
      log.Fatal("Invalid task #")
    }
    cliAction.actionType = "update"
    cliAction.taskIndex = *updateIndex
    cliAction.task.text = *newTaskText
    cliAction.task.completed = *toggleComplete
    return

  } else if *deleteIndex != -1 {

    if *deleteIndex < -1 {
      log.Fatal("Invalid task #")
    }

    cliAction.actionType = "delete"
    cliAction.taskIndex = *deleteIndex
    return

  } else {

    cliAction.actionType = "read"

  }

}

func check(e error, msg ...string) {

  if e != nil {

    if len(msg) > 0 {
      fmt.Println(msg)
    }

    log.Fatal(e)
  }

}

func printTasks(tasks []task) {
  fmt.Printf("total tasks: %d\n", len(tasks))
  for i, task := range tasks {
    fmt.Printf("(%d) %s [%t]\n", i + 1, task.text, task.completed)
  }
}

func recordsToTasks(records [][]string) []task {

  tasks := make([]task, len(records))

  for index, record := range records {

    b, err := strconv.ParseBool(record[1])
    check(err)

    tasks[index] = makeTask(record[0], b)

  }

  return tasks

}

func makeTask(title string, completed bool) task {
  return task{ title, completed }
}

func getTasksFromFile(filename string) []task {

  taskFileBytes, err := ioutil.ReadFile(TASK_FILENAME)
  check(err)

  csvReader := csv.NewReader(strings.NewReader(string(taskFileBytes)))
  records, err := csvReader.ReadAll()
  check(err, CSV_PARSE_MSG)

  tasks := recordsToTasks(records)

  return tasks

}

func writeTasksToFile(tasks []task, filename string) {

  if err := os.Truncate(filename, 0); err != nil {
    log.Fatalln("Could not truncate the file")
  }

  f, err := os.OpenFile(filename, os.O_RDWR | os.O_CREATE, 0755)
  check(err, "Could not open file for appending.")

  w := csv.NewWriter(f)

  for _, task := range tasks {
    record := make([]string, 2)
    record[0] = task.text
    record[1] = strconv.FormatBool(task.completed)
    if err := w.Write(record); err != nil {
      log.Fatalln("Could not write to file.")
    }
  }

  w.Flush()

  return
}

func appendTask(task task, filename string) {

  record := make([]string,2)
  record[0] = task.text
  record[1] = strconv.FormatBool(task.completed)

  f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  check(err, "Could not open file for appending.")

  w := csv.NewWriter(f)
  if err := w.Write(record); err != nil {
    log.Fatalln("Could not write to file.")
  }

  w.Flush()

  return
}

func updateTask(index int, update task, filename string) {

  if (index < 0) {
    return
  }

  tasks := getTasksFromFile(filename)

  if index >= len(tasks) {
    return
  }

  if len(update.text) > 0 {
    tasks[index].text = update.text
  }

  tasks[index].completed = !tasks[index].completed


  writeTasksToFile(tasks, filename)

  return


}

func deleteTask(index int, filename string) {


  tasks := getTasksFromFile(filename)
  if index >= len(tasks) {
    return
  }

  tasks = append(tasks[:index], tasks[index + 1:]...)

  if err := os.Truncate(filename, 0); err != nil {
    log.Fatalln("Could not truncate the file")
  }

  f, err := os.OpenFile(filename, os.O_RDWR | os.O_CREATE, 0755)
  check(err, "Could not open file for appending.")

  w := csv.NewWriter(f)
  for _, task := range tasks {
    record := make([]string, 2)
    record[0] = task.text
    record[1] = strconv.FormatBool(task.completed)
    if err := w.Write(record); err != nil {
      log.Fatalln("Could not write to file.")
    }
  }

  w.Flush()

  return
}

func main() {

  var cliAction action
  getArgs(&cliAction)

  if (cliAction.actionType == "read") {
    printTasks(getTasksFromFile(TASK_FILENAME))
    return
  }
  if (cliAction.actionType == "add") {
    appendTask(cliAction.task, TASK_FILENAME)
    return
  }
  if (cliAction.actionType == "delete") {
    deleteTask(cliAction.taskIndex - 1, TASK_FILENAME)
    return
  }
  if (cliAction.actionType == "update") {
    updateTask(cliAction.taskIndex - 1, cliAction.task, TASK_FILENAME)
    return
  }

}
