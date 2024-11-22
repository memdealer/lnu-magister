package models

/*
---
default:
  vm_default_specs: &vm_default_specs
    cpu: 4
    memory: 6144
    display: "1024x768"


Tart:
  - name: "test-runner"
    state: "running" # stopped, running, absent
    image: "tart-runner"
    specs: *vm_default_specs

  - name: "test-runner-2"
    state: "running" # stopped, running, absent
    image: "tart-runner"
    specs: *vm_default_specs
*/

type Specs struct {
	Cpu     int    `yaml:"cpu"`
	Memory  int    `yaml:"memory"`
	Display string `yaml:"display"`
}

type Runner struct {
	HostName string
	Name     string `yaml:"name"`
	State    string `yaml:"state"`
	Image    string `yaml:"image"`
	Specs    Specs  `yaml:"specs"`
}

type Hostname struct {
	HostName string
}

type TartState struct {
	HostName string   // Unique ID for the state
	Runners  []Runner `yaml:"TartState"`
}
