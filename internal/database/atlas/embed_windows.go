//go:build windows
// +build windows

package atlas

import _ "embed"

//go:embed windows/atlas-windows-amd64-latest.exe
var atlasBinary []byte
