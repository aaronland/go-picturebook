package multi

import (
	"errors"
	"fmt"
	"strings"
)

type KeyValueFlag struct {
	Key   string
	Value string
}

type KeyValue []*KeyValueFlag

func (e *KeyValue) String() string {
	return fmt.Sprintf("%v", *e)
}

func (e *KeyValue) Set(value string) error {

	value = strings.Trim(value, " ")
	kv := strings.Split(value, "=")

	if len(kv) != 2 {
		return errors.New("Invalid key=value argument")
	}

	a := KeyValueFlag{
		Key:   kv[0],
		Value: kv[1],
	}

	*e = append(*e, &a)
	return nil
}

func (e *KeyValue) Get() interface{} {
	return *e
}

func (e *KeyValue) ToMap() map[string]string {

	m := make(map[string]string)

	for _, arg := range *e {
		m[arg.Key] = arg.Value
	}

	return m
}
