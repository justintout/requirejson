package requiredjson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func Unmarshal(data []byte, v interface{}) error {
	val := reflect.ValueOf(v).Elem()
	t := val.Type()
	if t.Kind() != reflect.Struct {
		return json.Unmarshal(data, v)
	}
	requiredFields := []string{}
	for i := 0; i < t.NumField(); i++ {
		field, ok := t.Field(i).Tag.Lookup("json")
		if !ok {
			continue
		}
		name := t.Field(i).Name
		var required bool
		if name, required = parseRequiredTag(field); required {
			fmt.Printf("name: %s, required: %v\n", name, required)
			requiredFields = append(requiredFields, name)
		}
	}
	fmt.Printf("required fields: %v\n", requiredFields)
	if len(requiredFields) == 0 {
		return json.Unmarshal(data, v)
	}
	req := make(map[string]interface{})
	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}
	fmt.Printf("unmarshaled map: %+v\n", req)
	for _, n := range requiredFields {
		if _, ok := req[n]; !ok {
			return fmt.Errorf("%s is required", n)
		}
	}
	return json.Unmarshal(data, v)
}

func parseRequiredTag(tag string) (name string, required bool) {
	name = ""
	if idx := strings.Index(tag, ","); idx != -1 {
		name = tag[:idx]
		opts := tag[idx+1:]
		for opts != "" {
			fmt.Println(opts)
			var next string
			i := strings.Index(opts, ",")
			if i >= 0 {
				opts, next = opts[:i], opts[i+1:]
			}
			if opts == "required" {
				return name, true
			}
			opts = next
		}
		return name, false
	}
	return "", false
}
