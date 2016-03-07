package launcher

import (
	"ambari"
	"cassandra"
	"kafka"
	"kube"
	"log"
	"os"
	"path/filepath"
	"rabbitmq"
	"spark"
	"util"

	"github.com/spf13/viper"
)

func LaunchComponents(allFlag bool, forceFlag bool, config util.Config) util.Statuses {

	if kube.ClusterIsUp() {
		os.Mkdir(filepath.Join(viper.GetString("BDP_CONFIG_DIR"), "tmp"), 0777)
		if allFlag || config.Ambari {
			ambari.Start(config, forceFlag)
		}
		if allFlag || config.Rabbitmq {
			rabbitmq.Start(config, forceFlag)
		}
		if allFlag || config.Kafka {
			kafka.Start(config, forceFlag)
		}
		if allFlag || config.Spark {
			spark.Start(config, forceFlag)
		}
		if allFlag || config.Cassandra {
			cassandra.Start(config, forceFlag)
		}
	} else {
		log.Println("Cluster is not running, run bdp start first")
	}
	return ComponentsStatuses()
}

func ComponentsStatuses() util.Statuses {
	statuses := util.InitStatusesStruct()
	statuses.Ambari = ambari.Status()
	statuses.Rabbitmq = rabbitmq.Status()
	statuses.Kafka = kafka.Status()
	statuses.Spark = spark.Status()
	statuses.Cassandra = cassandra.Status()
	return statuses
}
