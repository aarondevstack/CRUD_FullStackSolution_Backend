//go:build prod
// +build prod

package config

import _ "embed"

//go:embed config_prod.yaml
var embeddedConfig []byte
