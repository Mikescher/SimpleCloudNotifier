package schema

import _ "embed"

type Def struct {
	SQL  string
	Hash string
}

//go:embed primary_1.ddl
var primarySchema1 string

//go:embed primary_2.ddl
var primarySchema2 string

//go:embed primary_3.ddl
var primarySchema3 string

//go:embed primary_4.ddl
var primarySchema4 string

//go:embed primary_5.ddl
var primarySchema5 string

//go:embed primary_migration_3_4.ddl
var PrimaryMigration_3_4 string

//go:embed primary_migration_4_5.ddl
var PrimaryMigration_4_5 string

//go:embed requests_1.ddl
var requestsSchema1 string

//go:embed logs_1.ddl
var logsSchema1 string

var PrimarySchema = map[int]Def{
	0: {"", ""},
	1: {primarySchema1, "f2b2847f32681a7178e405553beea4a324034915a0c5a5dc70b3c6abbcc852f2"},
	2: {primarySchema2, "07ed1449114416ed043084a30e0722a5f97bf172161338d2f7106a8dfd387d0a"},
	3: {primarySchema3, "65c2125ad0e12d02490cf2275f0067ef3c62a8522edf9a35ee8aa3f3c09b12e8"},
	4: {primarySchema4, "cb022156ab0e7aea39dd0c985428c43cae7d60e41ca8e9e5a84c774b3019d2ca"},
	5: {primarySchema5, "04bd0d4a81540f69f10c8f8cd656a1fdf852d4ef7a2ab2918ca6369b5423b1b6"},
}

var PrimarySchemaVersion = 5

var RequestsSchema = map[int]Def{
	0: {"", ""},
	1: {requestsSchema1, "ebb0a5748b605e8215437413b738279670190ca8159b6227cfc2aa13418b41e9"},
}

var RequestsSchemaVersion = 1

var LogsSchema = map[int]Def{
	0: {"", ""},
	1: {logsSchema1, "65fba477c04095effc3a8e1bb79fe7547b8e52e983f776f156266eddc4f201d7"},
}

var LogsSchemaVersion = 1
