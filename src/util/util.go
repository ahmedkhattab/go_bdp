package util

import (
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/viper"
)

type Config struct {
	AmbariNodes        int
	CassandraNodes     int
	RabbitmqNodes      int
	SparkWorkers       int
	KafkaNodes         int
	AmbariBlueprint    string
	AmbariBlueprintURL string
}

type Slave struct {
	AmbariSlaveName string
}

//ConfigStruct creates an instance of the config structure out of the config
//parameters to be used by the template engine to generate kubernetes resources
//config files
func ConfigStruct() Config {
	return Config{viper.GetInt("ambari.AMBARI_NODES"),
		viper.GetInt("cassandra.CASSANDRA_NODES"),
		viper.GetInt("rabbitmq.RABBITMQ_NODES"),
		viper.GetInt("spark.SPARK_WORKERS"),
		viper.GetInt("kafka.KAFKA_NODES"),
		viper.GetString("ambari.AMBARI_BLUEPRINT"),
		viper.GetString("ambari.AMBARI_BLUEPRINT_URL")}
}

//SetEnvVars sets the environment variables needed for the kube-up script based
//on values provided in the yaml config file
func SetEnvVars() {
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

//GenerateConfig generates a kubernetes configuration instance based on the
//input template and the input configuration data.
//Assumes a directory structure: $BDP_CONFIG_DIR/componentName/filename
//e.g. $BDP_CONFIG_DIR/Ambari/ambari-slave.json
func GenerateConfig(templateName string, component string, data interface{}) {
	tmpl, err := template.ParseFiles(filepath.Join(viper.GetString("BDP_CONFIG_DIR"), component, templateName))
	if err != nil {
		log.Fatalf("GenerateConfig: Error parsing templates: %s \n", err)
	}

	outputFile, err := os.Create(filepath.Join(viper.GetString("BDP_CONFIG_DIR"), "tmp", templateName))
	err = tmpl.ExecuteTemplate(outputFile, templateName, data)
	if err != nil {
		log.Fatalf("GenerateConfig: Error generating configuration from template: %s \n", err)
	}
}

func SetPID(pName string) {
	os.Create(filepath.Join(viper.GetString("BDP_CONFIG_DIR"), "tmp", pName+".pid"))
}

func ReleasePID(pName string) {
	os.Remove(filepath.Join(viper.GetString("BDP_CONFIG_DIR"), "tmp", pName+".pid"))
}

func IsRunning(pName string) bool {
	_, err := os.Stat(filepath.Join(viper.GetString("BDP_CONFIG_DIR"), "tmp", pName+".pid"))
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
