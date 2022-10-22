package main

import "fmt"

var Name string = "Plugin1"

type task struct {
	Parameters map[string]string
}

func (task) Exec() {
	fmt.Println("Plugin1 Exec()")
	fmt.Println("Running Plugin1")
}

func (task) GetName() string {
	return Name
}

func NewInstance(parameters map[string]string) interface{} {
	return task{
		Parameters: parameters,
	}
}
