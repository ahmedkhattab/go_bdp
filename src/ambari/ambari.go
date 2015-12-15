package ambari

import (
  "kube"
  "time"
  "fmt"
)

func Clean_up() {
  fmt.Println("Ambari: cleaning up ...")
  fmt.Printf(kube.Delete_resource("rc","amb-slave-controller"))
  fmt.Printf(kube.Delete_resource("svc","ambari"))
  fmt.Printf(kube.Delete_resource("svc","consul"))
  fmt.Printf(kube.Delete_resource("pods","amb-server.service.consul"))
  fmt.Printf(kube.Delete_resource("pods","amb-consul"))
  fmt.Printf(kube.Delete_resource("pods","amb-shell"))
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
