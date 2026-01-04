//go:build linux
// +build linux

package mysql

import _ "embed"

//go:embed assets/linux/mysql-8.4.7-linux-x86_64.zip
var mysqlZip []byte
