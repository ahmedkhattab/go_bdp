package launcher

import (
	"ambari"
	"cassandra"
	"fmt"
	"kafka"
	"kube"
	"os"
	"path/filepath"
	"rabbitmq"
	"spark"
	"util"

	"github.com/spf13/viper"
)

func LaunchComponents(allFlag bool, forceFlag bool, config util.Config) string {
	stdout := ""
	if kube.ClusterIsUp() {
		os.Mkdir(filepath.Join(viper.GetString("BDP_CONFIG_DIR"), "tmp"), 0777)
		if allFlag || config.Ambari {
			ambari.Start(config, forceFlag)
			stdout += fmt.Sprintf("Ambari UI accessible through http://%s:31313\n", kube.PodPublicIP("amb-server"))
		}
		if allFlag || config.Rabbitmq {
			rabbitmq.Start(config, forceFlag)
			stdout += fmt.Sprintf("RabbitMQ UI accessible through http://%s:31316\n", kube.PodPublicIP("rabbitmq"))
		}
		if allFlag || config.Kafka {
			kafka.Start(config, forceFlag)
			stdout += fmt.Sprintf("Kafka accessible through http://%s:31318\n", kube.PodPublicIP("kafka"))
		}
		if allFlag || config.Spark {
			spark.Start(config, forceFlag)
			stdout += fmt.Sprintf("Spark UI accessible through http://%s:31314\n", kube.PodPublicIP("spark-master"))
		}
		if allFlag || config.Cassandra {
			cassandra.Start(config, forceFlag)
			stdout += fmt.Sprintf("Cassandra accessible through %s:31317\n", kube.PodPublicIP("cassandra"))
		}
		fmt.Println(kube.GetPods())
		fmt.Print(stdout)
		return stdout + "\n" + kube.GetPods()
	}
	fmt.Println("Cluster is not running, run bdp start first")
	return "Cluster is not running, run bdp start first"
}
