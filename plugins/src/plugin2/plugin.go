package main

import (
	"fmt"
	"os/exec"
)

var Name string = "Plugin2"

type task struct {
	Parameters map[string]string
}

func (p task) Exec() {
	fmt.Println("Plugin2 Exec() - Executing shell command")
	shell := "bash"

	cmd := exec.Command(shell, "-c", p.Parameters["command"])
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Print the output
	fmt.Println(string(stdout))
}

func (task) GetName() string {
	return Name
}

func NewInstance(parameters map[string]string) interface{} {
	return task{
		Parameters: parameters,
	}
}
