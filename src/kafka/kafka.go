package kafka

import (
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
}

func Start(config util.Config) {
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

	//util.GenerateConfig("spark-worker-controller.json", "spark", config)
	log.Println("Kafka: Launching kafka service")
	kube.CreateResource(viper.GetString("BDP_CONFIG_DIR") + "/kafka/kafka-service.json")
	kube.CreateResource(viper.GetString("BDP_CONFIG_DIR") + "/kafka/kafka.json")

	log.Println("Kafka: Waiting for kafka client to start...")
	for {
		pending := kube.PendingPods()
		if pending == 0 {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}
}
