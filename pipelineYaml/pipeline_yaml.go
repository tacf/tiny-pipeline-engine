package pipelineYaml

import (
	"log"
	"os"
	
	"gopkg.in/yaml.v2"
)

type TaskYaml struct {
	Name string `yaml:"task"`
	Parameters map[string]string
}

type TasksYaml = []TaskYaml

type PipelineYaml struct {
	Tasks TasksYaml `yaml:"steps"`
}

func check(e error) {
    if e != nil {
		log.Fatalf("error: %v", e)
        panic(e)
    }
}

func LoadYaml(filePath string) PipelineYaml {
	yamlObject := PipelineYaml{}

	dat, err := os.ReadFile(filePath)
    check(err)
    
	err = yaml.Unmarshal([]byte(string(dat)), &yamlObject)
	check(err)

	return yamlObject
}
