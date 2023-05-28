package schema

import _ "embed"

//go:embed primary_1.ddl
var PrimarySchema1 string

const PrimaryHash1 = "f2b2847f32681a7178e405553beea4a324034915a0c5a5dc70b3c6abbcc852f2"

//go:embed primary_2.ddl
var PrimarySchema2 string

const PrimaryHash2 = "07ed1449114416ed043084a30e0722a5f97bf172161338d2f7106a8dfd387d0a"

//go:embed primary_3.ddl
var PrimarySchema3 string

const PrimaryHash3 = "a4851e7953d3423622555cde03d2d0ea2ca367fbe28aa3e363771f1d04bed90a"

//go:embed requests_1.ddl
var RequestsSchema1 string

const RequestsHash1 = "ebb0a5748b605e8215437413b738279670190ca8159b6227cfc2aa13418b41e9"

//go:embed logs_1.ddl
var LogsSchema1 string

const LogsHash1 = "65fba477c04095effc3a8e1bb79fe7547b8e52e983f776f156266eddc4f201d7"
