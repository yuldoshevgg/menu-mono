package helper

import (
	"encoding/json"
	"fmt"
	"log"
)

func MarshalUnMarshal(input, output any) error {
	rawBytes, err := json.Marshal(input)
	if err != nil {
		log.Println("error while marshaling", fmt.Errorf("%w", err))
		return err
	}

	return json.Unmarshal(rawBytes, output)
}
