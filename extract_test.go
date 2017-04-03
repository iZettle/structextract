package structextract

import (
	"reflect"
	"testing"
)

type testStruct struct {
	Field1 string      `json:"field_1"`
	Field2 string      `json:"field_2"`
	Field3 bool        `json:"field_3"`
	Field4 interface{} `json:"field_4"`
}

func fakeData() *Extractor {
	ts := testStruct{
		Field1: "hello",
		Field2: "world",
		Field3: true,
		Field4: "2016-10-10",
	}

	return New(&ts)
}
func fakeIgnoredData() *Extractor {
	ts := testStruct{
		Field1: "hello",
		Field2: "world",
		Field3: true,
		Field4: "2016-10-10",
	}
	ext := New(&ts).
		IgnoreField("Field2").
		IgnoreField("Field4")
	return ext
}
func TestExtractor_Names(t *testing.T) {

	ext := fakeData()
	exp := []string{
		"Field1",
		"Field2",
		"Field3",
		"Field4",
	}
	res, _ := ext.Names()

	if !reflect.DeepEqual(res, exp) {
		t.FailNow()
	}

}

func TestExtractor_Names_Invalid_Struct(t *testing.T) {
	test := "test"
	ext := New(&test)

	_, err := ext.Names()
	if err == nil {
		t.Fatal("Passed value is not a valid stract")
	}

}
func TestExtractor_NamesFromTag(t *testing.T) {
	ext := fakeData()
	exp := []string{
		"field_1",
		"field_2",
		"field_3",
		"field_4",
	}
	res, _ := ext.NamesFromTag("json")

	if !reflect.DeepEqual(res, exp) {
		t.FailNow()
	}
}

func TestExtractor_NamesFromTag_Invalid_Struct(t *testing.T) {
	test := []string{"fail", "fail2"}
	ext := New(&test)

	_, err := ext.NamesFromTag("json")
	if err == nil {
		t.Fatal("Passed value is not a valid stract")
	}

}

func TestExtractor_Values(t *testing.T) {
	ext := fakeData()
	exp := []interface{}{
		"hello",
		"world",
		true,
		"2016-10-10",
	}
	res, _ := ext.Values()

	if !reflect.DeepEqual(res, exp) {
		t.FailNow()
	}
}

func TestExtractor_Values_Invalid_Struct(t *testing.T) {
	test := []string{"fail", "fail2"}
	ext := New(&test)

	_, err := ext.Values()
	if err == nil {
		t.Fatal("Passed value is not a valid stract")
	}

}

func TestExtractor_FieldValueMap(t *testing.T) {
	ext := fakeData()
	exp := map[string]interface{}{
		"Field1": "hello",
		"Field2": "world",
		"Field3": true,
		"Field4": "2016-10-10",
	}
	res, _ := ext.FieldValueMap()

	if !reflect.DeepEqual(res, exp) {
		t.FailNow()
	}
}

func TestExtractor_FieldValueMap_Invalid_Struct(t *testing.T) {
	test := []string{"fail", "fail2"}
	ext := New(&test)

	_, err := ext.FieldValueMap()
	if err == nil {
		t.Fatal("Passed value is not a valid stract")
	}

}

func TestExtractor_FieldValueFromTagMap(t *testing.T) {
	ext := fakeData()
	exp := map[string]interface{}{
		"field_1": "hello",
		"field_2": "world",
		"field_3": true,
		"field_4": "2016-10-10",
	}
	res, _ := ext.FieldValueFromTagMap("json")

	if !reflect.DeepEqual(res, exp) {
		t.FailNow()
	}

}

func TestExtractor_FieldValueFromTagMap_Invalid_Struct(t *testing.T) {
	test := []string{"fail", "fail2"}
	ext := New(&test)

	_, err := ext.FieldValueFromTagMap("json")
	if err == nil {
		t.Fatal("Passed value is not a valid stract")
	}

}

func TestExtractor_GetFieldNamesIgnore(t *testing.T) {
	ext := fakeIgnoredData()
	exp := []string{
		"Field1",
		"Field3",
	}
	res, _ := ext.Names()

	if !reflect.DeepEqual(res, exp) {
		t.FailNow()
	}
}
func TestExtractor_GetFieldValueFromTagMapIgnore(t *testing.T) {
	ext := fakeIgnoredData()
	exp := map[string]interface{}{
		"field_1": "hello",
		"field_3": true,
	}
	res, _ := ext.FieldValueFromTagMap("json")

	if !reflect.DeepEqual(res, exp) {
		t.FailNow()
	}
}

func TestExtractor_GetFieldValueMapIgnore(t *testing.T) {
	ext := fakeIgnoredData()
	exp := map[string]interface{}{
		"Field1": "hello",
		"Field3": true,
	}
	res, _ := ext.FieldValueMap()

	if !reflect.DeepEqual(res, exp) {
		t.FailNow()
	}
}

func TestExtractor_GetFieldNamesFromTagIgnore(t *testing.T) {
	ext := fakeIgnoredData()
	exp := []string{
		"field_1",
		"field_3",
	}
	res, _ := ext.NamesFromTag("json")

	if !reflect.DeepEqual(res, exp) {
		t.FailNow()
	}
}

func TestExtractor_ValuesIgnore(t *testing.T) {
	ext := fakeIgnoredData()
	exp := []interface{}{
		"hello",
		true,
	}
	res, _ := ext.Values()

	if !reflect.DeepEqual(res, exp) {
		t.FailNow()
	}
}

func TestExtractor_FieldValueFromTagMapWrongTag(t *testing.T) {
	ext := fakeData()
	exp := map[string]interface{}{}
	res, _ := ext.FieldValueFromTagMap("json2")

	if !reflect.DeepEqual(res, exp) {
		t.FailNow()
	}

}
