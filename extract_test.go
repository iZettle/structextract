package structextract

import (
	"reflect"
	"testing"
)

type testStruct struct {
	Field1 string      `json:"field_1" db:"field1"`
	Field2 string      `json:"field_2" db:"field2"`
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
		IgnoreField("Field2", "Field4")
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
		t.Fatal("Passed value is not a valid struct")
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
		t.Fatal("Passed value is not a valid struct")
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
		t.Fatal("Passed value is not a valid struct")
	}

}

func TestExtractor_ValuesFromTag(t *testing.T) {
	ext := fakeData()
	exp := []interface{}{
		"hello",
		"world",
	}
	res, _ := ext.ValuesFromTag("db")

	expectedLength := len(exp)
	if len(res) != expectedLength {
		t.Fatalf("Number of values do not match: expected:%d, got:%d", expectedLength, len(res))
	}

	if !reflect.DeepEqual(res, exp) {
		t.FailNow()
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
		t.Fatal("Passed value is not a valid struct")
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
		t.Fatal("Passed value is not a valid struct")
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

func TestExtractor_IgnoreField_NotValidStruct(t *testing.T) {

	notAStruct := []string{"test", "test2"}
	ext := New(notAStruct).IgnoreField("test")

	if ext.ignoredFields != nil {
		t.Fatalf("not valid struct error was expected")
	}
}

func TestExtractor_IgnoreField_NotValidField(t *testing.T) {
	fk := fakeData()
	fk.IgnoreField("NotAValidField")
	exp := []interface{}{
		"hello",
		"world",
		true,
		"2016-10-10",
	}
	res, _ := fk.Values()
	if !reflect.DeepEqual(res, exp) {
		t.Fatalf("unexpected struct")
	}
}
