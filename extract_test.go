package structextract

import (
	"reflect"
	"strings"
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

func TestExtractor_NamesFromTagWithPrefix(t *testing.T) {
	ext := fakeData()
	prefix := "default_"
	res, err := ext.NamesFromTagWithPrefix("json", prefix)
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if !strings.Contains(res[1], prefix) {
		t.Fatalf("prefix was not applied")
	}

}

func TestExtractor_NamesFromTagWithPrefix_Empty_Tag(t *testing.T) {

	ext := fakeData()
	prefix := "default_"
	out, err := ext.NamesFromTagWithPrefix("", prefix)
	if err != nil {
		t.Fatalf("error should be null")
	}
	if len(out) != 0 {
		t.Fatalf("no objects was expected")
	}

}
func TestExtractor_NamesFromTagWithPrefix_No_Prefix(t *testing.T) {
	ext := fakeData()
	resWith, err := ext.NamesFromTagWithPrefix("json", "")
	if err != nil {
		t.Fatalf("unexpected error")
	}
	resWithOut, err := ext.NamesFromTag("json")
	if err != nil {
		t.Fatalf("unexpected error")
	}

	if !reflect.DeepEqual(resWith, resWithOut) {
		t.Fatalf("slices micmatch")
	}

}

func TestExtractor_NamesFromTagWithPrefix_InvalidStruct(t *testing.T) {
	test := []string{"fail", "fail2"}
	ext := New(&test)

	_, err := ext.NamesFromTagWithPrefix("json", "default-")
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
	ext := fakeData().IgnoreField("Field4")
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

func TestExtractor_ValuesFromTag_Invalid_Struct(t *testing.T) {
	test := []string{"fail", "fail2"}
	ext := New(&test)

	_, err := ext.ValuesFromTag("json")
	if err == nil {
		t.Fatal("Passed value is not a valid struct")
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
func TestExtractor_GetFieldNamesFromTagWithPrefixIgnore(t *testing.T) {
	ext := fakeIgnoredData()
	exp := []string{
		"default_field_1",
		"default_field_3",
	}
	res, _ := ext.NamesFromTagWithPrefix("json", "default_")

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

func TestTagMapping(t *testing.T) {
	type test struct {
		FieldA string `json:"fieldA" sql:"field_a"`
		FieldB string `json:"fieldB" sql:"field_b"`
		FieldC string `sql:"field_c"`
		FieldD string `json:"fieldD"`
	}

	ext := New(&test{})
	expected := map[string]string{
		"fieldA": "field_a",
		"fieldB": "field_b",
	}

	mapping, err := ext.TagMapping("json", "sql")
	if err != nil {
		t.Fatalf("encountered error (%s) whilst mapping tagged fields", err)
	}

	for k, v := range expected {
		value, ok := mapping[k]
		if !ok {
			t.Errorf("mapping for key %s not found", k)
			continue
		}
		if value != v {
			t.Errorf("expected mapping mapping to be %s -> %s, got %s -> %s", k, v, k, value)
		}
	}

}

func TestTagMapping_ignoredFields(t *testing.T) {
	type test struct {
		FieldA string `json:"fieldA" sql:"field_a"`
		FieldB string `json:"fieldB" sql:"field_b"`
		FieldC string `sql:"field_c"`
		FieldD string `json:"fieldD"`
	}

	ext := New(&test{}).IgnoreField("FieldA")
	expected := map[string]string{
		"fieldB": "field_b",
	}

	mapping, err := ext.TagMapping("json", "sql")
	if err != nil {
		t.Fatalf("encountered error (%s) whilst mapping tagged fields", err)
	}

	for k, v := range expected {
		value, ok := mapping[k]
		if !ok {
			t.Errorf("mapping for key %s not found", k)
			continue
		}
		if value != v {
			t.Errorf("expected mapping mapping to be %s -> %s, got %s -> %s", k, v, k, value)
		}
	}

}

func TestTagMapping_invalidStruct(t *testing.T) {
	type test struct {
		Field string
	}

	ext := New(test{})
	_, err := ext.TagMapping("json", "sql")
	if err == nil {
		t.Fatal("Passed value is not a valid struct")
	}
}

func TestEmbeddedStructs_togglingBehaviour(t *testing.T) {
	type Embed struct {
		AnotherField string
	}
	type Outer struct {
		Embed
		Field string
	}

	ts := Outer{Embed{"another"}, "some"}

	ext := New(&ts)
	v, err := ext.Values()
	if err != nil {
		t.Fatal("failed to get values when not using embedded structs")
	}

	if len(v) != 1 {
		t.Fatalf("expected a single value, got %d", len(v))
	}

	if !reflect.DeepEqual(v[0], "some") {
		t.Fatalf("expected the single value to be 'some', got %v", v[0])
	}

	v, err = ext.UseEmbeddedStructs(true).Values()
	if err != nil {
		t.Fatalf("failed to get values when using embedded structs: %s", err)
	}

	if len(v) != 2 {
		t.Fatalf("expected to get 2 values when using embedded struct")
	}
}

type basicTypes struct {
	BoolType           bool     `custom:"boolType" custom_two:"boolTypeTwo,omitempty"`
	StringType         string   `custom:"stringType,omitempty"`
	IntType            int      `custom:"intType,omitempty"`
	ByteType           []byte   `custom:"byteType,omitempty"`
	Float64Type        float64  `custom:"float64Type,omitempty"`
	BoolTypePtr        *bool    `custom:"boolTypePtr,omitempty"`
	StringTypePtr      *string  `custom:"stringTypePtr,omitempty"`
	IntTypePtr         *int     `custom:"intTypePtr,omitempty"`
	Float64TypePtr     *float64 `custom:"float64TypePtr,omitempty"`
	FieldWithNoOmitTag string   `custom:"fieldWithNoOmitTag"`
}

func TestExtractor_FieldValueFromTagMapOmitempty(t *testing.T) {
	testBool := false
	testStr := "test"
	testInt := 6
	testFloat64 := 1.2

	tests := []struct {
		name     string
		structIn basicTypes
		expected map[string]interface{}
	}{
		{
			name:     "all fields empty, expect not omitted fields: boolType and fieldWithNoOmitTag",
			structIn: basicTypes{},
			expected: map[string]interface{}{
				"boolType":           false,
				"fieldWithNoOmitTag": "",
			},
		},
		{
			name: "all fields initialised, expect all fields back",
			structIn: basicTypes{
				BoolType:           true,
				StringType:         testStr,
				IntType:            1,
				ByteType:           []byte("test"),
				Float64Type:        1.2,
				BoolTypePtr:        &testBool,
				StringTypePtr:      &testStr,
				IntTypePtr:         &testInt,
				Float64TypePtr:     &testFloat64,
				FieldWithNoOmitTag: testStr,
			},
			expected: map[string]interface{}{
				"boolType":           true,
				"stringType":         testStr,
				"intType":            1,
				"byteType":           []byte("test"),
				"float64Type":        1.2,
				"boolTypePtr":        &testBool,
				"stringTypePtr":      &testStr,
				"intTypePtr":         &testInt,
				"float64TypePtr":     &testFloat64,
				"fieldWithNoOmitTag": testStr,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := New(&test.structIn).FieldValueFromTagMap("custom")
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(test.expected, result) {
				t.Fatalf("want %v, got %v", test.expected, result)
			}
		})
	}
}

func TestExtractor_NamesFromTagOmitempty(t *testing.T) {
	testBool := false
	testStr := "test"
	testInt := 6
	testFloat64 := 1.2

	tests := []struct {
		name     string
		structIn basicTypes
		expected []string
	}{
		{
			name:     "all fields empty, expect not omitted fields: boolType and fieldWithNoOmitTag",
			structIn: basicTypes{},
			expected: []string{
				"boolType",
				"fieldWithNoOmitTag",
			},
		},
		{
			name: "all fields initialised, expect all fields back",
			structIn: basicTypes{
				BoolType:           true,
				StringType:         testStr,
				IntType:            1,
				ByteType:           []byte("test"),
				Float64Type:        1.2,
				BoolTypePtr:        &testBool,
				StringTypePtr:      &testStr,
				IntTypePtr:         &testInt,
				Float64TypePtr:     &testFloat64,
				FieldWithNoOmitTag: testStr,
			},
			expected: []string{
				"boolType",
				"stringType",
				"intType",
				"byteType",
				"float64Type",
				"boolTypePtr",
				"stringTypePtr",
				"intTypePtr",
				"float64TypePtr",
				"fieldWithNoOmitTag",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := New(&test.structIn).NamesFromTag("custom")
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(test.expected, result) {
				t.Fatalf("want %v, got %v", test.expected, result)
			}
		})
	}
}

func TestExtractor_NamesFromTagWithPrefixOmitempty(t *testing.T) {
	testBool := false
	testStr := "test"
	testInt := 6
	testFloat64 := 1.2

	tests := []struct {
		name     string
		structIn basicTypes
		expected []string
	}{
		{
			name:     "all fields empty, expect not omitted fields: boolType and fieldWithNoOmitTag",
			structIn: basicTypes{},
			expected: []string{
				"test_boolType",
				"test_fieldWithNoOmitTag",
			},
		},
		{
			name: "all fields initialised, expect all fields back",
			structIn: basicTypes{
				BoolType:           true,
				StringType:         testStr,
				IntType:            1,
				ByteType:           []byte("test"),
				Float64Type:        1.2,
				BoolTypePtr:        &testBool,
				StringTypePtr:      &testStr,
				IntTypePtr:         &testInt,
				Float64TypePtr:     &testFloat64,
				FieldWithNoOmitTag: testStr,
			},
			expected: []string{
				"test_boolType",
				"test_stringType",
				"test_intType",
				"test_byteType",
				"test_float64Type",
				"test_boolTypePtr",
				"test_stringTypePtr",
				"test_intTypePtr",
				"test_float64TypePtr",
				"test_fieldWithNoOmitTag",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := New(&test.structIn).NamesFromTagWithPrefix("custom", "test_")
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(test.expected, result) {
				t.Fatalf("want %v, got %v", test.expected, result)
			}
		})
	}
}

func TestExtractor_ValuesFromTagOmitempty(t *testing.T) {
	testBool := false
	testStr := "test"
	testInt := 6
	testFloat64 := 1.2

	tests := []struct {
		name     string
		structIn basicTypes
		expected []interface{}
	}{
		{
			name:     "all fields empty, expect not omitted fields: boolType and fieldWithNoOmitTag",
			structIn: basicTypes{},
			expected: []interface{}{
				false,
				"",
			},
		},
		{
			name: "all fields initialised, expect all fields back",
			structIn: basicTypes{
				BoolType:           true,
				StringType:         testStr,
				IntType:            1,
				ByteType:           []byte("test"),
				Float64Type:        1.2,
				BoolTypePtr:        &testBool,
				StringTypePtr:      &testStr,
				IntTypePtr:         &testInt,
				Float64TypePtr:     &testFloat64,
				FieldWithNoOmitTag: testStr,
			},
			expected: []interface{}{
				true,
				testStr,
				1,
				[]byte("test"),
				1.2,
				&testBool,
				&testStr,
				&testInt,
				&testFloat64,
				testStr,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := New(&test.structIn).ValuesFromTag("custom")
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(test.expected, result) {
				t.Fatalf("want %v, got %v", test.expected, result)
			}
		})
	}
}

