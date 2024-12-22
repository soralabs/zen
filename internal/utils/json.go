package utils

import (
	"encoding/json"
	"regexp"
)

var markdownCodeBlockRegex = regexp.MustCompile("(?s)^```[a-z]*\n?(.*?)\n?```$")

// SmartUnmarshal removes Markdown code fences from the input data (if present)
// and then unmarshals the remaining JSON into `v`.
func SmartUnmarshal(data []byte, v interface{}) error {
	if len(data) == 0 {
		return nil
	}

	strData := string(data)
	if matches := markdownCodeBlockRegex.FindStringSubmatch(strData); len(matches) > 1 {
		strData = matches[1]
	}

	return json.Unmarshal([]byte(strData), v)
}
