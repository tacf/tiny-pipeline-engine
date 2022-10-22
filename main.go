package main

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"plugin"

	"tiagoacf.com/pipelineYaml"
	"tiagoacf.com/types"
)

type Plugins = map[string]types.Plugin
type ExecutableTasks = []struct {
	types.Plugin
	pipelineYaml.TaskYaml
}

func check(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
		panic(e)
	}
}

func loadTasks(plugins Plugins, tasks pipelineYaml.TasksYaml) ExecutableTasks {
	executableTasks := make(ExecutableTasks, len(tasks))
	i := 0
	for _, v := range tasks {
		plugin, found := plugins[v.Name]
		if !found {
			panic(fmt.Sprintf("Unable to locate specified task '%s'", v.Name))
		}
		executableTasks[i] = struct {
			types.Plugin
			pipelineYaml.TaskYaml
		}{plugin, v}
		i++
	}
	return executableTasks
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

func loadPlugins(pluginPaths []string) Plugins {
	plugs := Plugins{}
	for _, path := range pluginPaths {
		p, err := plugin.Open(path)
		check(err)

		GetPluginIface, err := p.Lookup("Initialize")
		check(err)

		pluginIface, err := GetPluginIface.(func() (interface{}, error))()
		fmt.Printf("[engine] Loading plugin @ %s, err: %v\n", path, err)

		plug, ok := pluginIface.(types.Plugin)
		if !ok {
			panic("[engine] Plugin loading failed due to 'Interface Implementation Error'")
		}
		fmt.Printf("[plugin %s] Loaded Plugin: %T %v\n", plug.GetName(), plug, plug)

		plugs[plug.GetName()] = plug
	}
	return plugs
}

func executePlugins(tasks ExecutableTasks) {
	for _, plugin := range tasks {
		args := plugin.TaskYaml.Parameters
		fmt.Printf("[engine] Running plugin: '%s'\n", plugin.GetName())

		// How do we handle exceptions here ? (Cross Language Exceptions)
		plugin.Exec(args)
	}
}

func main() {
	pipelineYaml := pipelineYaml.LoadYaml("./pipeline.yaml")

	// Plugin loader operations

	// Locate libs in plugin folder ('.so' files)
	pluginFiles := locatePluginsFile("./bin/plugins/")

	// Load Plugins
	plugins := loadPlugins(pluginFiles)

	// Pipeline execution operations

	// Load Tasks
	executableTasks := loadTasks(plugins, pipelineYaml.Tasks)

	executePlugins(executableTasks)
}
