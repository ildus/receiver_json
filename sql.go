package receiver_json

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type ActionType int

const (
	INSERT ActionType = iota
	UPDATE ActionType = iota
	DELETE ActionType = iota
)

// Represents row created by decoder_json
type Query struct {
	Action   ActionType             `json:"a"`
	Data     map[string]interface{} `json:"d"`
	Clause   map[string]interface{} `json:"c"`
	Relation string                 `json:"r"`
}

/* Get representation of value */
func getRepr(val interface{}) (*string, error) {
	var repr string
	switch val.(type) {
	case nil:
		repr = "NULL"
	case string:
		sval := val.(string)
		sval = strings.Replace(sval, "'", "\\'", -1)
		repr = "'" + sval + "'"
	case int:
		repr = strconv.Itoa(val.(int))
	case float32:
		repr = strconv.FormatFloat(float64(val.(float32)), 'E', 6, 32)
	case float64:
		repr = strconv.FormatFloat(val.(float64), 'E', 15, 64)
	case bool:
		if val.(bool) {
			repr = "TRUE"
		} else {
			repr = "FALSE"
		}
	case []interface{}:
		repr := "ARRAY[%s]"
		items := make([]string, 0)
		for item, _ := range val.([]interface{}) {
			itemRepr, err := getRepr(item)
			if err != nil {
				return nil, err
			}
			items = append(items, *itemRepr)
		}
		repr = fmt.Sprintf(repr, strings.Join(items, ","))
	default:
		return nil, errors.New(fmt.Sprintf("Cannot convert value: %s", val))
	}
	return &repr, nil
}

/* Generate equality from field name and its value */
func generateEqual(name string, val interface{}) (*string, error) {
	var result string
	repr, err := getRepr(val)
	if err != nil {
		return nil, err
	}
	result = fmt.Sprintf("%s=%s", name, *repr)
	return &result, nil
}

/* Generate concatenation of field=val from pairs */
func generateEqualPairs(clause *map[string]interface{}, connector string) (*string, error) {
	var conditions []string
	for name, val := range *clause {
		equality, err := generateEqual(name, val)
		if err != nil {
			return nil, err
		}
		conditions = append(conditions, *equality)
	}
	result := strings.Join(conditions, connector)
	return &result, nil
}

/* Generate field and value list for INSERT query */
func generateFieldAndValues(data *map[string]interface{}) (*string, *string, error) {
	var fields, values []string
	for name, val := range *data {
		fields = append(fields, name)
		repr, err := getRepr(val)
		if err != nil {
			return nil, nil, err
		}
		values = append(values, *repr)
	}
	fieldsConcat := strings.Join(fields, ",")
	valuesConcat := strings.Join(values, ",")
	return &fieldsConcat, &valuesConcat, nil
}

func (q *Query) GenerateSQL() (*string, error) {
	var query string
	var clause *string
	var err error

	if q.Action != INSERT {
		if len(q.Clause) == 0 {
			return nil, errors.New("condition required")
		}
		clause, err = generateEqualPairs(&q.Clause, " AND ")
		if err != nil {
			return nil, err
		}
	}

	if q.Action != DELETE && len(q.Data) == 0 {
		return nil, errors.New("data required")
	}

	switch q.Action {
	case INSERT:
		var fields, values *string
		format := "INSERT INTO %s (%s) VALUES (%s);"
		fields, values, err = generateFieldAndValues(&q.Data)
		if err != nil {
			return nil, err
		}
		query = fmt.Sprintf(format, q.Relation, *fields, *values)
	case UPDATE:
		format := "UPDATE %s SET %s WHERE %s;"
		var pairs *string
		pairs, err = generateEqualPairs(&q.Data, ", ")
		if err != nil {
			return nil, err
		}
		query = fmt.Sprintf(format, q.Relation, *pairs, *clause)
	case DELETE:
		format := "DELETE FROM %s WHERE %s;"
		query = fmt.Sprintf(format, q.Relation, *clause)
	}

	return &query, nil
}
