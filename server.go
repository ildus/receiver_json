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

// Represents row created by decoder_json
type Query struct {
	Action string                 `json:"a"`
	Data   map[string]interface{} `json:"d"`
	Clause map[string]interface{} `json:"c"`
}

type Configuration struct {
	PgConnectionStr string `json:"pg_connection_str"`
	SlotName        string `json:"slot_name"`
	HekaAddr        string `json:"heka_addr"`
}

type Rules struct {
	Constants map[string]int         `json:"constants"`
	Tables    map[string]interface{} `json:"tables"`
}

func parseRow(row string) {

}

func loadRules() {
	confData, err := ioutil.ReadFile(*rulesLocation)
	if err != nil {
		logger.Fatal(err)
	}

	err = json.Unmarshal(confData, &rules)
	if err != nil {
		logger.Fatal("Rules decoding error: ", err)
	}
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