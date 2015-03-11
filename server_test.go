package receiver_json

import (
	//"fmt"
	"github.com/ildus/golog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func SetUp() {
	logger = golog.Default
}

func TestRulesLoading(t *testing.T) {
	SetUp()
	rules := loadRules("test_rules.json")
	assert.Exactly(t, 0, rules.Constants["table.skip"])
	assert.Exactly(t, 0, rules.Constants["field.pk_with_prefix"])

	table1info := rules.Tables["public.table1"]
	field1info := table1info.(map[string]interface{})["field1"]
	assert.Equal(t, 0, field1info.(map[string]interface{})["type"].(float64))
	assert.Equal(t, "i", field1info.(map[string]interface{})["prefix"].(string))
	//assert.Exactly(t, 0, rules.Tables.(map[string]map[string]interface{})["public.table1"]["field1"]["type"])
}
