package spark

import (
	"fmt"
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

func Start(config util.Config, forceDeploy bool) {
	if !forceDeploy {
		if util.IsRunning("spark") {
			log.Println("Spark: already running, skipping start ...")
			return
		}
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
	time.Sleep(5 * time.Second)
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
	log.Println("Spark: Done!")

}

func Status() util.Status {
	status := util.Status{State: false, Message: "Not Running", URL: ""}
	if util.IsRunning("spark") {
		status.State = true
		status.Message = fmt.Sprintf("Spark UI accessible through ")
		status.URL = fmt.Sprintf("http://%s:31314", kube.PodPublicIP("spark-master"))
	}
	return status
}

func RunApp(gitRep string, pathToJar string) {

	cloneCmd := "'git clone " + gitRep + "'"
	_, err := kube.ExecOnPod("spark-driver", cloneCmd)
	if err != "" {
		CleanDriver("bdp_apps/")
	}
	paramsJoined := ""
	/*
		if len(params) > 0 {
			paramsJoined = strings.Join(params, " ")
		}*/
	submitCmd := "'spark-submit " + pathToJar + " " + paramsJoined + "'"
	kube.ExecOnPod("spark-driver", submitCmd)

}

func CleanDriver(rep string) {
	rmCmd := "'rm -rf " + rep + "'"
	kube.ExecOnPod("spark-driver", rmCmd)
}
