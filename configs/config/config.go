package config

import (
	"flag"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

var (
	Conf *Config
)

type Config struct {
	Db     ConfDB     `mapstructure:"database"`
	Web    ConfWeb    `mapstructure:"web"`
	Logger ConfLogger `mapstructure:"logger"`
}

type ConfDB struct {
	Postgres ConfPostgres `mapstructure:"postgres"`
}

type ConfPostgres struct {
	DriverName string `mapstructure:"driver_name"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	DbName     string `mapstructure:"db_name"`
	SslMode    string `mapstructure:"ssl_mode"`
	Host       string `mapstructure:"host"`
	MaxConn    string `mapstructure:"max_conn"`
}

type ConfWeb struct {
	Server ConfServer `mapstructure:"server"`
}

type ConfServer struct {
	Address  string `mapstructure:"address"`
	Port     string `mapstructure:"port"`
	Host     string `mapstructure:"host"`
	Protocol string `mapstructure:"protocol"`
}

type ConfApi struct {
	Address string `mapstructure:"address"`
	Port    string `mapstructure:"port"`
	Host    string `mapstructure:"host"`
}

type ConfLogger struct {
	GinFilePath    string `mapstructure:"gin_file"`
	CommonFilePath string `mapstructure:"common_file"`
	GinLevel       string `mapstructure:"gin_level"`
	CommonLevel    string `mapstructure:"common_level"`
	StdoutLog      bool   `mapstructure:"stdout_log"`
}

func NewConfig() *Config {
	setDefaultDb()
	setDefaultWeb()
	setDefaultLog()

	confDir, confFile, confExt := splitPath(getConfPath())
	Lg("configs", "newConfig").
		Infof("Config dir: '%s', configs file name: '%s', configs ext: '%s'", confDir, confFile, confExt)

	viper.SetConfigName(confFile)
	viper.SetConfigType(confExt)
	viper.AddConfigPath(confDir)

	err := viper.ReadInConfig()
	if err != nil {
		Lg("configs", "newConfig").Fatal("Fatal error configs file ", err)
	}

	conf := new(Config)

	er := viper.Unmarshal(conf)
	if er != nil {
		Lg("configs", "newConfig").Fatal("Fatal error configs file:", er)
	}

	return conf
}

func setDefaultDb() {
	viper.SetDefault("database.postgres.driver_name", "default")
	viper.SetDefault("database.postgres.username", "default")
	viper.SetDefault("database.postgres.password", "0000")
	viper.SetDefault("database.postgres.db_name", "default")
	viper.SetDefault("database.postgres.ssl_mode", "disable")
	viper.SetDefault("database.postgres.host", "disable")
}

func setDefaultWeb() {
	viper.SetDefault("web.server.address", "")
	viper.SetDefault("web.server.port", "0000")
	viper.SetDefault("web.server.host", "pinterest-tp.tk")
	viper.SetDefault("web.server.protocol", "http://")

	viper.SetDefault("web.chat_srv.address", "localhost")
	viper.SetDefault("web.chat_srv.port", "8000")
	viper.SetDefault("web.chat_srv.host", "pinterest-tp.tk")
	viper.SetDefault("web.chat_srv.protocol", "http://")

	viper.SetDefault("web.static.dir_img", "img")
	viper.SetDefault("web.static.url_img", "img")
	viper.SetDefault("web.static.dir_avt", "/static/avatar/")
}

func setDefaultLog() {
	viper.SetDefault("logger.gin_file", "/var/log/pinterest/gin.log")
	viper.SetDefault("logger.common_file", "/var/log/pinterest/common.log")
	viper.SetDefault("logger.gin_level", "debug")
	viper.SetDefault("logger.common_level", "debug")
	viper.SetDefault("logger.stdout_log", true)
}

func getConfPath() string {
	argPath := ""
	flag.StringVar(&argPath, "configs", "", "Specify configs path")
	flag.Parse()

	if argPath != "" {
		return argPath
	}

	if confPath, ok := os.LookupEnv("CONF_PATH"); ok {
		return confPath
	}

	return "./configs/yaml/configs.yaml"
}

func splitPath(path string) (dir, fileName, ext string) {
	dir, file := filepath.Split(path)
	strs := strings.Split(file, ".")

	if len(strs) == 2 {
		return dir, strs[0], strs[1]
	} else if len(strs) == 1 {
		return dir, strs[0], ""
	} else {
		return dir, "", ""
	}
}
