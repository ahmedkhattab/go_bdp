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
	AmbariBlueprint    string
	AmbariBlueprintURL string
}

type Slave struct {
	AmbariSlaveName string
}

func ConfigStruct() Config {
	return Config{viper.GetInt("AMBARI_NODES"),
		viper.GetInt("CASSANDRA_NODES"),
		viper.GetInt("RABBITMQ_NODES"),
		viper.GetInt("SPARK_WORKERS"),
		viper.GetString("AMBARI_BLUEPRINT"),
		viper.GetString("AMBARI_BLUEPRINT_URL")}
}

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

//GenerateConfig generates a kubernetes configuration instance based on the input template
//and the input configuration data.
//assumed directory structure: $BDP_CONFIG_DIR/componentName/filename
//e.g. $BDP_CONFIG_DIR/Ambari/ambari-slave.json
func GenerateConfig(templateName string, component string, data Config) {

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
