package object

type IndexerType = string

const (
	IndexerGolang       = "GOLANG"
	IndexerJavaJUnit    = "JUNIT"
	IndexerPythonPytest = "PYTEST"
)

// key: case signature
// value: influenced methods
type CaseTagCache = map[string]map[string]interface{}
