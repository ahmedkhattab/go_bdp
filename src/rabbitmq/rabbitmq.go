package rabbitmq

import (
	"kube"
	"log"
	"time"
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
	kube.CreateResource("~/BDP/rabbitmq/rabbitmq-controller.json")
	kube.CreateResource("~/BDP/rabbitmq/rabbitmq-service.json")

	log.Println("Rabbitmq: Waiting for Rabbitmq pods to start...")
	for {
		pending := kube.PendingPods()
		if pending == 0 {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}
}
