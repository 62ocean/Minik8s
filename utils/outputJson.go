package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func OutputJson(name string, data any) error {
	bs, _ := json.Marshal(data)
	var out bytes.Buffer
	err := json.Indent(&out, bs, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%s=%v\n", name, out.String())

	return nil
}
