package rabbitmq

import (
	"kube"
	"log"
	"os"
	"time"
	"util"

	"github.com/spf13/viper"
)

func CleanUp() {
	log.Println("Rabbitmq: cleaning up cluster...")
	kube.DeleteResource("rc", "rabbitmq-controller")
	kube.DeleteResource("svc", "rabbitmq")

	for {
		remaining := kube.RemainingPods("rabbitmq")
		if remaining == 0 {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}
}

func Start() {
	CleanUp()

	log.Println("Rabbitmq: Launching rabbitmq")

	rc := util.LoadRC(os.Getenv("BDP_CONFIG_DIR") + "/rabbitmq/rabbitmq-controller.json")
	rc.Spec.Replicas = viper.GetInt("RABBITMQ_NODES")
	util.SaveRC(os.Getenv("BDP_CONFIG_DIR")+"/tmp/rabbitmq-controller.json", rc)
	kube.CreateResource(os.Getenv("BDP_CONFIG_DIR") + "/tmp/rabbitmq-controller.json")

	log.Println("Rabbitmq: Waiting for Rabbitmq pods to start...")
	for {
		pending := kube.PendingPods()
		if pending == 0 {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}

	kube.CreateResource(os.Getenv("BDP_CONFIG_DIR") + "/rabbitmq/rabbitmq-service.json")

}
