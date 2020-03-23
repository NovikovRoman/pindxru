package pindxru

var cTest *Client

const (
	testdata    = "testdata"
	testZipFile = "test.zip"
	testDbfFile = "test.dbf"
)

func init() {
	cTest = NewClient(nil)
}
