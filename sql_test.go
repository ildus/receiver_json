package receiver_json

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepr(t *testing.T) {
	repr, _ := getRepr(float32(1.1234234))
	assert.Equal(t, "1.123423E+00", *repr, "float fail")

	repr, _ = getRepr(float64(1.12342e10))
	assert.Equal(t, "1.123420000000000E+10", *repr, "float fail")

	repr, _ = getRepr(2344)
	assert.Equal(t, "2344", *repr, "int fail")

	repr, _ = getRepr("str''")
	assert.Equal(t, "'str\\'\\''", *repr, "str fail")

	repr, _ = getRepr(true)
	assert.Equal(t, "TRUE", *repr, "bool 'true' fail")

	repr, _ = getRepr(false)
	assert.Equal(t, "FALSE", *repr, "bool 'false' fail")

	repr, _ = getRepr(nil)
	assert.Equal(t, "NULL", *repr, "NULL fail")
}

func TestGenerateEqual(t *testing.T) {
	eq, _ := generateEqual("val1", 5)
	assert.Equal(t, "val1=5", *eq, "fail")
}

func TestPairsGeneration(t *testing.T) {
	m := &map[string]interface{}{
		"val1": 5,
		"val2": true,
	}
	pairs, _ := generateEqualPairs(m, " AND ")
	cond := *pairs == "val1=5 AND val2=TRUE" || *pairs == "val2=TRUE AND val1=5"
	assert.True(t, cond, "fail")

	pairs, _ = generateEqualPairs(m, ", ")
	cond = *pairs == "val1=5, val2=TRUE" || *pairs == "val2=TRUE, val1=5"
	assert.True(t, cond, "fail")
}

func TestGenerateFieldAndValues(t *testing.T) {
	m := &map[string]interface{}{
		"val1": 5,
		"val2": true,
	}
	fields, values, _ := generateFieldAndValues(m)
	if *fields == "val1,val2" {
		assert.Equal(t, "5,TRUE", *values, "values fail")
	} else if *fields == "val2,val1" {
		assert.Equal(t, "TRUE,5", *values, "values fail")
	} else {
		t.Fail()
	}
}

func TestSQLGeneration(t *testing.T) {
	var err error
	var actual *string

	assert.Equal(t, 0, INSERT)
	assert.Equal(t, 1, UPDATE)
	assert.Equal(t, 2, DELETE)

	q1 := &Query{
		Action:   INSERT,
		Relation: "public.t",
	}

	expected := "INSERT INTO public.t (val1) VALUES (0);"
	actual, err = q1.GenerateSQL()
	assert.Equal(t, "data required", err.Error(), "fail")

	q1.Data = map[string]interface{}{"val1": 0}
	actual, _ = q1.GenerateSQL()
	assert.Equal(t, expected, *actual, "INSERT query fail")

	q2 := &Query{
		Action:   UPDATE,
		Data:     map[string]interface{}{"val1": 0},
		Relation: "public.t",
	}

	expected = "UPDATE public.t SET val1=0 WHERE id=2;"
	actual, err = q2.GenerateSQL()
	assert.Equal(t, "condition required", err.Error(), "fail")

	q2.Clause = map[string]interface{}{"id": 2}
	actual, _ = q2.GenerateSQL()
	assert.Equal(t, expected, *actual, "UPDATE query fail")

	expected = "DELETE FROM public.t2 WHERE id=3;"
	q3 := &Query{
		Action:   DELETE,
		Clause:   map[string]interface{}{"id": 3},
		Relation: "public.t2",
	}
	actual, _ = q3.GenerateSQL()
	assert.Equal(t, expected, *actual, "DELETE query fail")
}
