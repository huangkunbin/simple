package simplelog

import "github.com/aliyun/aliyun-log-go-sdk/producer"

type OptionFunc func(opt *slsOption)

type slsOption struct {
	*producer.Producer
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	Project         string
	LogStore        string
	Topic           string
	NodeName        string
	NodeIp          string
	PodName         string
	PodIp           string
}

func EndPoint(endPoint, ak, as string) OptionFunc {
	return func(opt *slsOption) {
		opt.Endpoint = endPoint
		opt.AccessKeyId = ak
		opt.AccessKeySecret = as
	}
}

func ProjectInfo(project, store string) OptionFunc {
	return func(opt *slsOption) {
		opt.Project = project
		opt.LogStore = store
	}
}

func Topic(topic string) OptionFunc {
	return func(opt *slsOption) {
		opt.Topic = topic
	}
}

func HostInfo(nodeName, nodeIp, podName, podIp string) OptionFunc {
	return func(opt *slsOption) {
		opt.NodeName = nodeName
		opt.NodeIp = nodeIp
		opt.PodName = podName
		opt.PodIp = podIp
	}
}
