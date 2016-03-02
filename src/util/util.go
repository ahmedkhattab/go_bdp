package util

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
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
	Ambari             bool
	Spark              bool
	Kafka              bool
	Cassandra          bool
	Rabbitmq           bool
}

type Status struct {
	State   bool
	Message string
}

type Statuses struct {
	Ambari    Status
	Spark     Status
	Kafka     Status
	Cassandra Status
	Rabbitmq  Status
}
type Slave struct {
	AmbariSlaveName string
}

func InitConfigStruct() Config {
	return Config{AmbariNodes: viper.GetInt("ambari.AMBARI_NODES"),
		CassandraNodes:     viper.GetInt("cassandra.CASSANDRA_NODES"),
		RabbitmqNodes:      viper.GetInt("rabbitmq.RABBITMQ_NODES"),
		SparkWorkers:       viper.GetInt("spark.SPARK_WORKERS"),
		KafkaNodes:         viper.GetInt("kafka.KAFKA_NODES"),
		AmbariBlueprint:    viper.GetString("ambari.AMBARI_BLUEPRINT"),
		AmbariBlueprintURL: viper.GetString("ambari.AMBARI_BLUEPRINT_URL")}
}

func InitStatusesStruct() Statuses {
	return Statuses{Spark: Status{false, "Not Running"},
		Cassandra: Status{false, "Not Running"},
		Ambari:    Status{false, "Not Running"},
		Rabbitmq:  Status{false, "Not Running"},
		Kafka:     Status{false, "Not Running"}}
}

//ConfigStruct creates an instance of the config structure out of the config
//parameters to be used by the template engine to generate kubernetes resources
//config files
func ConfigStruct() Config {
	return Config{AmbariNodes: viper.GetInt("ambari.AMBARI_NODES"),
		CassandraNodes:     viper.GetInt("cassandra.CASSANDRA_NODES"),
		RabbitmqNodes:      viper.GetInt("rabbitmq.RABBITMQ_NODES"),
		SparkWorkers:       viper.GetInt("spark.SPARK_WORKERS"),
		KafkaNodes:         viper.GetInt("kafka.KAFKA_NODES"),
		AmbariBlueprint:    viper.GetString("ambari.AMBARI_BLUEPRINT"),
		AmbariBlueprintURL: viper.GetString("ambari.AMBARI_BLUEPRINT_URL"),
		Ambari:             viper.InConfig("ambari"),
		Kafka:              viper.InConfig("kafka"),
		Cassandra:          viper.InConfig("cassandra"),
		Rabbitmq:           viper.InConfig("rabbitmq"),
		Spark:              viper.InConfig("spark")}
}

func (config *Config) SetState(component string, value bool) {
	s := reflect.ValueOf(config).Elem()
	s.FieldByName(component).SetBool(value)
}

func SetDefaultConfig() {
	viper.SetConfigType("toml")
	viper.SetConfigName("defaults")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../src/bdp")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error loading config file: %s \n", err)
	}
	for key, value := range viper.AllSettings() {
		viper.SetDefault(key, value)
	}
	setEnvVars()
}

//SetEnvVars sets the environment variables needed for the kube-up script based
//on values provided in the yaml config file
func setEnvVars() {
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
	if err != nil {
		log.Fatalf("GenerateConfig: Error creating config file: %s \n", err)
	}
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
