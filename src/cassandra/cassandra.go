package cassandra

import (
	"kube"
	"log"
	"time"
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
	kube.CreateResource("~/BDP/cassandra/cassandra-controller.yaml")
	kube.CreateResource("~/BDP/cassandra/cassandra-service.yaml")

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
