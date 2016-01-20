package cassandra

import (
	"kube"
	"log"
	"os"
	"time"

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
	kube.CreateResource(os.Getenv("BDP_CONFIG_DIR") + "/cassandra/cassandra-controller.yaml")
	kube.CreateResource(os.Getenv("BDP_CONFIG_DIR") + "/cassandra/cassandra-service.yaml")
	if viper.GetInt("CASSANDRA_NODES") != 2 {
		kube.ScaleController("cassandra", viper.GetInt("CASSANDRA_NODES"))
	}
	log.Println("Cassandra: Waiting for cassandra pods to start...")
	for {
		pending := kube.PendingPods()
		if pending == 0 {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}

}
