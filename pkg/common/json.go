package common

import "encoding/json"

func ToMarshallString(v interface{}) string {
	b, _ := json.Marshal(v)

	return string(b)
}

func ToMarshallStringWithIndent(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "\t")

	return string(b)
}
