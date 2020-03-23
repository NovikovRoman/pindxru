package pindxru

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/LindsayBradford/go-dbf/godbf"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var reNumberChanges = regexp.MustCompile(`(?si)(\d+)\sзапис`)

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

func getListUpdates(b []byte) (updates []Updates, err error) {
	utf8, err := charset.NewReader(bytes.NewReader(b), "utf8")
	if err != nil {
		return
	}

	if b, err = ioutil.ReadAll(utf8); err != nil {
		return
	}

	var doc *goquery.Document
	if doc, err = goquery.NewDocumentFromReader(bytes.NewReader(b)); err != nil {
		return
	}

	table := &goquery.Selection{}
	doc.Find("table:last-child tr").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		if i != 1 {
			return true
		}

		selection.Find("table table").EachWithBreak(func(i int, selection *goquery.Selection) bool {
			if i == 2 {
				table = selection
				return false
			}
			return true
		})

		return false
	})

	rows := table.Find("table tr")
	updates = make([]Updates, rows.Length()-1) // тк первая строка заголовочная
	rows.EachWithBreak(func(i int, selection *goquery.Selection) bool {
		if i == 0 {
			return true
		}

		cols := selection.Find("td")
		if cols.Length() != 4 {
			err = fmt.Errorf("Количество столбцов не равено 4. Возможно изменилась верстка страницы %s. ",
				listUpdatesURL)
		}
		index := i - 1

		updates[index] = Updates{}
		cols.Each(func(j int, selection *goquery.Selection) {
			switch j {
			case 0:
				updates[index].Date, _ = time.Parse("02.01.2006", selection.Text())

			case 2:
				updates[index].Url = rootURL + "/" + selection.Find("a").AttrOr("href", "")
				m := reNumberChanges.FindAllStringSubmatch(selection.Text(), 1)
				if len(m) > 0 {
					updates[index].NumberRecords, _ = strconv.Atoi(m[0][1])
				}
			}
		})

		return true
	})

	return
}
