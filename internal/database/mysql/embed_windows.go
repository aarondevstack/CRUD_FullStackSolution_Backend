//go:build windows
// +build windows

package mysql

import _ "embed"

//go:embed assets/windows/mysql-8.4.7-winx64.zip
var mysqlZip []byte
