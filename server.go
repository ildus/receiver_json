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
	conf          *Configuration
	logger        *golog.Logger
)

type Configuration struct {
	PgConnectionStr string `json:"pg_connection_str"`
	SlotName        string `json:"slot_name"`
	HekaAddr        string `json:"heka_addr"`
}

func loadConfiguration(path string) *Configuration {
	var result Configuration
	confData, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Fatal(err)
	}

	err = json.Unmarshal(confData, &result)
	if err != nil {
		logger.Fatal("Configuration decoding error: ", err)
	}

	if len(result.PgConnectionStr) == 0 {
		logger.Fatal("Connection string is required")
	}
	if len(result.SlotName) == 0 {
		logger.Fatal("Slot name is required")
	}
	return &result
}

func initDatabase(connstr string) (db *sql.DB) {
	var err error
	db, err = sql.Open("postgres", connstr+"?sslmode=disable")
	if err != nil {
		logger.Fatal(err)
	}
	return db
}

func fetchRecords(db *sql.DB, rules *Rules) {

}

// Load configuration, init connection to database and start fetching from slot
func Init() {
	logger = golog.Default
	conf = loadConfiguration(*confLocation)
	if len(conf.HekaAddr) > 0 {
		logger.Enable(appenders.Heka(golog.Conf{
			"addr":         conf.HekaAddr,
			"proto":        "udp",
			"env_version":  "2",
			"message_type": "receiver_json",
		}))
	}
	db := initDatabase(conf.PgConnectionStr)
	rules := loadRules(*rulesLocation)
	go fetchRecords(db, rules)
}

func main() {
	Init()
}
