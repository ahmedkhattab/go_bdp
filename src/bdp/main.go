package main

import (
  "fmt"
  "kube"
  "os"
  "ambari"
)

func main() {
  os.Setenv("KUBE", "/home/khattab/kubernetes-1.1.2/cluster/kubectl.sh")
	fmt.Println("Hello, world.")
  /*
  fmt.Println(kube.Get_remaining_pods("amb"))
  fmt.Println(kube.Get_pending_pods())

  fmt.Println(kube.Get_pod_ip("amb-server.service.consul"))
  fmt.Println(kube.Get_pod_host_ip("amb-server.service.consul"))
  fmt.Println(kube.Get_pod_status("amb-server.service.consul"))
  fmt.Println(kube.Get_pod_public_ip("amb-server.service.consul"))
  */
  ambari.Start()
  fmt.Println(kube.Get_pod_public_ip("amb-server.service.consul"))


}
