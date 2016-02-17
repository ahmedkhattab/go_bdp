package main

import (
	"ambari"
	"cassandra"
	"flag"
	"fmt"
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

func main() {

	components := []string{"ambari", "cassandra", "rabbitmq", "spark", "kafka"}
	viper.SetConfigType("toml")
	viper.SetConfigName("defaults")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../src/bdp")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error loading config file: %s \n", err)
	}
	util.SetEnvVars()

	if len(os.Args) == 1 {
		fmt.Println("usage: bdp <command> [<args>]")
		fmt.Println("Commands: ")
		fmt.Println("\tstart   starts the cluster")
		fmt.Println("\tstop    stops the cluster")
		fmt.Println("\trestart stops the current cluster and restarts a new one")
		fmt.Println("\treset   removes all deployed components")
		fmt.Println("\tinfo    lists all pods running on the cluster")
		fmt.Println("\tdeploy  deploys bdp components on a running cluster")
		return
	}

	deployCommand := flag.NewFlagSet("deploy", flag.ExitOnError)
	componentsMap := make(map[string]*bool)
	for _, component := range components {
		componentsMap[component] = deployCommand.Bool(component, false, "")
	}
	allFlag := deployCommand.Bool("all", false, "")
	confFlag := deployCommand.String("conf", "", "path to the json config file")
	clusterFlag := deployCommand.String("cluster", "", "the kubernetes context to use for deployment")

	switch os.Args[1] {
	case "start":
		kube.StartCluster()
	case "stop":
		kube.StopCluster()
	case "restart":
		kube.StopCluster()
		kube.StartCluster()
	case "deploy":
		deployCommand.Parse(os.Args[2:])
	case "reset":
		kube.ResetCluster()
	case "info":
		fmt.Println(kube.GetPods())
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
		if *confFlag != "" {
			log.Printf("Loading configuration file %s \n", *confFlag)
			viper.SetConfigFile(*confFlag)

			err := viper.MergeInConfig()
			if err != nil {
				log.Fatalf("Error loading config file: %s \n", err)
			}
			for _, component := range components {
				if viper.IsSet(component) {
					*componentsMap[component] = true
				}
			}
		}

		if *clusterFlag != "" {
			kube.SetContext(*clusterFlag)
		}

		config := util.ConfigStruct()
		fmt.Println(config)
		launchComponents(componentsMap, allFlag, config)
	}
}

func launchComponents(componentsMap map[string]*bool, allFlag *bool, config util.Config) {
	stdout := ""
	if kube.ClusterIsUp() {
		os.Mkdir(filepath.Join(viper.GetString("BDP_CONFIG_DIR"), "tmp"), 0777)
		if *componentsMap["ambari"] || *allFlag {
			ambari.Start(config)
			stdout += fmt.Sprintf("Ambari UI accessible through http://%s:31313\n", kube.PodPublicIP("amb-server.service.consul"))
		}
		if *componentsMap["rabbitmq"] || *allFlag {
			rabbitmq.Start(config)
			stdout += fmt.Sprintf("RabbitMQ UI accessible through http://%s:31316\n", kube.PodPublicIP("spark-master"))
		}
		if *componentsMap["kafka"] || *allFlag {
			kafka.Start(config)
			stdout += fmt.Sprintf("Kafka accessible through http://%s:31318\n", kube.PodPublicIP("spark-master"))
		}
		if *componentsMap["spark"] || *allFlag {
			spark.Start(config)
			stdout += fmt.Sprintf("Spark UI accessible through http://%s:31314\n", kube.PodPublicIP("spark-master"))
		}
		if *componentsMap["cassandra"] || *allFlag {
			cassandra.Start(config)
			stdout += fmt.Sprintf("Cassandra accessible through %s:31317\n", kube.PodPublicIP("spark-master"))
		}
		fmt.Println(kube.GetPods())
		fmt.Print(stdout)
	} else {
		fmt.Println("Cluster is not running, run bdp start first")
	}
}

func test() {
	spark.CleanUp()
}
