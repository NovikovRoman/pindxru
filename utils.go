package pindxru

import (
	"archive/zip"
	"fmt"
	"github.com/LindsayBradford/go-dbf/godbf"
	"io/ioutil"
	"net/http"
)

func dbfToPIndx(table *godbf.DbfTable) ([]PIndx, error) {
	postIndexes := make([]PIndx, table.NumberOfRecords())

	for row := 0; row < table.NumberOfRecords(); row++ {
		p, err := createPIndx(table.GetRowAsSlice(row))
		if err != nil {
			return nil, fmt.Errorf("row %d: %s", row, err)
		}
		postIndexes[row] = p
	}

	return postIndexes, nil
}

func dbfToNPIndx(table *godbf.DbfTable) ([]NPIndx, error) {
	postIndexes := make([]NPIndx, table.NumberOfRecords())

	for row := 0; row < table.NumberOfRecords(); row++ {
		p, err := createNPIndx(table.GetRowAsSlice(row))
		if err != nil {
			return nil, fmt.Errorf("row %d: %s", row, err)
		}
		postIndexes[row] = p
	}

	return postIndexes, nil
}

func readZipFile(zf *zip.File) (body []byte, error error) {
	f, err := zf.Open()
	if err != nil {
		return
	}

	defer func() {
		if derr := f.Close(); derr != nil {
			err = derr
		}
	}()

	body, err = ioutil.ReadAll(f)
	return
}

func getBody(resp *http.Response) (body []byte, err error) {
	defer func() {
		if derr := resp.Body.Close(); derr != nil {
			err = derr
		}
	}()

	body, err = ioutil.ReadAll(resp.Body)
	return
}
