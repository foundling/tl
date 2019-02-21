package main

import (
	"tl/cli"
)

func main() {

	action := cli.ParseAction()
	records := cli.Init(action.TaskFilepath)

	cli.Run(records, action)

}
