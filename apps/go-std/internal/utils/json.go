package utils

import "encoding/json"

func PrettyJSON(data interface{}) string {
	pretty, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err.Error()
	}
	return string(pretty)
}
