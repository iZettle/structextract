package structextract

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

// Extractor holds the struct that we want to extract data from
type Extractor struct {
	StructAddr         interface{} // StructAddr: struct address
	ignoredFields      []string    // ignoredFields: an array with all the fields to be ignored
	useEmbeddedStructs bool
}

// New returns a new Extractor struct
// the parameter have to be a pointer to a struct
func New(s interface{}) *Extractor {
	return &Extractor{
		StructAddr:         s,
		ignoredFields:      nil,
		useEmbeddedStructs: false,
	}
}

//Names returns an array with all the field names (with the same order) as defined on the struct
func (e *Extractor) Names() (out []string, err error) {

	if err := e.isValidStruct(); err != nil {
		return nil, err
	}

	s := reflect.ValueOf(e.StructAddr).Elem()
	fields := e.fields(s)
	for _, field := range fields {
		out = append(out, field.name)
	}

	return
}

//NamesFromTag returns an array with all the tag names for each field
func (e *Extractor) NamesFromTag(tag string) (out []string, err error) {

	if err := e.isValidStruct(); err != nil {
		return nil, err
	}

	s := reflect.ValueOf(e.StructAddr).Elem()
	fields := e.fields(s)

	for _, field := range fields {
		if val, ok := field.tags.Lookup(tag); ok {
			out = append(out, val)
		}
	}

	return
}

//NamesFromTagWithPrefix returns an array with all the tag names for each field including the given prefix
func (e *Extractor) NamesFromTagWithPrefix(tag string, prefix string) (out []string, err error) {

	if err := e.isValidStruct(); err != nil {
		return nil, err
	}

	s := reflect.ValueOf(e.StructAddr).Elem()
	fields := e.fields(s)

	for _, field := range fields {
		val, ok := field.tags.Lookup(tag)
		if !ok {
			continue
		}
		out = append(out, strings.TrimSpace(prefix+val))
	}

	return
}

//Values returns an interface array with all the values
func (e *Extractor) Values() (out []interface{}, err error) {

	if err := e.isValidStruct(); err != nil {
		return nil, err
	}

	s := reflect.ValueOf(e.StructAddr).Elem()
	fields := e.fields(s)

	for _, field := range fields {
		out = append(out, field.value.Interface())

	}

	return
}

//ValuesFromTag returns an interface array with all the values of fields with the given tag
func (e *Extractor) ValuesFromTag(tag string) (out []interface{}, err error) {

	if err := e.isValidStruct(); err != nil {
		return nil, err
	}

	s := reflect.ValueOf(e.StructAddr).Elem()
	fields := e.fields(s)

	for _, field := range fields {
		if _, ok := field.tags.Lookup(tag); ok {
			out = append(out, field.value.Interface())
		}

	}

	return
}

// FieldValueMap returns a string to interface map,
// key: field as defined on the struct
// value: the value of the field
func (e *Extractor) FieldValueMap() (out map[string]interface{}, err error) {

	if err := e.isValidStruct(); err != nil {
		return nil, err
	}

	out = make(map[string]interface{})
	s := reflect.ValueOf(e.StructAddr).Elem()
	fields := e.fields(s)

	for _, field := range fields {
		out[field.name] = field.value.Interface()
	}

	return
}

// FieldValueFromTagMap returns a string to interface map that uses as key the tag name,
// key: tag name for the given field
// value: the value of the field
func (e *Extractor) FieldValueFromTagMap(tag string) (out map[string]interface{}, err error) {

	if err := e.isValidStruct(); err != nil {
		return nil, err
	}

	out = make(map[string]interface{})
	s := reflect.ValueOf(e.StructAddr).Elem()
	fields := e.fields(s)

	for _, field := range fields {
		if val, ok := field.tags.Lookup(tag); ok {
			out[val] = field.value.Interface()
		}

	}

	return
}

// TagMapping returns a map that maps tagged fields from one tag to another.
// This can help with mapping partial JSON objects to some other kind of a
// mapping, such as SQL. It only maps existing field pairs, if either field
// does not have a tag, it's left out.
func (e *Extractor) TagMapping(from, to string) (out map[string]string, err error) {
	if err := e.isValidStruct(); err != nil {
		return nil, err
	}

	out = make(map[string]string)
	s := reflect.ValueOf(e.StructAddr).Elem()
	fields := e.fields(s)

	for _, field := range fields {
		fromTag, fromOk := field.tags.Lookup(from)
		toTag, toOk := field.tags.Lookup(to)
		if toOk && fromOk {
			out[fromTag] = toTag
		}
	}

	return
}

// IgnoreField checks if the given fields are valid based on the given struct,
// then append them on the ignore list
// e.g. ext := structextract.New(&business).IgnoreField("ID","DateModified")
func (e *Extractor) IgnoreField(fd ...string) *Extractor {

	if err := e.isValidStruct(); err != nil {
		return e
	}
	for _, field := range fd {
		if e.isFieldNameValid(field) {
			e.ignoredFields = append(e.ignoredFields, field)
		}
	}

	return e
}

// UseEmbeddedStructs toggles the usage of embedded structs
func (e *Extractor) UseEmbeddedStructs(use bool) *Extractor {
	e.useEmbeddedStructs = use
	return e
}

// GetChangesetForTag takes a dataset and tries to change the keys from one
// tag to another. This is helpful when trying to create partial updates to a
// row in a database.
func (e *Extractor) GetChangesetForTag(data map[string]interface{}, inputTag, outputTag string) (out map[string]interface{}, err error) {
	if err = e.isValidStruct(); err != nil {
		return
	}

	out = make(map[string]interface{})
	s := reflect.ValueOf(e.StructAddr).Elem()
	fields := e.fields(s)
	for _, field := range fields {

		inputFieldTag, ok := field.tags.Lookup(inputTag)
		if !ok {
			continue
		}
		outputFieldTag, ok := field.tags.Lookup(outputTag)
		if !ok {
			continue
		}
		value, ok := data[inputFieldTag]
		if !ok {
			continue
		}
		out[outputFieldTag] = value
	}

	return
}

// ApplyMap takes a map of arbitrary values and tries to apply them to a
// struct, matching keys with the given input tag. If the input tag is empty,
// the key will be matched against the field name
func (e *Extractor) ApplyMap(data map[string]interface{}, inputTag string) (interface{}, error) {
	err := e.isValidStruct()
	if err != nil {
		return nil, err
	}

	validFields := map[string]bool{}
	s := reflect.ValueOf(e.StructAddr).Elem()
	fields := e.fields(s)
	for _, field := range fields {

		var inputFieldTag string
		var ok bool
		if inputTag == "" {
			inputFieldTag = field.name
		} else {
			inputFieldTag, ok = field.tags.Lookup(inputTag)
			if !ok {
				continue
			}
		}

		_, ok = data[inputFieldTag]
		if !ok {
			continue
		}

		fieldName := field.name
		validFields[fieldName] = true

	}

	// Create a new zero value of the type stored in e.StructAddr
	out := reflect.New(s.Type()).Interface()
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// dirty hax, as some fields may be deserialized differently depending
	// on whether they're being deserialized into a map[string]interface{}
	// or a user specified type.
	err = json.Unmarshal(bytes, out)
	if err != nil {
		return nil, err
	}

	outRef := reflect.ValueOf(out).Elem()
	partiallyApply(s, outRef, validFields)

	return out, nil
}

// This function applies field values from a to b if the field name is not in
// the ignored fields.
// It is expected that both a and b are of the same type, othwerise panics will ensue
func partiallyApply(a, b reflect.Value, ignoredFields map[string]bool) {
	for i := 0; i < a.NumField(); i++ {
		if a.Type().Field(i).Anonymous {
			partiallyApply(a.Field(i), b.Field(i), ignoredFields)
			continue
		}
		fieldName := a.Type().Field(i).Name
		ok := ignoredFields[fieldName]
		// check if field is already set because it was in data
		if ok {
			continue
		}

		b.Field(i).Set(a.Field(i))
	}
}

func (e *Extractor) isFieldNameValid(fn string) bool {

	s := reflect.ValueOf(e.StructAddr).Elem()
	fields := e.fields(s)

	for _, field := range fields {
		if field.name == fn {
			return true
		}
	}

	return false
}

func isIgnored(a string, list []string) bool {
	for _, l := range list {
		if l == a {
			return true
		}
	}
	return false
}

func (e *Extractor) isValidStruct() error {

	stVal := reflect.ValueOf(e.StructAddr)
	if stVal.Kind() != reflect.Ptr || stVal.IsNil() {
		return errors.New("struct passed is not valid, a pointer was expected")
	}
	structVal := stVal.Elem()
	if structVal.Kind() != reflect.Struct {
		return errors.New("struct passed is not valid, a pointer to struct was expected")
	}

	return nil
}

type field struct {
	value reflect.Value
	name  string
	tags  reflect.StructTag
}

// This function returns a slice of fields of a struct
// as reflect.Value, even fields of embedded structs
func (e *Extractor) fields(s reflect.Value) []field {
	fields := make([]field, 0, s.NumField())

	for i := 0; i < s.NumField(); i++ {
		if isIgnored(s.Type().Field(i).Name, e.ignoredFields) {
			continue
		}

		if s.Type().Field(i).Anonymous {
			if e.useEmbeddedStructs {
				fields = append(fields, e.fields(s.Field(i))...)
			}
			continue
		}

		tag := s.Type().Field(i).Tag
		name := s.Type().Field(i).Name
		value := s.Field(i)
		fields = append(fields, field{value, name, tag})
	}

	return fields
}
