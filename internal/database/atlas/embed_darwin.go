//go:build darwin
// +build darwin

package atlas

import _ "embed"

//go:embed darwin/atlas-darwin-arm64-latest
var atlasBinary []byte
