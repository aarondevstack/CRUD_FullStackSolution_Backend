package rbac

import _ "embed"

//go:embed model.conf
var ModelConf []byte

//go:embed policy.csv
var PolicyCSV []byte
