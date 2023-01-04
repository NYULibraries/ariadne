package main

import (
	"ariadne/cmd"
	"ariadne/util"
)

func main() {
	util.InitSentry()
	cmd.Execute()
}
