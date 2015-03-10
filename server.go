package receiver_json

import (
	"database/sql"
	"encoding/json"
	"flag"
	"github.com/ildus/golog"
	"github.com/ildus/golog/appenders"
	_ "github.com/lib/pq"
	"io/ioutil"
)

var (
	confLocation  = flag.String("conf", "settings.json", "Settings location")
	rulesLocation = flag.String("rules-conf", "rules.json", "Rules location")
	rules         Rules
	conf          Configuration
	db            *sql.DB
	logger        *golog.Logger
)

type Configuration struct {
	PgConnectionStr string `json:"pg_connection_str"`
	SlotName        string `json:"slot_name"`
	HekaAddr        string `json:"heka_addr"`
}

func loadConfiguration() {
	confData, err := ioutil.ReadFile(*confLocation)
	if err != nil {
		logger.Fatal(err)
	}

	err = json.Unmarshal(confData, &conf)
	if err != nil {
		logger.Fatal("Configuration decoding error: ", err)
	}
}

func initDatabase() {
	var err error
	db, err = sql.Open("postgres", conf.PgConnectionStr+"?sslmode=disable")
	if err != nil {
		logger.Fatal(err)
	}
}

func Init() {
	logger = golog.Default
	if len(conf.HekaAddr) > 0 {
		logger.Enable(appenders.Heka(golog.Conf{
			"addr":         conf.HekaAddr,
			"proto":        "udp",
			"env_version":  "2",
			"message_type": "receiver_json",
		}))
	}

	loadConfiguration()
	loadRules()
	initDatabase()
}

func main() {
	Init()
}
