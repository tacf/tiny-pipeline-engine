package main

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/tacf/tiny-pipeline-engine/pipelineYaml"
)

type Plugin interface {
	Exec()
	GetName() string
}

type Arguments = map[string]string
type PluginLibs = map[string]func(Arguments) interface{}
type ExecutableTasks = []Plugin

func check(e error) {
	if e != nil {
		log.Fatalf("[ERROR][engine]: %v", e)
		panic(e)
	}
}

func locatePluginsFile(pluginBaseDir string) []string {
	var filepaths []string
	check(filepath.WalkDir(pluginBaseDir, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ".so" {
			filepaths = append(filepaths, s)
		}
		return nil
	}))
	return filepaths
}

func instanciateTasks(pluginLibs PluginLibs, tasks pipelineYaml.TasksYaml) ExecutableTasks {
	executableTasks := make(ExecutableTasks, len(tasks))
	i := 0
	for _, v := range tasks {
		pluginLib, found := pluginLibs[strings.ToLower(v.Name)]
		if !found {
			panic(fmt.Sprintf("Unable to locate specified Plugin Library '%s'", v.Name))
		}
		executableTasks[i] = pluginLib(v.Parameters).(Plugin)
		fmt.Printf("[engine] New  plugin <%s> instance created\n", strings.ToLower(v.Name))
		i++
	}
	return executableTasks
}

func loadPluginLibs(pluginPaths []string) PluginLibs {
	plugLibs := PluginLibs{}
	for _, path := range pluginPaths {
		p, err := plugin.Open(path)
		check(err)

		pluginName, err := p.Lookup("Name")
		check(err)

		pluginNewInstance, err := p.Lookup("NewInstance")
		check(err)

		pluginNewInstanceHandler := pluginNewInstance.(func(map[string]string) interface{})

		pName := *pluginName.(*string)
		plugLibs[strings.ToLower(pName)] = pluginNewInstanceHandler
		fmt.Printf("[engine] Plugin Lib '%s' loaded\n", strings.ToLower(pName))
	}
	return plugLibs
}

func executeTasks(executableTasks ExecutableTasks) {
	for _, executeTask := range executableTasks {
		fmt.Printf("[engine] Running plugin: '%s'\n", executeTask.GetName())

		// How do we handle exceptions here ? (Cross Language Exceptions)
		executeTask.Exec()
	}
}

func main() {
	pipelineYaml := pipelineYaml.LoadYaml("./pipeline.yaml")

	// Locate and load libs in plugin folder ('.so' files)
	pluginLibs := loadPluginLibs(locatePluginsFile("./bin/plugins/"))

	// Load Tasks
	executableTasks := instanciateTasks(pluginLibs, pipelineYaml.Tasks)

	// Execute the YAML workflow
	executeTasks(executableTasks)
}
