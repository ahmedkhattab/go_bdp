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
		Labels interface{} `json:"labels,omitempty"`
		Name   string      `json:"name"`
	} `json:"metadata"`
	Spec struct {
		Replicas int       `json:"replicas"`
		Selector *selector `json:"selector,omitempty"`
		Template struct {
			Metadata struct {
				Labels interface{} `json:"labels,omitempty"`
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
					Resources    *resources     `json:"resources,omitempty"`
					VolumeMounts []*volumeMount `json:"volumeMounts,omitempty"`
				} `json:"containers"`
				Volumes []*volume `json:"volumes,omitempty"`
			} `json:"spec"`
		} `json:"template"`
	} `json:"spec"`
}

type volume struct {
	EmptyDir struct{} `json:"emptyDir,omitempty"`
	Name     string   `json:"name,omitempty"`
}

type volumeMount struct {
	MountPath string `json:"mountPath,omitempty"`
	Name      string `json:"name,omitempty"`
}

type resources struct {
	Limits *limits `json:"limits,omitempty"`
}

type limits struct {
	CPU string `json:"cpu,omitempty"`
}

type envVar struct {
	Name      string     `json:"name,omitempty"`
	Value     string     `json:"value,omitempty"`
	ValueFrom *valueFrom `json:"valueFrom,omitempty"`
}

type valueFrom struct {
	FieldRef struct {
		FieldPath string `json:"fieldPath,omitempty"`
	} `json:"fieldRef,omitempty"`
}

type selector struct {
	Name string `json:"name,omitempty"`
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
	return rc
}

func SaveRC(filename string, rc Replicationcontroller) {
	data, err := json.Marshal(rc)
	if err != nil {
		log.Fatal("Error:", err)
	}
	err = ioutil.WriteFile(filename, data, 0644)
}
