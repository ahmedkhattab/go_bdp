package main

import (
	"codcli/launcher"
	"flag"
	"fmt"
	"kube"
	"log"
	"os"
	"reflect"
	"spark"
	"util"

	"github.com/spf13/viper"
)

func main() {

	util.SetDefaultConfig()

	if len(os.Args) == 1 {
		fmt.Println("usage: cod <command> [<args>]")
		fmt.Println("Commands: ")
		fmt.Println("\tstart   starts the kubernetes cluster")
		fmt.Println("\tstop    stops the kubernetes cluster")
		fmt.Println("\trestart stops the current cluster and restarts a new one")
		fmt.Println("\treset   removes all deployed components")
		fmt.Println("\tinfo    displays the cluster info")
		fmt.Println("\tpods    lists all pods running on the cluster")
		fmt.Println("\tdeploy  deploys bdp components on a running cluster")
		return
	}

	deployCommand := flag.NewFlagSet("deploy", flag.ExitOnError)

	allFlag := deployCommand.Bool("all", false, "")
	confFlag := deployCommand.String("conf", "", "")
	clusterFlag := deployCommand.String("cluster", "", "")
	forceFlag := deployCommand.Bool("f", false, "force the deployment by removing any running instances of the components")

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
	case "pods":
		fmt.Println(kube.GetPods())
	case "info":
		fmt.Println(kube.ClusterInfo())
	case "test":
		test()
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}

	if deployCommand.Parsed() {
		if len(os.Args[2:]) == 0 {
			fmt.Println("usage: decap deploy [<args>]")
			fmt.Println("args: ")
			fmt.Println("  -conf    path to a toml config file")
			fmt.Println("  -f       force deployment (removes running components first)")
			fmt.Println("  -cluster the kubernetes context to use for deployment")
			fmt.Println("  -all     to deploy all components with default configuration")
			os.Exit(2)
		}
		if *confFlag != "" {
			log.Printf("Loading configuration file %s \n", *confFlag)
			viper.SetConfigFile(*confFlag)

			err := viper.ReadInConfig()
			if err != nil {
				log.Fatalf("Error loading config file: %s \n", err)
			}
			//fmt.Println(viper.AllSettings())
		}

		if *clusterFlag != "" {
			kube.SetContext(*clusterFlag)
		}
		config := util.ConfigStruct()
		//fmt.Println(config)
		statuses := launcher.LaunchComponents(*allFlag, *forceFlag, config)
		v := reflect.ValueOf(statuses)
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).FieldByName("State").Bool() {
				fmt.Print(v.Field(i).FieldByName("Message"))
				fmt.Print(v.Field(i).FieldByName("URL"))
				fmt.Println("")
			}
		}
	}
}

func test() {
	spark.CleanUp()
}