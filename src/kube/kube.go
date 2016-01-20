package kube

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

func kube(cmd string, args ...string) *exec.Cmd {
	if len(args) > 0 {
		argsJoined := strings.Join(args, " ")
		return exec.Command("sh", "-c", viper.GetString("KUBE_PATH")+"/cluster/kubectl.sh"+" "+cmd+" "+argsJoined)
	}
	return exec.Command("sh", "-c", viper.GetString("KUBE_PATH")+"/cluster/kubectl.sh"+" "+cmd)
}

func pipeCommands(commands ...*exec.Cmd) ([]byte, error) {
	for i, command := range commands[:len(commands)-1] {
		out, err := command.StdoutPipe()
		if err != nil {
			return nil, err
		}
		command.Start()
		commands[i+1].Stdin = out
	}
	var stderr bytes.Buffer
	cmd := commands[len(commands)-1]
	cmd.Stderr = &stderr
	final, err := cmd.Output()
	if err != nil {
		fmt.Println(stderr.String())
		return nil, err
	}
	return final, nil
}

//ClusterInfo executes kubectl cluster-info
func ClusterInfo() string {
	out, err := kube("cluster-info").Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

//ClusterIsUp checks if a kubernetes cluster is up and running
func ClusterIsUp() bool {
	_, err := kube("cluster-info").Output()
	if err != nil {
		return false
	}
	return true
}

//GetPods executes kubectl get pods
func GetPods() string {
	out, err := kube("get pods").Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

//PodNames returns the names of all pods starting with certain prefix
func PodNames(prefix string) []string {
	pods := kube("get pods", "--no-headers")
	cut := exec.Command("cut", "-d", " ", "-f", "1")
	grep := exec.Command("grep", "amb-slave")
	out, err := pipeCommands(pods, cut, grep)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(out), "\n")
}

//PodStatus returns the status of the input pod
func PodStatus(pod string) string {
	out, err := kube("get pod", pod, "-o template", "--template={{.status.phase}}").Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

//PodIP returns the ip of the input pod
func PodIP(pod string) string {
	out, err := kube("get pod", pod, "-o template", "--template={{.status.podIP}}").Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

//PodHostIP returns the ip of the host of the input pod
func PodHostIP(pod string) string {
	out, err := kube("get pod", pod, "-o template", "--template={{.status.hostIP}}").Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

//PodHostName returns the node name of the input pod
func PodHostName(pod string) string {
	out, err := kube("get pod", pod, "-o template", "--template={{.spec.nodeName}}").Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

//PodPublicIP returns the public ip of the host of the input pod
func PodPublicIP(pod string) string {
	hostName := PodHostName(pod)
	out, err := kube("get node", hostName, "-o template", "'--template={{(index .status.addresses 2).address}}'").Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

func ServiceIP(service string) string {
	out, err := kube("get service", service, "-o=template", "--template={{.spec.clusterIP}}").Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

func PendingPods() int {
	pods := kube("get pods")
	pending := exec.Command("grep", "Pending")
	count := exec.Command("wc", "-l")
	out, err := pipeCommands(pods, pending, count)
	if err != nil {
		log.Fatal(err)
	}
	num, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		log.Fatal(err)
	}
	return num
}

func RemainingPods(prefix string) int {
	pods := kube("get pods", "--no-headers")
	grepPrefix := exec.Command("grep", prefix)
	count := exec.Command("wc", "-l")
	out, err := pipeCommands(pods, grepPrefix, count)
	if err != nil {
		log.Fatal(err)
	}
	num, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		log.Fatal(err)
	}
	return num
}

func DeleteResource(rType string, rName string) string {
	log.Println("Attempting to delete " + rName)
	cmd := kube("delete", rType, rName)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		if !strings.Contains(stderr.String(), "not found") {
			log.Println(stderr.String())
		}
	}
	return string(out)
}

func CreateResource(path string) string {
	cmd := kube("create", "-f", path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		log.Println(stderr.String())
	}
	return string(out)
}

func ScaleController(rcName string, size int) string {
	replicasParam := fmt.Sprintf("--replicas=%d", size)
	out, err := kube("scale", replicasParam, "rc", rcName).Output()
	if err != nil {
		log.Println(err)
	}
	return string(out)
}
func ExecOnPod(pod string, command string) string {
	cmd := kube("exec", pod, "--", "/bin/sh", "-c", command)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		log.Println(err)
		log.Println(stderr.String())
	}
	return string(out)
}
