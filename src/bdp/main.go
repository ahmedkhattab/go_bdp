package main

import (
	"ambari"
	"cassandra"
	"flag"
	"fmt"
	"io/ioutil"
	"kube"
	"log"
	"net/http"
	"os"
	"rabbitmq"
	"spark"
	"strings"
	"util"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatalf("Error loading config file: %s \n", err)
	}

	config := util.ConfigStruct()
	util.SetEnvVars()

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
		kube.StartCluster()
	case "stop":
		kube.StopCluster()
	case "deploy":
		deployCommand.Parse(os.Args[2:])
	case "reset":
		resetCluster()
	case "test":
		test(config)
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
				ambari.Start(config)
				stdout += fmt.Sprintf("Ambari UI accessible through http://%s:31313\n", kube.PodPublicIP("amb-server.service.consul"))
			}
			if *sparkFlag || *allFlag {
				spark.Start(config)
				stdout += fmt.Sprintf("Spark UI accessible through http://%s:31314\n", kube.PodPublicIP("spark-master"))
			}
			if *cassandraFlag || *allFlag {
				cassandra.Start(config)
				stdout += fmt.Sprintf("Cassandra accessible through %s:31317\n", kube.PodPublicIP("spark-master"))
			}
			if *rabbitmqFlag || *allFlag {
				rabbitmq.Start(config)
				stdout += fmt.Sprintf("RabbitMQ UI accessible through http://%s:31316\n", kube.PodPublicIP("spark-master"))
			}
			fmt.Println(kube.GetPods())
			fmt.Print(stdout)
		} else {
			fmt.Println("Cluster is not running, run bdp start first")
		}
	}
}

func resetCluster() {
	if kube.ClusterIsUp() {
		ambari.CleanUp()
		spark.CleanUp()
		rabbitmq.CleanUp()
		cassandra.CleanUp()
	}
}

func test(config util.Config) {
	//	curlAmbari("hosts")
	//fmt.Println(kube.PodPublicIP("amb-server.service.consul"))
	//	kube.DeleteResource("pod", kube.PodNames("amb-slave")[0])
	//	curlAmbari("hosts")
	slavePods := kube.PodNames("amb-slave")
	for v := 0; v < len(slavePods); v++ {
		if slavePods[v] != "" {
			podname := strings.Split(slavePods[v], ".")[0]
			cmd := fmt.Sprintf("'curl -X PUT -d \"{\\\"Node\\\": \\\"%s\\\",\\\"Address\\\": \\\"$(hostname -I)\\\",\\\"Service\\\": {\\\"Service\\\": \\\"%s\\\"}}\" http://$CONSUL_SERVICE_HOST:8500/v1/catalog/register'", podname)
			fmt.Println(cmd)
			kube.ExecOnPod(slavePods[v], cmd)
		}
	}
}

func curlAmbari(url string) {
	urlFull := fmt.Sprintf("http://%s:%s/api/v1/%s", kube.PodPublicIP("amb-server.service.consul"), "31313", url)
	request, err := http.NewRequest("GET", urlFull, nil)
	if err != nil {
		log.Fatalf("Error creating http request: %s \n", err)
	}
	request.SetBasicAuth("admin", "admin")
	client := &http.Client{}
	resp, err := client.Do(request)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("Error performing http request: %s \n", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
