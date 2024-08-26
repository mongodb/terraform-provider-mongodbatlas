package schemafunc

import (
	"encoding/json"
	"log"
	"reflect"
)

func EqualJSON(old, newStr, errContext string) bool {
	var j, j2 any

	if old == "" {
		old = "{}"
	}

	if newStr == "" {
		newStr = "{}"
	}
	if err := json.Unmarshal([]byte(old), &j); err != nil {
		log.Printf("[ERROR] cannot unmarshal old %s json %v", errContext, err)
		return false
	}
	if err := json.Unmarshal([]byte(newStr), &j2); err != nil {
		log.Printf("[ERROR] cannot unmarshal new %s json %v", errContext, err)
		return false
	}
	return reflect.DeepEqual(&j, &j2)
}
