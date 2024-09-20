package schema

import (
	"embed"
	_ "embed"
)

type Def struct {
	SQL       string
	Hash      string
	MigScript string
}

//go:embed *.ddl
//go:embed *.sql
var assets embed.FS

var PrimarySchema = map[int]Def{
	0: {"", "", ""},
	1: {readDDL("primary_1.ddl"), "f2b2847f32681a7178e405553beea4a324034915a0c5a5dc70b3c6abbcc852f2", ""},
	2: {readDDL("primary_2.ddl"), "07ed1449114416ed043084a30e0722a5f97bf172161338d2f7106a8dfd387d0a", ""},
	3: {readDDL("primary_3.ddl"), "65c2125ad0e12d02490cf2275f0067ef3c62a8522edf9a35ee8aa3f3c09b12e8", ""},
	4: {readDDL("primary_4.ddl"), "cb022156ab0e7aea39dd0c985428c43cae7d60e41ca8e9e5a84c774b3019d2ca", readMig("primary_migration_3_4.sql")},
	5: {readDDL("primary_5.ddl"), "9d6217ba4a3503cfe090f72569367f95a413bb14e9effe49ffeabbf255bce8dd", readMig("primary_migration_4_5.sql")},
	6: {readDDL("primary_6.ddl"), "8e83d20bcd008082713f248ae8cd558335a37a37ce90bd8c86e782da640ee160", readMig("primary_migration_5_6.sql")},
	7: {readDDL("primary_7.ddl"), "90d8dbc460afe025f9b74cda5c16bb8e58b178df275223bd2531907a8d8c36c3", readMig("primary_migration_6_7.sql")},
	8: {readDDL("primary_8.ddl"), "746f6005c7a573b8816e5993ecd1d949fe2552b0134ba63bab8b4d5b2b5058ad", readMig("primary_migration_7_8.sql")},
}

var PrimarySchemaVersion = len(PrimarySchema) - 1

var RequestsSchema = map[int]Def{
	0: {"", "", ""},
	1: {readDDL("requests_1.ddl"), "ebb0a5748b605e8215437413b738279670190ca8159b6227cfc2aa13418b41e9", ""},
}

var RequestsSchemaVersion = len(RequestsSchema) - 1

var LogsSchema = map[int]Def{
	0: {"", "", ""},
	1: {readDDL("logs_1.ddl"), "65fba477c04095effc3a8e1bb79fe7547b8e52e983f776f156266eddc4f201d7", ""},
}

var LogsSchemaVersion = len(LogsSchema) - 1

func readDDL(name string) string {
	data, err := assets.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func readMig(name string) string {
	data, err := assets.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return string(data)
}
