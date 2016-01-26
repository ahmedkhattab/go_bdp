package cassandra

import (
	"kube"
	"log"
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

func Start(config util.Config) {
	CleanUp()

	log.Println("Cassandra: Launching cassandra pods")
	util.GenerateConfig("cassandra-controller.json", "cassandra", config)
	kube.CreateResource(viper.GetString("BDP_CONFIG_DIR") + "/tmp/cassandra-controller.json")

	log.Println("Cassandra: Waiting for cassandra pods to start...")
	for {
		pending := kube.PendingPods()
		if pending == 0 {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}
	kube.CreateResource(viper.GetString("BDP_CONFIG_DIR") + "/cassandra/cassandra-service.json")

}
