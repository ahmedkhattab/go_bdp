package main

import (
	"ambari"
	"cassandra"
	"flag"
	"fmt"
	"kube"
	"log"
	"os"
	"os/exec"
	"rabbitmq"
	"spark"
	"util"

	"github.com/spf13/viper"
)

func main() {
	os.Setenv("BDP_CONFIG_DIR", "/home/khattab/BDP")
	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatalf("Error loading config file: %s \n", err)
	}

	setEnvironment()

	if len(os.Args) == 1 {
		fmt.Println("usage: bdp <command> [<args>]")
		fmt.Println("Commands: ")
		fmt.Println("\tstart   starts the cluster")
		fmt.Println("\tstop    stops the cluster")
		fmt.Println("\treset   removes all deployed components")
		fmt.Println("\tdeploy  deploys bdp components on a running cluster")
		return
	}
	deployCommand := flag.NewFlagSet("deploy", flag.ExitOnError)
	cassandraFlag := deployCommand.Bool("cassandra", false, "")
	rabbitmqFlag := deployCommand.Bool("rabbitmq", false, "")
	ambariFlag := deployCommand.Bool("ambari", false, "")
	sparkFlag := deployCommand.Bool("spark", false, "")
	allFlag := deployCommand.Bool("all", false, "")

	switch os.Args[1] {
	case "start":
		startCluster()
	case "stop":
		stopCluster()
	case "deploy":
		deployCommand.Parse(os.Args[2:])
	case "reset":
		cleanUpCluster()
	case "test":
		test()
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}

	if deployCommand.Parsed() {
		if len(os.Args[2:]) == 0 {
			*allFlag = true
		}
		stdout := ""
		if kube.ClusterIsUp() {
			if *ambariFlag || *allFlag {
				ambari.Start()
				stdout += fmt.Sprintf("Ambari UI accessible through http://%s:31313\n", kube.PodPublicIP("amb-server.service.consul"))
			}
			if *sparkFlag || *allFlag {
				spark.Start()
				stdout += fmt.Sprintf("Spark UI accessible through http://%s:31314\n", kube.PodPublicIP("spark-master"))
			}
			if *cassandraFlag || *allFlag {
				cassandra.Start()
				stdout += fmt.Sprintf("Cassandra accessible through %s:31317\n", kube.PodPublicIP("spark-master"))
			}
			if *rabbitmqFlag || *allFlag {
				rabbitmq.Start()
				stdout += fmt.Sprintf("RabbitMQ UI accessible through http://%s:31316\n", kube.PodPublicIP("spark-master"))
			}
			fmt.Println(kube.GetPods())
			fmt.Print(stdout)
		} else {
			fmt.Println("Cluster is not running, run bdp start first")
		}

	}
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
		fmt.Println("Cluster is already running")
		return true
	}
	cmd := exec.Command("sh", "-c", viper.GetString("KUBE_PATH")+"/cluster/kube-up.sh")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func stopCluster() bool {
	if kube.ClusterIsUp() {
		cmd := exec.Command("sh", "-c", viper.GetString("KUBE_PATH")+"/cluster/kube-down.sh")
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
			return false
		}
		return true
	}
	fmt.Println("Cluster is already stopped")
	return true
}

func deployComponents() {
	if kube.ClusterIsUp() {
		//	ambari.Start()
		spark.Start()
		//	rabbitmq.Start()
		//cassandra.Start()

		//	fmt.Printf("Ambari UI accessible through http://%s:31313\n", kube.PodPublicIP("amb-server.service.consul"))
		fmt.Printf("Spark UI accessible through http://%s:31314\n", kube.PodPublicIP("spark-master"))
		//	fmt.Printf("RabbitMQ UI accessible through http://%s:31316\n", kube.PodPublicIP("spark-master"))
		//	fmt.Printf("Cassandra accessible through %s:31317\n", kube.PodPublicIP("spark-master"))

		fmt.Println(kube.GetPods())
	} else {
		fmt.Println("Cluster is not running, run bdp start first")
	}
}

func cleanUpCluster() {
	if kube.ClusterIsUp() {
		ambari.CleanUp()
		spark.CleanUp()
		rabbitmq.CleanUp()
		cassandra.CleanUp()
	}
}

func test() {
	rc := util.LoadRC("/home/khattab/BDP/Ambari/ambari-slave.json")
	util.SaveRC("amb.json", rc)

	rc = util.LoadRC("/home/khattab/BDP/spark/spark-worker-controller.json")
	util.SaveRC("spark.json", rc)

	rc = util.LoadRC("/home/khattab/BDP/rabbitmq/rabbitmq-controller.json")
	util.SaveRC("rabbitmq.json", rc)

	rc = util.LoadRC("/home/khattab/BDP/cassandra/cassandra-controller.json")
	util.SaveRC("cassandra.json", rc)
}
