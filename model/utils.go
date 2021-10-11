package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"strconv"
	"strings"
)

const (
	OpenApi3 = "openapi"
	OpenApi2 = "swagger"
	AsyncApi = "asyncapi"
)

func ExtractSpecInfo(spec []byte) (*SpecInfo, error) {
	var parsedSpec map[string]interface{}
	specVersion := &SpecInfo{}
	runes := []rune(strings.TrimSpace(string(spec)))
	if runes[0] == '{' && runes[len(runes)-1] == '}' {
		// try JSON
		err := json.Unmarshal(spec, &parsedSpec)
		if err != nil {
			return nil, fmt.Errorf("unable to parse specification: %s", err.Error())
		}
		specVersion.Version = "json"
	} else {
		// try YAML
		err := yaml.Unmarshal(spec, &parsedSpec)
		if err != nil {
			return nil, fmt.Errorf("unable to parse specification: %s", err.Error())
		}
		specVersion.Version = "yaml"
	}

	// check for specific keys
	if parsedSpec[OpenApi3] != nil {
		specVersion.SpecType = OpenApi3
		version, majorVersion := parseVersionTypeData(parsedSpec[OpenApi3])

		// double check for the right version, people mix this up.
		if majorVersion < 3 {
			return nil, errors.New("spec is defined as an openapi spec, but is using a swagger (2.0), or unknown version")
		}
		specVersion.Version = version
	}
	if parsedSpec[OpenApi2] != nil {
		specVersion.SpecType = OpenApi2
		version, majorVersion := parseVersionTypeData(parsedSpec[OpenApi2])

		// I am not certain this edge-case is very frequent, but let's make sure we handle it anyway.
		if majorVersion > 2 {
			return nil, errors.New("spec is defined as a swagger (openapi 2.0) spec, but is an openapi 3 or unkown version")
		}
		specVersion.Version = version
	}
	if parsedSpec[AsyncApi] != nil {
		specVersion.SpecType = AsyncApi
		version, majorVersion := parseVersionTypeData(parsedSpec[AsyncApi])

		// so far there is only 2 as a major release of AsyncAPI
		if majorVersion > 2 {
			return nil, errors.New("spec is defined as asyncapi, but has a major version that is invalid")
		}
		specVersion.Version = version

	}

	if specVersion.SpecType == "" {
		return nil, errors.New("spec type not supported by vaccum, sorry")
	}
	return specVersion, nil
}

func parseVersionTypeData(d interface{}) (string, int) {
	switch d.(type) {
	case int:
		return strconv.Itoa(d.(int)), d.(int)
	case float64:
		return strconv.FormatFloat(d.(float64), 'f', 2, 32), int(d.(float64))
	case bool:
		if d.(bool) {
			return "true", 0
		}
		return "false", 0
	case []string:
		return "multiple versions detected", 0
	}
	r := []rune(strings.TrimSpace(fmt.Sprintf("%v", d)))
	return string(r), int(r[0]) - '0'
}