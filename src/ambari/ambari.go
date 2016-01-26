package ambari

import (
	"fmt"
	"kube"
	"log"
	"time"
	"util"

	"github.com/spf13/viper"
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

func Start(config util.Config) {
	CleanUp()

	log.Println("Ambari: Launching consul")
	kube.CreateResource(viper.GetString("BDP_CONFIG_DIR") + "/ambari/consul.json")
	kube.CreateResource(viper.GetString("BDP_CONFIG_DIR") + "/ambari/consul-service.json")

	log.Println("Ambari: Waiting for consul server to start...")
	for {
		serverState := kube.PodStatus("amb-consul")
		if serverState == "Running" {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}

	log.Println("Ambari: Launching Ambari server")
	kube.CreateResource(viper.GetString("BDP_CONFIG_DIR") + "/ambari/ambari-hdfs.json")
	kube.CreateResource(viper.GetString("BDP_CONFIG_DIR") + "/ambari/ambari-service.json")

	log.Println("Ambari: Waiting for ambari server to start...")
	for {
		serverState := kube.PodStatus("amb-server.service.consul")
		if serverState == "Running" {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}

	time.Sleep(10 * time.Second)
	log.Println("Ambari: registering consul services")

	ambariServiceIP := kube.ServiceIP("ambari")
	cmd := fmt.Sprintf("'curl -X PUT -d \"{\\\"Node\\\": \\\"ambari-8080\\\",\\\"Address\\\": \\\"%s\\\",\\\"Service\\\": {\\\"Service\\\": \\\"ambari-8080\\\"}}\" http://$CONSUL_SERVICE_HOST:8500/v1/catalog/register'",
		ambariServiceIP)
	kube.ExecOnPod("amb-server.service.consul", cmd)
	cmd = fmt.Sprintf("'curl -X PUT -d \"{\\\"Node\\\": \\\"amb-server\\\",\\\"Address\\\": \\\"%s\\\",\\\"Service\\\": {\\\"Service\\\": \\\"amb-server\\\"}}\" http://$CONSUL_SERVICE_HOST:8500/v1/catalog/register'",
		ambariServiceIP)
	kube.ExecOnPod("amb-server.service.consul", cmd)

	log.Println("Ambari: launching ambari slaves...")

	util.GenerateConfig("ambari-slave.json", "ambari", config)

	kube.CreateResource(viper.GetString("BDP_CONFIG_DIR") + "/tmp/ambari-slave.json")

	for {
		pending := kube.PendingPods()
		if pending == 0 {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}

	log.Println("Ambari: creating ambari cluster using blueprint: multi-node-hdfs")
	kube.CreateResource(viper.GetString("BDP_CONFIG_DIR") + "/ambari/ambari-shell.json")

	slavePods := kube.PodNames("amb-slave")
	for v := 0; v < len(slavePods); v++ {
		if slavePods[v] != "" {
			cmd = "'curl -X PUT -d \"{\\\"Node\\\": \\\"$(hostname)\\\",\\\"Address\\\": \\\"$(hostname -I)\\\",\\\"Service\\\": {\\\"Service\\\": \\\"$(hostname)\\\"}}\" http://$CONSUL_SERVICE_HOST:8500/v1/catalog/register'"
			kube.ExecOnPod(slavePods[v], cmd)
		}
	}

}
