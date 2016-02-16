package spark

import (
	"kube"
	"log"
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
	util.ReleasePID("spark")
}

func Start(config util.Config) {
	if util.IsRunning("spark") {
		log.Println("Spark: already running, skipping start ...")
		return
	}
	CleanUp()

	log.Println("Spark: Launching spark master")
	kube.CreateResource(viper.GetString("BDP_CONFIG_DIR") + "/spark/spark-master.json")

	log.Println("Spark: Waiting for spark master to start...")
	for {
		serverState := kube.PodStatus("spark-master")
		if serverState == "Running" {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}
	kube.CreateResource(viper.GetString("BDP_CONFIG_DIR") + "/spark/spark-master-service.json")

	log.Println("Spark: Launching spark workers")

	util.GenerateConfig("spark-worker-controller.json", "spark", config)
	kube.CreateResource(viper.GetString("BDP_CONFIG_DIR") + "/tmp/spark-worker-controller.json")

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
	kube.CreateResource(viper.GetString("BDP_CONFIG_DIR") + "/spark/spark-driver.json")
	util.SetPID("spark")
}
