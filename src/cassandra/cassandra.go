package cassandra

import (
	"fmt"
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
	util.ReleasePID("cassandra")
}

func Start(config util.Config, forceDeploy bool) {
	if !forceDeploy {
		if util.IsRunning("cassandra") {
			log.Println("Cassandra: already running, skipping start ...")
			return
		}
	}
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
	util.SetPID("cassandra")
	log.Println("Cassandra: Done!")

}

func Status() util.Status {
	status := util.Status{false, "Not Running", ""}
	if util.IsRunning("cassandra") {
		status.State = true
		status.Message = fmt.Sprintf("Cassandra accessible through %s:31317", kube.PodPublicIP("cassandra"))
	}
	return status
}
