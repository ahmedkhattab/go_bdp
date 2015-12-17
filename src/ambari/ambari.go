package ambari

import (
	"fmt"
	"kube"
	"log"
	"time"
)

func CleanUp() {
	log.Println("Ambari: cleaning up cluster...")
	kube.DeleteResource("rc", "amb-slave-controller")
	kube.DeleteResource("svc", "ambari")
	kube.DeleteResource("svc", "consul")
	kube.DeleteResource("pods", "amb-server.service.consul")
	kube.DeleteResource("pods", "amb-consul")
	kube.DeleteResource("pods", "amb-shell")
	for {
		remaining := kube.RemainingPods("amb")
		if remaining == 0 {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}
}

func Start() {
	CleanUp()

	log.Println("Ambari: Launching consul")
	kube.CreateResource("~/BDP/Ambari/consul.json")
	kube.CreateResource("~/BDP/Ambari/consul-service.json")

	log.Println("Ambari: Waiting for consul server to start...")
	for {
		server_state := kube.PodStatus("amb-consul")
		if server_state == "Running" {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}

	log.Println("Ambari: Launching Ambari server")
	kube.CreateResource("~/BDP/Ambari/ambari-hdfs.json")
	kube.CreateResource("~/BDP/Ambari/ambari-service.json")

	log.Println("Ambari: Waiting for ambari server to start...")
	for {
		server_state := kube.PodStatus("amb-server.service.consul")
		if server_state == "Running" {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}

	time.Sleep(10 * time.Second)
	log.Println("Ambari: registering consul services")

	ambari_service_ip := kube.ServiceIP("ambari")
	cmd := fmt.Sprintf("'curl -X PUT -d \"{\\\"Node\\\": \\\"ambari-8080\\\",\\\"Address\\\": \\\"%s\\\",\\\"Service\\\": {\\\"Service\\\": \\\"ambari-8080\\\"}}\" http://$CONSUL_SERVICE_HOST:8500/v1/catalog/register'",
		ambari_service_ip)
	kube.ExecOnPod("amb-server.service.consul", cmd)
	cmd = fmt.Sprintf("'curl -X PUT -d \"{\\\"Node\\\": \\\"amb-server\\\",\\\"Address\\\": \\\"%s\\\",\\\"Service\\\": {\\\"Service\\\": \\\"amb-server\\\"}}\" http://$CONSUL_SERVICE_HOST:8500/v1/catalog/register'",
		ambari_service_ip)
	kube.ExecOnPod("amb-server.service.consul", cmd)

	log.Println("Ambari: launching ambari slaves...")
	kube.CreateResource("~/BDP/Ambari/ambari-slave.json")
	for {
		pending := kube.PendingPods()
		if pending == 0 {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}

	log.Println("Ambari: creating ambari cluster using blueprint: multi-node-hdfs")
	kube.CreateResource("~/BDP/Ambari/ambari-shell.json")

	slave_pods := kube.PodNames("amb-slave")
	for v := 0; v < len(slave_pods); v++ {
		if slave_pods[v] != "" {
			cmd = "'curl -X PUT -d \"{\\\"Node\\\": \\\"$(hostname)\\\",\\\"Address\\\": \\\"$(hostname -I)\\\",\\\"Service\\\": {\\\"Service\\\": \\\"$(hostname)\\\"}}\" http://$CONSUL_SERVICE_HOST:8500/v1/catalog/register'"
			kube.ExecOnPod(slave_pods[v], cmd)
		}
	}

}
