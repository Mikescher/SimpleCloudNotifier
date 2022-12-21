package schema

import _ "embed"

//go:embed schema_1.ddl
var Schema1 string

//go:embed schema_2.ddl
var Schema2 string

//go:embed schema_3.ddl
var Schema3 string
