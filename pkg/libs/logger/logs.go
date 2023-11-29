package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/astaxie/beego/logs"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendgridConf struct {
	APIKey   string
	Endpoint string
	Host     string
	From     string
	To       string
}

// Agent logs agent
var Agent *logs.BeeLogger
var sendgridConf = &SendgridConf{
	From: "select-core-team@liveramp.com",
	To:   "david.chen@liveramp.com",
}

type FileLogConf struct {
	Filename string `json:"filename"`
}

func init() {
	Initialize("console", `{"color":true}`, logs.LevelDebug, "{}")
	// Initialize("file", `{"filename":"/Users/hading/Workspace/New_SafeHeaven/ae-copilot/tmp/logs.log","daily":true,"maxdays":14,"level":7}`, logs.LevelDebug, "{}")
}

func Initialize(logType, logConf string, logLevel int, sgConf string) {
	logs.Info("Start to initialize logger with type: %s, conf: %s, level: %d, sendgridConf: %v.", logType, logConf, logLevel, sgConf)
	Agent = NewLogger(logType, logConf, logLevel)
	if err := json.Unmarshal([]byte(sgConf), &sendgridConf); err != nil {
		logs.Warning("Failed to initialize sendgrid config, error: %v.", err)
	}
}

func UpdateSendgridConf(sgConf string) {
	if err := json.Unmarshal([]byte(sgConf), &sendgridConf); err != nil {
		logs.Warning("Failed to update sendgrid config, error: %v.", err)
	}
}

func NewLogger(logType, logConf string, logLevel int) *logs.BeeLogger {
	var lt string
	// console、file、conn、es
	switch logType {
	case logs.AdapterFile:
		lt = logs.AdapterFile
		flc := new(FileLogConf)
		if err := json.Unmarshal([]byte(logConf), flc); err == nil {
			os.MkdirAll(path.Dir(flc.Filename), 0750)
		}
	case logs.AdapterConn:
		lt = logs.AdapterConn
	case logs.AdapterEs:
		lt = logs.AdapterEs
	default:
		lt = logs.AdapterConsole
	}
	logs.SetLevel(logLevel)
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	if err := logs.SetLogger(lt, logConf); err != nil {
		logs.Warning("Failed to reset logger, error: %v.", err)
		return logs.GetBeeLogger()
	}
	logs.Async(1e3)

	return logs.GetBeeLogger()
}

func SetLevel(level int) {
	logs.SetLevel(level)
}

func PrintStack(name string, r interface{}) {
	NoticeIssueToDevelopTeam(fmt.Sprintf("%s recovered panic error: %v, details: %s", name, r, RuntimeStack()))
}

func RuntimeStack() string {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return string(buf[:n])
}

func NoticeIssueToDevelopTeam(msg string) {
	logs.Error("[Exception Handling]: %s", msg)
}

func GenerateExceptionMsg(format string, a ...interface{}) string {
	return fmt.Sprintf("[Exception Handling]: %s", fmt.Sprintf(format, a))
}

func NoticeIssueViaEmail(subject, body string) {
	logs.Warning("[Select Core Alerts]: [%s] %s", subject, body)
	if sendgridConf.APIKey == "" {
		logs.Warning("Invalid API key of Sendgrid")
		return
	}
	from := mail.NewEmail("LiveRamp Safe Haven", sendgridConf.From)
	to := mail.NewEmail("", sendgridConf.To)
	message := mail.NewSingleEmail(from, subject, to, body, body)
	client := sendgrid.NewSendClient(sendgridConf.APIKey)
	if _, err := client.Send(message); err != nil {
		logs.Warning("Failed to send email via Sendgrid API, error: %v.", err)
	}
}
