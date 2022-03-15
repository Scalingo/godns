package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
)

var (
	settings Settings
)

var LogLevelMap = map[string]int{
	"DEBUG":  LevelDebug,
	"INFO":   LevelInfo,
	"NOTICE": LevelNotice,
	"WARN":   LevelWarn,
	"ERROR":  LevelError,
}

type Settings struct {
	Version      string
	Debug        bool
	Server       DNSServerSettings       `toml:"server"`
	ResolvConfig ResolvSettings          `toml:"resolv"`
	Redis        RedisSettings           `toml:"redis"`
	Memcache     MemcacheSettings        `toml:"memcache"`
	Log          LogSettings             `toml:"log"`
	Cache        CacheSettings           `toml:"cache"`
	Hosts        HostsSettings           `toml:"hosts"`
	Zones        map[string]ZoneSettings `toml:"zones"`
}

type ResolvSettings struct {
	Timeout        int
	Interval       int
	SetEDNS0       bool
	ServerListFile string `toml:"server-list-file"`
	ResolvFile     string `toml:"resolv-file"`
	IPv6           bool   `toml:"ipv6"`
}

type DNSServerSettings struct {
	Host string
	Port int
}

type RedisSettings struct {
	Host     string
	Port     int
	DB       int
	Password string
}

type MemcacheSettings struct {
	Servers []string
}

func (s RedisSettings) Addr() string {
	return s.Host + ":" + strconv.Itoa(s.Port)
}

type LogSettings struct {
	Stdout bool
	File   string
	Level  string
}

func (ls LogSettings) LogLevel() int {
	l, ok := LogLevelMap[ls.Level]
	if !ok {
		panic("Config error: invalid log level: " + ls.Level)
	}
	return l
}

type CacheSettings struct {
	Backend  string
	Expire   int
	Maxcount int
}

type ZoneSettings struct {
	Name        string `toml:"name"`
	Ns          string `toml:"ns"`
	Mbox        string `toml:"mbox"`
	Serial      uint32 `toml:"serial"`
	Refresh     uint32 `toml:"refresh"`
	Retry       uint32 `toml:"retry"`
	Expire      uint32 `toml:"expire"`
	NegcacheTtl uint32 `toml:"negcache-ttl"`
	SoaTtl      uint32 `toml:"soa-ttl"`
}

type HostsSettings struct {
	Enable          bool
	HostsFile       string `toml:"host-file"`
	RedisEnable     bool   `toml:"redis-enable"`
	RedisKey        string `toml:"redis-key"`
	TTL             uint32 `toml:"ttl"`
	RefreshInterval uint32 `toml:"refresh-interval"`
	Zone            string `toml:"zone"`
	ZoneNs          string `toml:"zone-ns"`
	ZoneMbox        string `toml:"zone-mbox"`
	ZoneSerial      uint32 `toml:"zone-serial"`
	ZoneRefresh     uint32 `toml:"zone-refresh"`
	ZoneRetry       uint32 `toml:"zone-retry"`
	ZoneExpire      uint32 `toml:"zone-expire"`
	ZoneNegcacheTtl uint32 `toml:"zone-negcache-ttl"`
	ZoneSoaTtl      uint32 `toml:"zone-soa-ttl"`
}

func init() {

	var configFile string
	var verbose bool

	flag.StringVar(&configFile, "c", "./etc/godns.conf", "Look for godns toml-formatting config file in this directory")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.Parse()

	if _, err := toml.DecodeFile(configFile, &settings); err != nil {
		fmt.Printf("%s is not a valid toml config file\n", configFile)
		fmt.Println(err)
		os.Exit(1)
	}

	if settings.Zones == nil {
		settings.Zones = make(map[string]ZoneSettings)
	}

	if len(settings.Hosts.Zone) > 0 {
		settings.Zones["default"] = ZoneSettings{
			Name:        settings.Hosts.Zone,
			Ns:          settings.Hosts.ZoneNs,
			Mbox:        settings.Hosts.ZoneMbox,
			Serial:      settings.Hosts.ZoneSerial,
			Refresh:     settings.Hosts.ZoneRefresh,
			Retry:       settings.Hosts.ZoneRetry,
			Expire:      settings.Hosts.ZoneExpire,
			NegcacheTtl: settings.Hosts.ZoneNegcacheTtl,
			SoaTtl:      settings.Hosts.ZoneSoaTtl,
		}
	}

	if verbose {
		settings.Log.Stdout = true
		settings.Log.Level = "DEBUG"
	}
}
