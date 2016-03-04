package kafka

import (
	"fmt"
	"kube"
	"log"
	"time"
	"util"

	"github.com/spf13/viper"
)

func CleanUp() {
	log.Println("Kafka: cleaning up cluster...")
	kube.DeleteResource("svc", "-l app=zk")
	kube.DeleteResource("pod", "-l app=zk")
	kube.DeleteResource("svc", "kafka")
	kube.DeleteResource("rc", "kafka-controller")

	for {
		remaining := kube.RemainingPods("zookeeper")
		if remaining == 0 {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}
	util.ReleasePID("kafka")
}

func Start(config util.Config, forceDeploy bool) {
	if !forceDeploy {
		if util.IsRunning("kafka") {
			log.Println("Kafka: already running, skipping start ...")
			return
		}
	}
	CleanUp()

	log.Println("Kafka: Launching zookeeper services")
	kube.CreateResource(viper.GetString("BDP_CONFIG_DIR") + "/kafka/zookeeper-service.json")
	log.Println("Kafka: Launching zookeeper servers")
	kube.CreateResource(viper.GetString("BDP_CONFIG_DIR") + "/kafka/zookeeper.json")

	log.Println("Kafka: Waiting for zookeeper server to start...")
	for {
		pending := kube.PendingPods()
		if pending == 0 {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}

	log.Println("Kafka: Launching kafka clients")

	log.Println("Kafka: Launching kafka service")
	kube.CreateResource(viper.GetString("BDP_CONFIG_DIR") + "/kafka/kafka-service.json")
	util.GenerateConfig("kafka.json", "kafka", config)
	kube.CreateResource(viper.GetString("BDP_CONFIG_DIR") + "/tmp/kafka.json")

	log.Println("Kafka: Waiting for kafka client to start...")
	for {
		pending := kube.PendingPods()
		if pending == 0 {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}
	util.SetPID("kafka")
	log.Println("Kafka: Done!")

}

func Status() util.Status {
	status := util.Status{false, "Not Running"}
	if util.IsRunning("kafka") {
		status.State = true
		status.Message = fmt.Sprintf("Kafka accessible through http://%s:31318\n", kube.PodPublicIP("kafka"))
	}
	return status
}
