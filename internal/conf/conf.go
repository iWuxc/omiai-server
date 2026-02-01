package conf

import (
	"fmt"
	"net/url"
	"omiai-server/pkg/config"
	"omiai-server/pkg/config/nacos"
	"omiai-server/pkg/config/source/file"
	"omiai-server/pkg/track"
	"os"
	"strings"
	"time"

	"github.com/iWuxc/go-wit/errors"
	"github.com/iWuxc/go-wit/log"
	"github.com/iWuxc/go-wit/utils"
	"github.com/spf13/cast"
)

var globalConfig = new(Config)

// Config Global conf .
type Config struct {
	Domain   *Domain           `json:"domain,omitempty"`
	Debug    bool              `json:"debug,omitempty"`
	Cron     bool              `json:"cron,omitempty"`
	Env      string            `json:"env"`
	Log      *Logger           `json:"log,omitempty"`
	Server   *Server           `json:"server,omitempty"`
	Database *Database         `json:"database,omitempty"`
	Cache    *Cache            `json:"cache"`
	Runtime  *Runtime          `json:"runtime"`
	Redis    *Redis            `json:"redis"`
	Queue    *Queue            `json:"queue"`
	Track    track.ManagerConf `json:"track"`
	Storage  *Storage          `json:"storage"`
	CronConf *Cron             `json:"cron_conf" mapstructure:"cron_conf"`
}

type Storage struct {
	Driver string `json:"driver"` // local, oss, cos
	OSS    OSS    `json:"oss"`
	COS    COS    `json:"cos"`
}

type OSS struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"access_key_id" mapstructure:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret" mapstructure:"access_key_secret"`
	BucketName      string `json:"bucket_name" mapstructure:"bucket_name"`
	Domain          string `json:"domain"`
}

type COS struct {
	BucketURL string `json:"bucket_url" mapstructure:"bucket_url"`
	Region    string `json:"region"`
	SecretID  string `json:"secret_id" mapstructure:"secret_id"`
	SecretKey string `json:"secret_key" mapstructure:"secret_key"`
}

type Cron struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DB       int    `json:"db"`
	Password string `json:"password"`
	CronName string `json:"cron_name"`
}

type Queue struct {
	Retry       int   `json:"retry"`       // 重试次数
	Timeout     int64 `json:"timeout"`     // 超时时间秒
	Retention   int64 `json:"retention"`   // 成功后消息保存时间
	Concurrency int   `json:"concurrency"` // 并发数
}

// Server conf .
type Server struct {
	RPCPort     int64 `json:"rpc_port,omitempty" mapstructure:"rpc_port"`
	HTTPPort    int64 `json:"http_port,omitempty" mapstructure:"http_port"`
	MonitorPort int64 `json:"monitor_port,omitempty" mapstructure:"monitor_port"`
}

// Logger . 日志相关配置
type Logger struct {
	Path             string `json:"path,omitempty"`
	Level            string `json:"level,omitempty"`
	MaxAge           uint32 `json:"max_age,omitempty" mapstructure:"max_age"` // 小时为单位
	RotationTimeHour uint32 `json:"rotation_time_hour"`                       // 日志轮转时间
}

// Database . 数据库相关配置
type Database struct {
	Debug bool `json:"debug"`
	// Default db conn conf
	Default struct {
		Driver string `json:"driver"`
		Source string `json:"source"`
	} `json:"default"`

	// DBConf db conn pool conf.
	DBConf struct {
		MaxOpenConn     int           `json:"max_open_conn" mapstructure:"max_open_conn"`
		MaxIdleConn     int           `json:"max_idle_conn" mapstructure:"max_idle_conn"`
		ConnMaxLifeTime time.Duration `json:"conn_max_life_time" mapstructure:"conn_max_life_time"`
	} `json:"db_conf" mapstructure:"db_conf"`
}

// Cache . 缓存相关配置
type Cache struct {
	Driver string `json:"driver"`
	URL    string `json:"url"`
}

// OSS . 阿里云 OSS 相关配置

type Domain struct {
	H5 string `json:"h5,omitempty"`
}

// Runtime conf.
type Runtime struct {
	Path     string `json:"path"`
	Download string `json:"download"`
}

// NacosSource 配置中心 .
type NacosSource struct {
	Endpoint    string `json:"endpoint"`
	Port        uint64 `json:"port"`
	DataID      string `json:"data_id"`
	Group       string `json:"group"`
	NameSpaceID string `json:"namespace_id"`
	LogLevel    string `json:"log_level"`
	UserName    string `json:"user_name"`
	Password    string `json:"password"`
	AccessKey   string `json:"access_key"`
	SecretKey   string `json:"secret_key"`
}

type Redis struct {
	Default struct {
		Url      string `json:"url"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Password string `json:"password"`
	} `json:"default"`
	Download struct {
		Url string `json:"url"`
	} `json:"download"`
}

type Nacos struct {
	Addr         string `json:"addr,omitempty"`
	NamespaceId  string `json:"namespace_id,omitempty"`
	Port         uint64 `json:"port,omitempty"`
	Name         string `json:"name,omitempty"`
	Timeout      uint64 `json:"timeout,omitempty" mapstructure:"timeout"`
	BeatInterval int64  `json:"beat_interval,omitempty" mapstructure:"beat_interval"`
	Log          Logger `json:"log,omitempty" mapstructure:"log"`
	UserName     string `json:"user_name"`
	Password     string `json:"password"`
}

func Init(conf, confCenter string) (f func(), e error) {
	// 加载配置文件
	s := source(conf, confCenter)
	c := config.New(config.WithSource(
		file.NewSource(strings.TrimRight(conf, "/")+"/config.yaml"),
		nacos.NewSource(
			nacos.WithEndpoint(s.Endpoint),
			nacos.WithPort(s.Port),
			nacos.WithDataID(s.DataID),
			nacos.WithGroup(s.Group),
			nacos.WithNamespaceID(s.NameSpaceID),
			nacos.WithLogLevel(s.LogLevel),
			nacos.WithUserName(s.UserName),
			nacos.WithPassword(s.Password),
			nacos.WithAccessKey(s.AccessKey),
			nacos.WithSecretKey(s.SecretKey),
		),
	))

	if e = c.Load(); e != nil {
		return nil, errors.Wrap(e, "load config failed")
	}

	if e = c.Scan(globalConfig); e != nil {
		return nil, errors.Wrap(e, "scan config failed")
	}

	if e = logInit(globalConfig.Log); e != nil {
		return
	}

	if e = cacheInit(globalConfig.Cache); e != nil {
		return
	}
	// 监听配置变动
	// watch(c, "cron")
	// watch(c, "llm")
	// watch(c, "sign_api_key")
	// watch(c, "private_deploy")
	// wathc, "fc_test")
	// watch(c, "meitu")
	return
}

func GetConfig() *Config {
	return globalConfig
}

// source .
func source(conf, confCenter string) NacosSource {
	var s NacosSource
	//启动命令
	if u, err := url.Parse(confCenter); len(confCenter) > 0 && err == nil {
		q := u.Query()
		s.Endpoint = u.Host
		s.Port = cast.ToUint64(q.Get("port"))
		s.DataID = q.Get("data_id")
		s.Group = q.Get("group")
		s.NameSpaceID = q.Get("namespace_id")
		s.LogLevel = q.Get("log_level")
		s.UserName = q.Get("user_name")
		s.Password = q.Get("password")
		s.AccessKey = q.Get("access_key")
		s.SecretKey = q.Get("secret_key")
		fmt.Printf("使用环启动命令nacos 连接地址: %s", confCenter)
		return s
	}

	// 环境变量
	nacosUrl := os.Getenv("nacos_url")
	if nacosUrl != "" {
		if u, err := url.Parse(nacosUrl); len(nacosUrl) > 0 && err == nil {
			q := u.Query()
			s.Endpoint = u.Host
			s.Port = cast.ToUint64(q.Get("port"))
			s.DataID = q.Get("data_id")
			s.Group = q.Get("group")
			s.NameSpaceID = q.Get("namespace_id")
			s.LogLevel = q.Get("log_level")
			s.UserName = q.Get("user_name")
			s.Password = q.Get("password")
			s.AccessKey = q.Get("access_key")
			s.SecretKey = q.Get("secret_key")
			fmt.Printf("使用环境变量连接nacos 连接地址: %s", nacosUrl)
			return s
		}
	}

	//本地
	f := strings.TrimRight(conf, "/") + "/nacos.yaml"
	if !utils.FileExist(f) {
		return s
	}

	c := config.New(config.WithSource(file.NewSource(f)))
	if err := c.Load(); err != nil {
		panic(errors.Wrap(err, "load nacos config failed"))
	}

	if err := c.Scan(&s); err != nil {
		panic(errors.Wrap(err, "scan nacos config failed"))
	}
	fmt.Printf("使用本地nacos.yaml连接nacos 连接")
	return s
}

// watch .
func watch(c config.ConfigInterface, key string) {
	if err := c.Watch(key, func(key string, value config.Value) {
		log.Printf("config(key=%s) changed: %s\n", key, value.Load())
		if err := c.Scan(globalConfig); err != nil {
			log.Error(err)
		}
	}); err != nil {
		panic(err)
	}
}

// value for globalConfig
func value(c config.ConfigInterface, key string) {
	fmt.Printf("=========== value result =============\n")
	v := c.Value(key).Load()
	fmt.Printf("key=%s, load: %+v\n\n", key, v)
}
