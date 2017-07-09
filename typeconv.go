package main

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/ugorji/go/codec"
	yaml "gopkg.in/yaml.v2"
)

type TypeConv struct {
}

func unmarshalObj(obj map[string]interface{}, in string, f func([]byte, interface{}) error) map[string]interface{} {
	err := f([]byte(in), &obj)
	if err != nil {
		log.Fatalf("Unable to unmarshal object %s: %v", in, err)
	}
	return obj
}

func unmarshalArray(obj []interface{}, in string, f func([]byte, interface{}) error) []interface{} {
	err := f([]byte(in), &obj)
	if err != nil {
		log.Fatalf("Unable to unmarshal array %s: %v", in, err)
	}
	return obj
}

func marshalObj(obj interface{}, f func(interface{}) ([]byte, error)) string {
	b, err := f(obj)
	if err != nil {
		log.Fatalf("Unable to marshal object %s: %v", obj, err)
	}

	return string(b)
}

func toJSONBytes(in interface{}) []byte {
	h := &codec.JsonHandle{}
	h.Canonical = true
	buf := new(bytes.Buffer)
	err := codec.NewEncoder(buf, h).Encode(in)
	if err != nil {
		log.Fatalf("Unable to marshal %s: %v", in, err)
	}
	return buf.Bytes()
}

func (t *TypeConv) JSON(in string) map[string]interface{} {
	obj := make(map[string]interface{})
	return unmarshalObj(obj, in, yaml.Unmarshal)
}

func (t *TypeConv) JSONArray(in string) []interface{} {
	obj := make([]interface{}, 1)
	return unmarshalArray(obj, in, yaml.Unmarshal)
}

func (t *TypeConv) ToJSON(in interface{}) string {
	return string(toJSONBytes(in))
}

func (t *TypeConv) toJSONPretty(indent string, in interface{}) string {
	out := new(bytes.Buffer)
	b := toJSONBytes(in)
	err := json.Indent(out, b, "", indent)
	if err != nil {
		log.Fatalf("Unable to indent JSON %s: %v", b, err)
	}

	return string(out.Bytes())
}

func (t *TypeConv) ToYAML(in interface{}) string {
	return marshalObj(in, yaml.Marshal)
}

func (t *TypeConv) YAML(in string) map[string]interface{} {
	obj := make(map[string]interface{})
	return unmarshalObj(obj, in, yaml.Unmarshal)
}

func (t *TypeConv) YAMLArray(in string) []interface{} {
	obj := make([]interface{}, 1)
	return unmarshalArray(obj, in, yaml.Unmarshal)
}

func (t *TypeConv) TOML(in string) interface{} {
	obj := make(map[string]interface{})
	return unmarshalObj(obj, in, toml.Unmarshal)
}

func (t *TypeConv) ToTOML(in interface{}) string {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(in)
	if err != nil {
		log.Fatalf("Unable to marshal %s: %v", in, err)
	}
	return string(buf.Bytes())
}
