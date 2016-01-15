package main

import (
	"fmt"
	"kube"
	"os"
	"os/exec"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")
	viper.SetConfigName("config")
	viper.AddConfigPath("$GOPATH/bin")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	setEnvironment()
	startCluster()

	fmt.Println(viper.GetString("KUBE_PATH"))

	/*ambari.Start()
	spark.Start()
	rabbitmq.Start()
	cassandra.Start()

	fmt.Printf("Ambari UI accessible through http://%s:31313\n", kube.PodPublicIP("amb-server.service.consul"))
	fmt.Printf("Spark UI accessible through http://%s:31314\n", kube.PodPublicIP("spark-master"))
	fmt.Printf("RabbitMQ UI accessible through http://%s:31316\n", kube.PodPublicIP("spark-master"))
	fmt.Printf("Cassandra accessible through %s:31317\n", kube.PodPublicIP("spark-master"))

	fmt.Println(kube.GetPods())
	*/
}

func setEnvironment() {
	os.Setenv("KUBERNETES_PROVIDER", viper.GetString("KUBERNETES_PROVIDER"))
	os.Setenv("KUBE_AWS_ZONE", viper.GetString("KUBE_AWS_ZONE"))
	os.Setenv("NUM_MINIONS", viper.GetString("NUM_MINIONS"))
	os.Setenv("MINION_SIZE", viper.GetString("MINION_SIZE"))
	os.Setenv("MASTER_SIZE", viper.GetString("MASTER_SIZE"))
	os.Setenv("AWS_S3_REGION", viper.GetString("AWS_S3_REGION"))
	os.Setenv("INSTANCE_PREFIX", viper.GetString("INSTANCE_PREFIX"))
	os.Setenv("AWS_S3_BUCKET", viper.GetString("AWS_S3_BUCKET"))
	os.Setenv("MINION_ROOT_DISK_SIZE", viper.GetString("MINION_ROOT_DISK_SIZE"))
	os.Setenv("MASTER_ROOT_DISK_SIZE", viper.GetString("MASTER_ROOT_DISK_SIZE"))
}

func startCluster() bool {
	if kube.ClusterIsUp() {
		return true
	} else {
		exec.Command("sh", "-c", viper.GetString("KUBE_PATH")+"cluster/kube-up.sh")
		return true
	}
}
