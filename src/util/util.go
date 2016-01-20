package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Replicationcontroller struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Labels interface{} `json:"labels"`
		Name   string      `json:"name"`
	} `json:"metadata"`
	Spec struct {
		Replicas int `json:"replicas"`
		Selector struct {
			Name string `json:"name"`
		} `json:"selector"`
		Template struct {
			Metadata struct {
				Labels interface{} `json:"labels"`
			} `json:"metadata"`
			Spec struct {
				Containers []struct {
					Command []string  `json:"command,omitempty"`
					Env     []*envVar `json:"env,omitempty"`
					Image   string    `json:"image"`
					Name    string    `json:"name"`
					Ports   []struct {
						ContainerPort int `json:"containerPort"`
					} `json:"ports"`
					Resources *resources `json:"resources,omitempty"`
				} `json:"containers"`
			} `json:"spec"`
		} `json:"template"`
	} `json:"spec"`
}
type resources struct {
	Limits *limits `json:"limits,omitempty"`
}

type limits struct {
	CPU string `json:"cpu,omitempty"`
}

type envVar struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

func LoadRC(jsonfile string) Replicationcontroller {
	content, err := ioutil.ReadFile(jsonfile)
	if err != nil {
		fmt.Print("Error:", err)
	}
	var rc Replicationcontroller
	err = json.Unmarshal(content, &rc)
	if err != nil {
		log.Fatal("Error:", err)
	}
	fmt.Println(rc)
	return rc
}

func SaveRC(filename string, rc Replicationcontroller) {
	data, err := json.Marshal(rc)
	if err != nil {
		log.Fatal("Error:", err)
	}
	err = ioutil.WriteFile(filename, data, 0644)
}
