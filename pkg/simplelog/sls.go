package simplelog

import (
	"fmt"
	"strconv"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
)

func NewSLS(optFunc ...OptionFunc) SLSLogger {
	opt := &slsOption{}
	for _, f := range optFunc {
		f(opt)
	}
	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.Endpoint = opt.Endpoint
	producerConfig.AccessKeyID = opt.AccessKeyId
	producerConfig.AccessKeySecret = opt.AccessKeySecret
	pdr := producer.InitProducer(producerConfig)
	pdr.Start()
	return &slsLogger{
		Producer: pdr,
		project:  opt.Project,
		logStore: opt.LogStore,
		topic:    opt.Topic,
		nodeName: opt.NodeName,
		nodeIp:   opt.NodeIp,
		podName:  opt.PodName,
		podIp:    opt.PodIp,
	}
}

type slsLogger struct {
	*producer.Producer
	project  string
	logStore string
	topic    string
	nodeName string
	nodeIp   string
	podName  string
	podIp    string
}

type SLSLogger interface {
	Logger
	SendLog(project, logstore, topic, source string, log *sls.Log) error
}

func NewSLSLoggerWithTopic(endpoint, accessKeyID, accessKeySecret,
	project, logstore, topic string) SLSLogger {
	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.Endpoint = endpoint
	producerConfig.AccessKeyID = accessKeyID
	producerConfig.AccessKeySecret = accessKeySecret
	pdr := producer.InitProducer(producerConfig)
	pdr.Start()
	return &slsLogger{
		Producer: pdr,
		project:  project,
		logStore: logstore,
		topic:    topic,
	}
}

func NewSLSLogger(endpoint, accessKeyID, accessKeySecret, project, logstore string) SLSLogger {
	return NewSLSLoggerWithTopic(endpoint, accessKeyID, accessKeySecret, project, logstore, "")
}

func (sl *slsLogger) sendLog(level LoggerLevel, code int, msg, stack string) {
	slsLog := producer.GenerateLog(uint32(time.Now().Unix()), map[string]string{
		"Level":    strconv.Itoa(int(level)),
		"Msg":      msg,
		"Code":     strconv.Itoa(code),
		"NodeName": sl.nodeName,
		"NodeIp":   sl.nodeIp,
		"PodName":  sl.podName,
		"PodIp":    sl.podIp,
		"Stack":    stack,
	})
	err := sl.SendLog(sl.project, sl.logStore, sl.topic, "", slsLog)
	if err != nil {
		fmt.Println("Send Log Error:", err)
	}
}

func (sl *slsLogger) Debug(v ...interface{}) {
	sl.sendLog(DebugLv, 0, fmt.Sprint(v...), "")
}

func (sl *slsLogger) Debugf(format string, v ...interface{}) {
	sl.sendLog(DebugLv, 0, fmt.Sprintf(format, v...), "")
}

func (sl *slsLogger) Info(v ...interface{}) {
	sl.sendLog(InfoLv, 0, fmt.Sprint(v...), "")
}

func (sl *slsLogger) Infof(format string, v ...interface{}) {
	sl.sendLog(InfoLv, 0, fmt.Sprintf(format, v...), "")
}

func (sl *slsLogger) Warn(v ...interface{}) {
	sl.sendLog(WarnLv, 0, fmt.Sprint(v...), "")
}

func (sl *slsLogger) Warnf(format string, v ...interface{}) {
	sl.sendLog(WarnLv, 0, fmt.Sprintf(format, v...), "")
}

func (sl *slsLogger) Error(v ...interface{}) {
	sl.sendLog(ErrorLv, 0, fmt.Sprint(v...), "")
}

func (sl *slsLogger) Errorf(format string, v ...interface{}) {
	sl.sendLog(ErrorLv, 0, fmt.Sprintf(format, v...), "")
}
