package tlog

import (
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
	clssdk "github.com/tencentcloud/tencentcloud-cls-sdk-go"
)

type tlog struct {
	accessKeyID, accessKeySecret string
	topic                        string
	opt                          Options

	producer *clssdk.AsyncProducerClient
}

func New(id, secret string, topic string, opts ...Option) *tlog {
	options := Options{
		// 默认入口
		endpoint: "na-siliconvalley.cls.tencentcs.com",
	}
	for _, opt := range opts {
		opt(&options)
	}
	return &tlog{
		accessKeyID:     id,
		accessKeySecret: secret,
		topic:           topic,
		opt:             options,
	}
}

func (t *tlog) Establish() error {
	producerConfig := clssdk.GetDefaultAsyncProducerClientConfig()
	producerConfig.Endpoint = t.opt.endpoint
	producerConfig.AccessKeyID = t.accessKeyID
	producerConfig.AccessKeySecret = t.accessKeySecret
	producerInstance, err := clssdk.NewAsyncProducerClient(producerConfig)
	if err != nil {
		return err
	}
	producerInstance.Start()
	t.producer = producerInstance
	return nil
}

func (t tlog) Fire(entry *logrus.Entry) error {
	return t.producer.SendLog(
		t.topic,
		clssdk.NewCLSLog(
			time.Now().Unix(),
			getContent(entry),
		), nil,
	)
}

func (t tlog) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
	}
}

func getContent(entry *logrus.Entry) map[string]string {
	var out = map[string]string{
		"Level":   entry.Level.String(),
		"Message": entry.Message,
		"Time":    entry.Time.Format("2006-01-02 15:04:05"),
	}

	d, err := json.Marshal(entry.Data)
	if err == nil {
		out["Fields"] = string(d)
	}

	v, err := json.Marshal(entry.Caller)
	if err == nil {
		out["Caller"] = string(v)
	}
	return out
}
