//go:build darwin
// +build darwin

package mysql

import _ "embed"

//go:embed assets/darwin/mysql-8.4.7-macos15-arm64.zip
var mysqlZip []byte
