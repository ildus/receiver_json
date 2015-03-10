package receiver_json

import (
	"encoding/json"
	"io/ioutil"
)

type Rules struct {
	Constants map[string]int         `json:"constants"`
	Tables    map[string]interface{} `json:"tables"`
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
