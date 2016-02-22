package ambari

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kube"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"util"

	"github.com/jmoiron/jsonq"
	"github.com/spf13/viper"
)

func CleanUp() {
	log.Println("Ambari: cleaning up cluster...")

	kube.DeleteResource("pods", "-l name=amb-slave")
	//kube.DeleteResource("rc", "amb-slave-controller")
	kube.DeleteResource("svc", "ambari")
	kube.DeleteResource("svc", "consul")
	kube.DeleteResource("svc", "namenode")
	kube.DeleteResource("pods", "amb-server.service.consul")
	kube.DeleteResource("pods", "amb-consul.service.consul")
	kube.DeleteResource("pods", "amb-shell")

	for {
		remaining := kube.RemainingPods("amb")
		if remaining == 0 {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}
	log.Println("Ambari: clean up done")
	util.ReleasePID("ambari")
}

func UpdateHosts(slavePods []string) {
	for v := 0; v < len(slavePods); v++ {
		if slavePods[v] != "" {
			hosts := ""
			for k := 0; k < len(slavePods); k++ {
				if v != k && slavePods[k] != "" {
					hosts += fmt.Sprintf("%s %s\\n", kube.PodIP(slavePods[k]), slavePods[k])
				}
			}
			cmd := fmt.Sprintf("\"echo $'%s' >> /etc/hosts\"", hosts)
			kube.ExecOnPod(slavePods[v], cmd)
		}
	}
}

func createSlaves(configFile string, component string) {
	slaves := viper.GetInt("ambari.AMBARI_NODES")
	for v := 1; v <= slaves; v++ {
		slave := util.Slave{fmt.Sprintf("amb-slave%d", v)}
		util.GenerateConfig(configFile, "ambari", slave)
		kube.CreateResource(filepath.Join(viper.GetString("BDP_CONFIG_DIR"), "tmp", configFile))
	}
}

func GetNamenode() string {
	url := fmt.Sprintf("http://%s:%s/api/v1/clusters/%s/services/HDFS/components/NAMENODE", kube.PodPublicIP("amb-server.service.consul"), "31313", viper.GetString("AMBARI_BLUEPRINT"))
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("GetNamenode: Error creating http request: %s \n", err)
	}
	request.SetBasicAuth("admin", "admin")
	client := &http.Client{}
	resp, err := client.Do(request)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("GetNamenode: Error performing http request: %s \n", err)
	}
	jsonObj := make(map[string]interface{})
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &jsonObj)
	if err != nil {
		log.Fatalf("GetNamenode: Error parsing json response: %s \n", err)
	}

	jq := jsonq.NewQuery(jsonObj)
	hostname, err := jq.String("host_components", "0", "HostRoles", "host_name")
	if err != nil {
		log.Fatalf("GetNamenode: could not find namenode hostname: %s \n", string(body))
	}

	log.Printf("Ambari: namenode running on pod: %s \n", hostname)
	return hostname
}

func Start(config util.Config, forceDeploy bool) {
	if !forceDeploy {
		if util.IsRunning("ambari") {
			log.Println("Ambari: already running, skipping start ...")
			return
		}
	}
	CleanUp()

	log.Println("Ambari: Launching consul")
	kube.CreateResource(viper.GetString("BDP_CONFIG_DIR") + "/ambari/consul.json")
	kube.CreateResource(viper.GetString("BDP_CONFIG_DIR") + "/ambari/consul-service.json")

	log.Println("Ambari: Waiting for consul server to start...")
	for {
		serverState := kube.PodStatus("amb-consul.service.consul")
		if serverState == "Running" {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}

	log.Println("Ambari: Launching Ambari server")
	kube.CreateResource(viper.GetString("BDP_CONFIG_DIR") + "/ambari/ambari-hdfs.json")
	kube.CreateResource(viper.GetString("BDP_CONFIG_DIR") + "/ambari/ambari-service.json")

	log.Println("Ambari: Waiting for ambari server to start...")
	for {
		serverState := kube.PodStatus("amb-server.service.consul")
		if serverState == "Running" {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}

	time.Sleep(10 * time.Second)
	log.Println("Ambari: registering consul services")

	ambariServiceIP := kube.ServiceIP("ambari")
	cmd := fmt.Sprintf("'curl -X PUT -d \"{\\\"Node\\\": \\\"ambari-8080\\\",\\\"Address\\\": \\\"%s\\\",\\\"Service\\\": {\\\"Service\\\": \\\"ambari-8080\\\"}}\" http://$CONSUL_SERVICE_HOST:8500/v1/catalog/register'",
		ambariServiceIP)
	kube.ExecOnPod("amb-server.service.consul", cmd)
	cmd = fmt.Sprintf("'curl -X PUT -d \"{\\\"Node\\\": \\\"amb-server\\\",\\\"Address\\\": \\\"%s\\\",\\\"Service\\\": {\\\"Service\\\": \\\"amb-server\\\"}}\" http://$CONSUL_SERVICE_HOST:8500/v1/catalog/register'",
		ambariServiceIP)
	kube.ExecOnPod("amb-server.service.consul", cmd)

	log.Println("Ambari: launching ambari slaves...")

	createSlaves("amb-slave.json", "ambari")

	for {
		pending := kube.PendingPods()
		if pending == 0 {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}

	slavePods := kube.PodNames("amb-slave")
	for v := 0; v < len(slavePods); v++ {
		if slavePods[v] != "" {
			podname := strings.Split(slavePods[v], ".")[0]
			cmd = fmt.Sprintf("'curl -X PUT -d \"{\\\"Node\\\": \\\"%[1]s\\\",\\\"Address\\\": \\\"$(hostname -I)\\\",\\\"Service\\\": {\\\"Service\\\": \\\"%[1]s\\\"}}\" http://$CONSUL_SERVICE_HOST:8500/v1/catalog/register'", podname)
			kube.ExecOnPod(slavePods[v], cmd)
		}
	}
	UpdateHosts(slavePods)

	log.Printf("Ambari: creating ambari cluster using blueprint: %s", viper.GetString("AMBARI_BLUEPRINT"))
	util.GenerateConfig("ambari-shell.json", "ambari", config)
	kube.CreateResource(viper.GetString("BDP_CONFIG_DIR") + "/tmp/ambari-shell.json")

	log.Println("Ambari: waiting to expose namenode service")
	time.Sleep(15 * time.Second)
	kube.Expose("pod", GetNamenode(), "--port=8020", "--target-port=8020", "--name=namenode")
	util.SetPID("ambari")
}
