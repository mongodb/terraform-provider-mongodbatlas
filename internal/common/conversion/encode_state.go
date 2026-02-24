package conversion

import (
	"encoding/base64"
	"fmt"
	"log"
	"sort"
	"strings"
)

func GetEncodedID(stateID, keyPosition string) string {
	id := ""
	if !hasMultipleValues(stateID) {
		return stateID
	}

	decoded := DecodeStateID(stateID)
	id = decoded[keyPosition]

	return id
}

func EncodeStateID(values map[string]string) string {
	encode := func(e string) string { return base64.StdEncoding.EncodeToString([]byte(e)) }
	encodedValues := make([]string, 0)

	// sort to make sure the same encoding is returned in case of same input
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		encodedValues = append(encodedValues, fmt.Sprintf("%s:%s", encode(key), encode(values[key])))
	}

	return strings.Join(encodedValues, "-")
}

func DecodeStateID(stateID string) map[string]string {
	decode := func(d string) string {
		decodedString, err := base64.StdEncoding.DecodeString(d)
		if err != nil {
			log.Printf("[WARN] error decoding state ID: %s", err)
		}

		return string(decodedString)
	}
	decodedValues := make(map[string]string)
	encodedValues := strings.SplitSeq(stateID, "-")

	for value := range encodedValues {
		keyValue := strings.Split(value, ":")
		if len(keyValue) > 1 {
			decodedValues[decode(keyValue[0])] = decode(keyValue[1])
		}
	}

	return decodedValues
}

func hasMultipleValues(value string) bool {
	if strings.Contains(value, "-") && strings.Contains(value, ":") {
		return true
	}

	return false
}
