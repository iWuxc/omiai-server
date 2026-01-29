package trace

import (
	"github.com/sirupsen/logrus"
)

const (
	GroupQueue group = "QUEUE"
	GroupCron  group = "CRON"
	GroupHttp  group = "SERVER"
)

type group string

type logCtx struct {
	group
}

func (l *logCtx) Fire(entry *logrus.Entry) error {
	//if entry != nil && entry.Context != nil && !reflect.ValueOf(entry.Context).IsNil() && entry.Context.Value("request_id") != nil {
	//	entry.Data["requestID"] = entry.Context.Value("request_id").(string)
	//}
	if entry != nil && entry.Context != nil && entry.Context.Value("request_id") != nil {
		entry.Data["requestID"] = entry.Context.Value("request_id").(string)
	}
	if entry != nil && entry.Context != nil && entry.Context.Value("api-key") != nil {
		entry.Data["apiKey"] = entry.Context.Value("api-key").(string)
	}

	if entry != nil && entry.Context != nil && entry.Context.Value("request_url") != nil {
		entry.Data["requestUrl"] = entry.Context.Value("request_url").(string)
	}
	entry.Data["group"] = l.group
	return nil
}

func (l *logCtx) Levels() []logrus.Level {
	return logrus.AllLevels
}

func NewLogCtx(group group) *logCtx {
	return &logCtx{
		group: group,
	}
}
