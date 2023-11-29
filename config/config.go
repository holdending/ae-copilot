package config

import (
	"os"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

var (
	Agent *configData
)

type configI interface {
	defaultString(string, string) string
	defaultBool(string, bool) bool
	defaultFloat(string, float64) float64
	defaultInt(string, int) int
}

type localConfig struct{}

func (*localConfig) defaultString(key, defaultValue string) string {
	return beego.AppConfig.DefaultString(key, defaultValue)
}

func (*localConfig) defaultBool(key string, defaultValue bool) bool {
	return beego.AppConfig.DefaultBool(key, defaultValue)
}

func (*localConfig) defaultFloat(key string, defaultValue float64) float64 {
	return beego.AppConfig.DefaultFloat(key, defaultValue)
}

func (*localConfig) defaultInt(key string, defaultValue int) int {
	return beego.AppConfig.DefaultInt(key, defaultValue)
}

type configData struct {
	AppName  string
	HTTPPort string

	LogType      string
	LogConf      string
	LogLevel     int
	SendgridConf string

	ScanIntervalTime int

	InPath     string
	RejectPath string

	GCSCredentials string
	Tenants        []string
}

func init() {
	var config configI = new(localConfig)
	Agent = &configData{}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = ""
	}

	Agent.AppName = config.defaultString("appname", hostname)
	Agent.HTTPPort = config.defaultString("http.port", "8080")

	Agent.LogType = config.defaultString("log.type", "console")
	Agent.LogConf = config.defaultString("log.config", `{"filename":"/Users/hading/Workspace/New_SafeHeaven/ae-copilot/tmp/logs.log","daily":true,"maxdays":14,"level":6}`)
	Agent.LogLevel = config.defaultInt("log.level", logs.LevelDebug)
	Agent.SendgridConf = config.defaultString("sendgrid.conf", `{"From":"select-core-team@liveramp.com","To":"david.chen@liveramp.com"}`)

	Agent.ScanIntervalTime = config.defaultInt("scan.interval.time.seconds", 10) // Seconds

	// Agent.GCSCredentials = config.defaultString("gcs.credentials", `{"ProjectID":"datalake-landing-eng-us-prod"}`)
	// Agent.RejectPath = "gs://lranalytics-au-endpoint-select-vm/%s/REJECT/"
	Agent.GCSCredentials = config.defaultString("gcs.credentials", `{"ProjectID":"select-eng-us-2pqa"}`)
	Agent.RejectPath = "gs://lr-select-vm-us-qa-temp/%s/%s"

	Agent.Tenants = strings.Split(config.defaultString("tenants", "721211"), ",")

	Agent.InPath = "%s/%s/%s"

}
