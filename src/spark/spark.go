package spark

import (
	"kube"
	"log"
	"os"
	"time"
	"util"

	"github.com/spf13/viper"
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
	kube.CreateResource(os.Getenv("BDP_CONFIG_DIR") + "/spark/spark-master.json")

	log.Println("Spark: Waiting for spark master to start...")
	for {
		serverState := kube.PodStatus("spark-master")
		if serverState == "Running" {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}
	kube.CreateResource(os.Getenv("BDP_CONFIG_DIR") + "/spark/spark-master-service.json")

	log.Println("Spark: Launching spark workers")
	rc := util.LoadRC(os.Getenv("BDP_CONFIG_DIR") + "/spark/spark-worker-controller.json")
	rc.Spec.Replicas = viper.GetInt("SPARK_WORKERS")
	util.SaveRC(os.Getenv("BDP_CONFIG_DIR")+"/tmp/spark-worker-controller.json", rc)
	kube.CreateResource(os.Getenv("BDP_CONFIG_DIR") + "/tmp/spark-worker-controller.json")

	log.Println("Spark: Waiting for spark workers to start...")
	for {
		pending := kube.PendingPods()
		if pending == 0 {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}

	log.Println("Spark: Launching spark driver")
	kube.CreateResource(os.Getenv("BDP_CONFIG_DIR") + "/spark/spark-driver.json")

}
