package schema

import _ "embed"

//go:embed schema_1.ddl
var PrimarySchema1 string

//go:embed schema_2.ddl
var PrimarySchema2 string

//go:embed schema_3.ddl
var PrimarySchema3 string
