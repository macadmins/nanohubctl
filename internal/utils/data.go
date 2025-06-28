package utils

import "encoding/json"

func PrettyJsonPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
