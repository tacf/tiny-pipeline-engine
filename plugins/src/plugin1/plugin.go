package main

import "fmt"

var Name string

type plugin struct {}

func (plugin) Exec(parameters map[string]string) {
        fmt.Println("Plugin1 Exec()")
        fmt.Println("Running Plugin1")
}

func (plugin) GetName() string {
        return Name
}


func Initialize() (f interface{}, err error) {
        Name = "plugin1"
        f = plugin{}
        return
}
