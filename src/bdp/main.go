package main

import (
	"ambari"
	"cassandra"
	"flag"
	"fmt"
	"kube"
	"os"
	"rabbitmq"
	"spark"
)

func main() {
	os.Setenv("KUBE", "/home/khattab/kubernetes-1.1.2/cluster/kubectl.sh")
	/*
	  fmt.Println(kube.Get_remaining_pods("amb"))
	  fmt.Println(kube.Get_pending_pods())

	  fmt.Println(kube.Get_pod_ip("amb-server.service.consul"))
	  fmt.Println(kube.Get_pod_host_ip("amb-server.service.consul"))
	  fmt.Println(kube.Get_pod_status("amb-server.service.consul"))
	  fmt.Println(kube.Get_pod_public_ip("amb-server.service.consul"))
	*/
	var verbose bool
	flag.BoolVar(&verbose, "v", false, "")

	ambari.Start()
	spark.Start()
	rabbitmq.Start()
	cassandra.Start()

	fmt.Printf("Ambari UI accessible through http://%s:31313\n", kube.PodPublicIP("amb-server.service.consul"))
	fmt.Printf("Spark UI accessible through http://%s:31314\n", kube.PodPublicIP("spark-master"))
	fmt.Printf("RabbitMQ UI accessible through http://%s:31316\n", kube.PodPublicIP("spark-master"))
	fmt.Printf("Cassandra accessible through %s:31317\n", kube.PodPublicIP("spark-master"))

	fmt.Println(kube.GetPods())

}
