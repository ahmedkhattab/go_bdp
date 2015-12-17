package spark

import (
	"kube"
	"log"
	"time"
)

func CleanUp() {
	log.Println("Spark: cleaning up cluster...")
	kube.DeleteResource("rc", "spark-worker-controller")
	kube.DeleteResource("svc", "spark-master")
	kube.DeleteResource("pod", "spark-master")
	kube.DeleteResource("pod", "spark-driver")
	for {
		remaining := kube.RemainingPods("spark")
		if remaining == 0 {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}
}

func Start() {
	CleanUp()

	log.Println("Spark: Launching spark master")
	kube.CreateResource("~/BDP/spark/spark-master.json")
	kube.CreateResource("~/BDP/spark/spark-master-service.json")

	log.Println("Spark: Waiting for spark master to start...")
	for {
		serverState := kube.PodStatus("spark-master")
		if serverState == "Running" {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}

	log.Println("Spark: Launching spark workers")
	kube.CreateResource("~/BDP/spark/spark-worker-controller.json")
	log.Println("Spark: Waiting for 3 spark workers to start...")
	for {
		pending := kube.PendingPods()
		if pending == 0 {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}

	log.Println("Spark: Launching spark driver")
	kube.CreateResource("~/BDP/spark/spark-driver.json")

}
