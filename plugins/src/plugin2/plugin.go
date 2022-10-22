package main

import (
	"fmt"
	"os/exec"
)

var Name string

type plugin struct{}

func (plugin) Exec(parameters map[string]string) {
	fmt.Println("Plugin2 Exec() - Executing shell command")
	shell := "bash"

	cmd := exec.Command(shell, "-c", parameters["command"])
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Print the output
	fmt.Println(string(stdout))
}

func (plugin) GetName() string {
	return Name
}

func Initialize() (f interface{}, err error) {
	Name = "plugin2"
	f = plugin{}
	return
}
