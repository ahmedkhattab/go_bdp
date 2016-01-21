package cassandra

import (
	"kube"
	"log"
	"os"
	"time"
	"util"

	"github.com/spf13/viper"
)

func CleanUp() {
	log.Println("Cassandra: cleaning up cluster...")
	kube.DeleteResource("rc", "cassandra")
	kube.DeleteResource("svc", "cassandra")

	for {
		remaining := kube.RemainingPods("cassandra")
		if remaining == 0 {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}
}

func Start() {
	CleanUp()

	log.Println("Cassandra: Launching cassandra pods")
	rc := util.LoadRC(os.Getenv("BDP_CONFIG_DIR") + "/cassandra/cassandra-controller.json")
	rc.Spec.Replicas = viper.GetInt("CASSANDRA_NODES")
	util.SaveRC(os.Getenv("BDP_CONFIG_DIR")+"/tmp/cassandra-controller.json", rc)
	kube.CreateResource(os.Getenv("BDP_CONFIG_DIR") + "/tmp/cassandra-controller.json")

	log.Println("Cassandra: Waiting for cassandra pods to start...")
	for {
		pending := kube.PendingPods()
		if pending == 0 {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}
	kube.CreateResource(os.Getenv("BDP_CONFIG_DIR") + "/cassandra/cassandra-service.json")

}
