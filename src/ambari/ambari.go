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
  fmt.Println("Ambari: registering consul services")

  ambari_service_ip := kube.Get_service_ip("ambari")
  cmd := fmt.Sprintf("'curl -X PUT -d \"{\\\"Node\\\": \\\"ambari-8080\\\",\\\"Address\\\": \\\"%s\\\",\\\"Service\\\": {\\\"Service\\\": \\\"ambari-8080\\\"}}\" http://$CONSUL_SERVICE_HOST:8500/v1/catalog/register'",
    ambari_service_ip)
  fmt.Println(kube.Exec_on_pod("amb-server.service.consul", cmd))
  cmd = fmt.Sprintf("'curl -X PUT -d \"{\\\"Node\\\": \\\"amb-server\\\",\\\"Address\\\": \\\"%s\\\",\\\"Service\\\": {\\\"Service\\\": \\\"amb-server\\\"}}\" http://$CONSUL_SERVICE_HOST:8500/v1/catalog/register'",
      ambari_service_ip)
  fmt.Println(kube.Exec_on_pod("amb-server.service.consul", cmd))

  fmt.Println("Ambari: launching ambari slaves")
  kube.Create_resource("~/BDP/Ambari/ambari-slave.json")
  for {
    pending := kube.Get_pending_pods()
    if pending == 0 {
      fmt.Println()
      break
    } else {
      fmt.Printf(".")
      time.Sleep(5 * time.Second)
    }
  }

  fmt.Println("Ambari: creating ambari cluster using blueprint: multi-node-hdfs")
  kube.Create_resource("~/BDP/Ambari/ambari-shell.json")

  slave_pods := kube.Get_pod_names("amb-slave")
  for v := 0; v < len(slave_pods); v++ {
    if slave_pods[v] != "" {
     cmd = "'curl -X PUT -d \"{\\\"Node\\\": \\\"$(hostname)\\\",\\\"Address\\\": \\\"$(hostname -I)\\\",\\\"Service\\\": {\\\"Service\\\": \\\"$(hostname)\\\"}}\" http://$CONSUL_SERVICE_HOST:8500/v1/catalog/register'"
     fmt.Println(kube.Exec_on_pod(slave_pods[v], cmd))
   }
  }

}
