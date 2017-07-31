package types

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/ugorji/go/codec"
	yaml "gopkg.in/yaml.v2"
)

func unmarshalObj(obj map[string]interface{}, in string, f func([]byte, interface{}) error) (map[string]interface{}, error) {
	err := f([]byte(in), &obj)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to unmarshal object %s: %v", in, err))
	}

	return obj, nil
}

func unmarshalArray(obj []interface{}, in string, f func([]byte, interface{}) error) ([]interface{}, error) {
	err := f([]byte(in), &obj)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to unmarshal array %s: %v", in, err))
	}

	return obj, nil
}

func marshalObj(obj interface{}, f func(interface{}) ([]byte, error)) (string, error) {
	b, err := f(obj)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Unable to marshal object %s: %v", obj, err))
	}

	return string(b), nil
}

func toJSONBytes(in interface{}) ([]byte, error) {
	h := &codec.JsonHandle{}
	h.Canonical = true
	buf := new(bytes.Buffer)
	err := codec.NewEncoder(buf, h).Encode(in)
	if err != nil {
		return []byte{}, errors.New(fmt.Sprintf("Unable to marshal %s: %v", in, err))
	}

	return buf.Bytes(), nil
}

func JSON(in string) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	return unmarshalObj(obj, in, yaml.Unmarshal)
}

func JSONArray(in string) ([]interface{}, error) {
	obj := make([]interface{}, 1)
	return unmarshalArray(obj, in, yaml.Unmarshal)
}

func ToJSON(in interface{}) (string, error) {
	b, err := toJSONBytes(in)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func ToYAML(in interface{}) (string, error) {
	return marshalObj(in, yaml.Marshal)
}

func YAML(in string) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	return unmarshalObj(obj, in, yaml.Unmarshal)
}

func YAMLArray(in string) ([]interface{}, error) {
	obj := make([]interface{}, 1)
	return unmarshalArray(obj, in, yaml.Unmarshal)
}

func TOML(in string) (interface{}, error) {
	obj := make(map[string]interface{})
	return unmarshalObj(obj, in, toml.Unmarshal)
}

func ToTOML(in interface{}) (string, error) {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(in)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Unable to marshal %s: %v", in, err))
	}

	return string(buf.Bytes()), nil
}
