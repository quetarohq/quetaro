package sqsmsg

import (
	"encoding/json"

	"github.com/pkg/errors"
)

func DecodeId(rawBody string, key string) (string, error) {
	var body map[string]any
	err := json.Unmarshal([]byte(rawBody), &body)

	if err != nil {
		return "", errors.Wrap(err, "failed to unmarshal to JSON")
	}

	rawPayload, ok := body["requestPayload"]

	if !ok {
		return "", errors.New("'requestPayload' not found in JSON")
	}

	payload, ok := rawPayload.(map[string]any)

	if !ok {
		return "", errors.New("unable to cast 'requestPayload' to map[string]any")
	}

	rawId, ok := payload[key]

	if !ok {
		return "", errors.Errorf("'%s' not found in 'requestPayload'", key)
	}

	id, ok := rawId.(string)

	if !ok {
		return "", errors.Errorf("unable to cast '%s' to string", key)
	}

	return id, nil
}
