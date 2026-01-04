//go:build !prod
// +build !prod

package config

import _ "embed"

//go:embed config_dev.yaml
var embeddedConfig []byte
