package kube

import (
  "os"
  "os/exec"
  "log"
  "strings"
  "fmt"
)

func kube(cmd string, args ...string) (*exec.Cmd) {
  if len(args) > 0 {
    args_s := strings.Join(args, " ")
    return exec.Command("sh", "-c", os.Getenv("KUBE") + " " + cmd + " " + args_s)
  }
  return exec.Command("sh", "-c", os.Getenv("KUBE") + " " + cmd)
}

func pipe_commands(commands ...*exec.Cmd) ([]byte, error) {
  for i, command := range commands[:len(commands) - 1] {
      out, err := command.StdoutPipe()
      if err != nil {
          return nil, err
      }
      command.Start()
      commands[i + 1].Stdin = out
  }
  final, err := commands[len(commands) - 1].Output()
  if err != nil {
      return nil, err
  }
  return final, nil
}

func Get_pod_status (pod string) string {
  out, err := kube("get pod", pod, "-o template", "--template={{.status.phase}}").Output()
  if err != nil {
			log.Fatal(err)
  }
  return string(out)
}

func Get_pod_ip (pod string) string {
  out, err := kube("get pod", pod, "-o template", "--template={{.status.podIP}}").Output()
  if err != nil {
			log.Fatal(err)
  }
  return string(out)
}

func Get_pod_host_ip (pod string) string {
  out, err := kube("get pod", pod, "-o template", "--template={{.status.hostIP}}").Output()
  if err != nil {
			log.Fatal(err)
  }
  return string(out)
}

func Get_pod_host_name (pod string) string {
  out, err := kube("get pod", pod, "-o template", "--template={{.spec.nodeName}}").Output()
  if err != nil {
			log.Fatal(err)
  }
  return string(out)
}

func Get_pod_public_ip (pod string) string {
  host_name := Get_pod_host_name(pod)
  out, err := kube("get node", host_name, "-o template", "'--template={{(index .status.addresses 2).address}}'").Output()
  if err != nil {
			log.Fatal(err)
  }
  return string(out)
}

func Get_node (pod string) string {
  out, err := kube("get pod", pod, "-o template", "--template={{.status.hostIP}}").Output()
  if err != nil {
			log.Fatal(err)
  }
  return string(out)
}

func Get_pending_pods () string {
  pods := kube("get pods")
  pending := exec.Command("grep", "Pending")
  count := exec.Command("wc", "-l")
  out, err := pipe_commands(pods, pending, count)
  if err != nil {
          log.Fatal(err)
      }
  return string(out)
}

func Create_resource (path string) string {
    out, err := kube("create", "-f", path).Output()
    if err != nil {
            log.Fatal(err)
        }
    return string(out)
}
