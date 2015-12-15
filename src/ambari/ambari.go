package ambari

import (
  "kube"
  "time"
  "fmt"
)

func Clean_up() {
  fmt.Println("Ambari: cleaning up ...")
  kube.Delete_resource("rc","amb-slave-controller")
  kube.Delete_resource("svc","ambari")
  kube.Delete_resource("svc","consul")
  kube.Delete_resource("pods","amb-server.service.consul")
  kube.Delete_resource("pods","amb-consul")
  kube.Delete_resource("pods","amb-shell")
  for {
    remaining := kube.Get_remaining_pods("amb")
    if remaining == 0 {
      break
    } else {
      fmt.Printf(".")
      time.Sleep(5 * time.Second)
    }
  }
}

func Start() {
  Clean_up()

  fmt.Println("Ambari: Launching consul")
  kube.Create_resource("~/BDP/Ambari/consul.json")
  kube.Create_resource("~/BDP/Ambari/consul-service.json")

  fmt.Println("Ambari: Waiting for consul server to start")
  for {
    server_state := kube.Get_pod_status("amb-consul")
    if server_state == "Running" {
      fmt.Println()
      break
    } else {
      fmt.Printf(".")
      time.Sleep(5 * time.Second)
    }
  }

  fmt.Println("Ambari: Launching Ambari server")
  kube.Create_resource("~/BDP/Ambari/ambari-hdfs.json")
  kube.Create_resource("~/BDP/Ambari/ambari-service.json")

  fmt.Println("Ambari: Waiting for ambari server to start")
  for {
    server_state := kube.Get_pod_status("amb-server.service.consul")
    if server_state == "Running" {
      fmt.Println()
      break
    } else {
      fmt.Printf(".")
      time.Sleep(5 * time.Second)
    }
  }

  time.Sleep(10 * time.Second)
  fmt.Printf("Ambari: registering consul services")

  ambari_service_ip := kube.Get_service_ip("ambari")
  cmd := fmt.Sprintf("/bin/sh -c 'curl -X PUT -d '{\"Node\": \"ambari-8080\",\"Address\": \"%s\",\"Service\": {\"Service\": \"ambari-8080\"}}' http://$CONSUL_SERVICE_HOST:8500/v1/catalog/register'", ambari_service_ip)
  fmt.Printf(cmd)
  kube.Exec_on_pod("amb-server.service.consul", cmd)
}
