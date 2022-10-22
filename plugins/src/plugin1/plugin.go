package main

import "fmt"

var Name string = "Plugin1"

type plugin struct {
	Parameters map[string]string
}

func (plugin) Exec() {
	fmt.Println("Plugin1 Exec()")
	fmt.Println("Running Plugin1")
}

func NewInstance(parameters map[string]string) interface{} {
	return plugin{
		Parameters: parameters,
	}
}
