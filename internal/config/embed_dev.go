//go:build dev
// +build dev

package config

import _ "embed"

//go:embed config_dev.yaml
var embeddedConfig []byte
