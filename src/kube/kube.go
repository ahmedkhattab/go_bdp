package kube

import (
  "os"
  "os/exec"
  "log"
  "strings"
  "strconv"
  "bytes"
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
  var stderr bytes.Buffer
  cmd := commands[len(commands) - 1]
  cmd.Stderr = &stderr
  final, err := cmd.Output()
  if err != nil {
    fmt.Println(stderr.String())
    return nil, err
  }
  return final, nil
}

func Get_pods () string {
  out, err := kube("get pods").Output()
  if err != nil {
			log.Fatal(err)
  }
  return string(out)
}

func Get_pod_names (prefix string) []string {
  pods := kube("get pods", "--no-headers")
  cut := exec.Command("cut", "-d", " ", "-f", "1")
  grep := exec.Command("grep", "amb-slave")
  out, err := pipe_commands(pods, cut, grep)
  if err != nil {
      log.Fatal(err)
  }
  return strings.Split(string(out), "\n")
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

func Get_service_ip (service string) string {
  out, err := kube("get service", service, "-o=template", "--template={{.spec.clusterIP}}").Output()
  if err != nil {
			log.Fatal(err)
  }
  return string(out)
}

func Get_pending_pods () int {
  pods := kube("get pods")
  pending := exec.Command("grep", "Pending")
  count := exec.Command("wc", "-l")
  out, err := pipe_commands(pods, pending, count)
  if err != nil {
      log.Fatal(err)
  }
  num, err := strconv.Atoi(strings.TrimSpace(string(out)))
  if err != nil {
      log.Fatal(err)
  }
  return num
}

func Get_remaining_pods (prefix string) int {
  pods := kube("get pods", "--no-headers")
  grep_prefix := exec.Command("grep", prefix)
  count := exec.Command("wc", "-l")
  out, err := pipe_commands(pods, grep_prefix, count)
  if err != nil {
      log.Fatal(err)
  }
  num, err := strconv.Atoi(strings.TrimSpace(string(out)))
  if err != nil {
      log.Fatal(err)
  }
  return num
}

func Delete_resource (r_type string, r_name string) string {
  out, err := kube("delete", r_type, r_name).Output()
  if err != nil {
      log.Println(err)
  }
  return string(out)
}

func Create_resource (path string) string {
  out, err := kube("create", "-f", path).Output()
  if err != nil {
      log.Println(err)
  }
  return string(out)
}

func Exec_on_pod (pod string, command string) string {
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
