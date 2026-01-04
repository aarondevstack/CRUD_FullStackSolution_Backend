//go:build linux
// +build linux

package atlas

import _ "embed"

//go:embed linux/atlas-linux-amd64-latest
var atlasBinary []byte
