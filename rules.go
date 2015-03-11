package receiver_json

import (
	"encoding/json"
	"io/ioutil"
)

type TableAction int

const (
	SKIP_TABLE          = iota
	MODIFY_TABLE_FIELDS = iota
	TRUNCATE_TABLE      = iota
)

type ColumnAction int

const (
	SKIP_COLUMN      = iota
	PREFIX_COLUMN_PK = iota
	REPLACE_VALUE    = iota
)

type Rules struct {
	Constants map[string]int         `json:"constants"`
	Tables    map[string]interface{} `json:"tables"`

	TableActions map[string]int
	FieldActions map[string]int
}

func loadRules(path string) *Rules {
	var rules Rules
	data, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Fatal(err)
	}

	err = json.Unmarshal(data, &rules)
	if err != nil {
		logger.Fatal("Rules decoding error: ", err)
	}
	return &rules
}

func (r *Rules) modifyFields(q *Query) {

}

func (r *Rules) Modify(q *Query) bool {
	if tableAction, ok := r.TableActions[q.Relation]; ok {
		switch TableAction {
		case SKIP_TABLE:
			return false
		case TRUNCATE_TABLE:
			//todo: execute truncate if rows exists
			return false
		case MODIFY_TABLE_FIELDS:
			r.modifyFields(q)
		}
	}
	return true
}
