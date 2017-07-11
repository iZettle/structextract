package structextract

import (
	"errors"
	"reflect"
)

// Extractor holds the struct that we want to extract data from
type Extractor struct {
	StructAddr    interface{} // StructAddr: struct address
	ignoredFields []string    // ignoredFields: an array with all the fields to be ignored
}

// New returns a new Extractor struct
// the parameter have to be a pointer to a struct
func New(s interface{}) *Extractor {
	return &Extractor{
		StructAddr:    s,
		ignoredFields: nil,
	}
}

//Names returns an array with all the field names (with the same order) as defined on the struct
func (e *Extractor) Names() (out []string, err error) {

	if err := e.isValidStruct(); err != nil {
		return nil, err
	}

	s := reflect.ValueOf(e.StructAddr).Elem()
	for i := 0; i < s.NumField(); i++ {
		if isIgnored(s.Type().Field(i).Name, e.ignoredFields) {
			continue
		}
		out = append(out, s.Type().Field(i).Name)
	}

	return
}

//NamesFromTag returns an array with all the tag names for each field
func (e *Extractor) NamesFromTag(tag string) (out []string, err error) {

	if err := e.isValidStruct(); err != nil {
		return nil, err
	}

	s := reflect.ValueOf(e.StructAddr).Elem()

	for i := 0; i < s.NumField(); i++ {
		if isIgnored(s.Type().Field(i).Name, e.ignoredFields) {
			continue
		}
		if val, ok := s.Type().Field(i).Tag.Lookup(tag); ok {
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

	for i := 0; i < s.NumField(); i++ {
		if isIgnored(s.Type().Field(i).Name, e.ignoredFields) {
			continue
		}
		if val, ok := s.Type().Field(i).Tag.Lookup(tag); ok {
			if prefix != "" {
				out = append(out, prefix+val)
				continue
			}
			out = append(out, val)
		}
	}

	return
}

//Values returns an interface array with all the values
func (e *Extractor) Values() (out []interface{}, err error) {

	if err := e.isValidStruct(); err != nil {
		return nil, err
	}

	s := reflect.ValueOf(e.StructAddr).Elem()
	for i := 0; i < s.NumField(); i++ {
		if isIgnored(s.Type().Field(i).Name, e.ignoredFields) {
			continue
		}
		out = append(out, s.Field(i).Interface())

	}

	return
}

//ValuesFromTag returns an interface array with all the values of fields with the given tag
func (e *Extractor) ValuesFromTag(tag string) (out []interface{}, err error) {

	if err := e.isValidStruct(); err != nil {
		return nil, err
	}

	s := reflect.ValueOf(e.StructAddr).Elem()
	for i := 0; i < s.NumField(); i++ {
		if isIgnored(s.Type().Field(i).Name, e.ignoredFields) {
			continue
		}
		if _, ok := s.Type().Field(i).Tag.Lookup(tag); ok {
			out = append(out, s.Field(i).Interface())
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
	for i := 0; i < s.NumField(); i++ {
		if isIgnored(s.Type().Field(i).Name, e.ignoredFields) {
			continue
		}
		out[s.Type().Field(i).Name] = s.Field(i).Interface()
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
	for i := 0; i < s.NumField(); i++ {
		if isIgnored(s.Type().Field(i).Name, e.ignoredFields) {
			continue
		}

		if val, ok := s.Type().Field(i).Tag.Lookup(tag); ok {
			out[val] = s.Field(i).Interface()
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

func (e *Extractor) isFieldNameValid(fn string) bool {

	s := reflect.ValueOf(e.StructAddr).Elem()

	for i := 0; i < s.NumField(); i++ {
		if s.Type().Field(i).Name == fn {
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
