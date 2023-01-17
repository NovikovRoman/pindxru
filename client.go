package pindxru

import (
	"archive/zip"
	"bytes"
	"errors"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/NovikovRoman/godbf"
	"golang.org/x/text/encoding/charmap"
)

const (
	rootURL        = "https://www.pochta.ru"
	listUpdatesURL = rootURL + "/support/database/ops"
)

var fileEncoding = charmap.CodePage866

// Client structure.
type Client struct {
	httpClient *http.Client
	transport  *http.Transport
}

// NewClient create new pindxru Client.
func NewClient(transport *http.Transport) *Client {
	c := &http.Client{}
	if transport != nil {
		c.Transport = transport
	}

	return &Client{
		httpClient: c,
		transport:  transport,
	}
}

func (c *Client) GetReferenceRows() (referenceRows ReferenceRows, err error) {
	var b []byte
	if b, err = c.loadPage(); err != nil {
		return
	}

	referenceRows, err = c.parseReferenceRows(b)
	return
}

func (c *Client) parseReferenceRows(b []byte) (referenceRows ReferenceRows, err error) {
	content := regexp.MustCompile(`(?si)<table[^>]*>.+?<td[^>]*>Обновленный эталонный справочник.+?</tr>(.+?)</table>`).FindSubmatch(b)
	if len(content) == 0 {
		err = errors.New("Не найден контент. ")
		return
	}

	re := regexp.MustCompile(`(?si)<tr[^>]*><td[^>]*>\s*(\d{2}\.\d{2}\.\d{4})\s*.+?<td[^>]*>\s*(\d+)\s*.+?<td[^>]*>.+?href="(.+?)".+?</a>.+?(\d+)(?:\s+|&nbsp;| )запис.+?href="(.+?)".+?</a>.+?(\d+)(?:\s+|&nbsp;| )запис`)
	rows := re.FindAllSubmatch(content[1], -1)
	referenceRows = make([]ReferenceRow, len(rows))
	for i, r := range rows {
		referenceRows[i] = ReferenceRow{
			Number: string(r[2]),
			Update: ReferenceFile{
				Url: rootURL + string(r[3]),
			},
			Full: ReferenceFile{
				Url: rootURL + string(r[5]),
			},
		}
		referenceRows[i].Date, _ = time.Parse("02.01.2006", string(r[1]))
		referenceRows[i].Update.Records, _ = strconv.Atoi(string(r[4]))
		referenceRows[i].Full.Records, _ = strconv.Atoi(string(r[6]))
	}
	return
}

// Indexes Возвращает все почтовые индексы из web-справочника.
func (c *Client) Indexes(referenceRows ReferenceRows, lastModified *time.Time) (indexes []PIndx, lastMod time.Time, err error) {
	var (
		b  []byte
		ok bool
	)

	if len(referenceRows) == 0 {
		return
	}

	if lastModified != nil {
		if ok, err = referenceRows.hasUpdates(*lastModified); err != nil || !ok {
			return
		}
	}

	lastRow, _ := referenceRows.LastRow()
	if b, err = c.downloadZip(lastRow.Full.Url); err != nil {
		return
	}

	lastMod = lastRow.Date
	indexes, err = c.unzipPIndex(b)
	return
}

// IndexesZip Загружает zip-файл со всеми почтовыми индексами.
func (c Client) IndexesZip(referenceRows ReferenceRows, fname string, perm os.FileMode, lastMod *time.Time) (modify time.Time, ok bool, err error) {
	var b []byte
	if b, modify, ok, err = c.getFullZip(referenceRows, lastMod); err != nil || !ok {
		return
	}

	if len(b) > 0 {
		err = os.WriteFile(fname, b, perm)
	}

	ok = err == nil
	return
}

// IndexesDbf Загружает dbf-файл со всеми почтовыми индексами.
func (c Client) IndexesDbf(referenceRows ReferenceRows, fname string, perm os.FileMode, lastMod *time.Time) (modify time.Time, ok bool, err error) {
	var b []byte
	if b, modify, ok, err = c.getFullZip(referenceRows, lastMod); err != nil || !ok {
		return
	}

	if b, err = c.unzipDbf(b); err != nil {
		return
	}

	if len(b) > 0 {
		err = os.WriteFile(fname, b, perm)
	}

	ok = err == nil
	return
}

// getFullZip Возвращает последнее полное обновление.
//
// Если не указана lastMod, то самая последняя запись.
//
// Если lastMod указана, то если есть запись после указаной даты.
func (c *Client) getFullZip(referenceRows ReferenceRows, lastMod *time.Time) (b []byte, modify time.Time, ok bool, err error) {
	if len(referenceRows) == 0 {
		return
	}

	if lastMod != nil {
		if ok, err = referenceRows.hasUpdates(*lastMod); err != nil || !ok {
			return
		}
	}

	ok = true
	lastRow, _ := referenceRows.LastRow()
	modify = lastRow.Date
	b, err = c.downloadZip(lastRow.Full.Url)
	return
}

// GetPackageIndexes получает изменения.
func (c Client) GetPackageIndexes(pack *Package) (lastMod time.Time, err error) {
	var b []byte
	if b, err = c.downloadZip(pack.Url); err != nil {
		return
	}

	pack.Indexes, lastMod, err = c.unzipNPIndx(b)
	return
}

// PackageZip загружает zip-файл пакета изменений.
func (c Client) PackageZip(pack Package, filename string, perm os.FileMode) (err error) {
	var b []byte
	if b, err = c.downloadZip(pack.Url); err != nil {
		return
	}
	err = os.WriteFile(filename, b, perm)
	return
}

// PackageDbf загружает dbf-файл пакета изменений.
func (c Client) PackageDbf(pack Package, filename string, perm os.FileMode) (err error) {
	var b []byte
	if b, err = c.downloadZip(pack.Url); err != nil {
		return
	}

	if b, err = c.unzipDbf(b); err != nil {
		return
	}

	err = os.WriteFile(filename, b, perm)
	return
}

func (c *Client) loadPage() (b []byte, err error) {
	var (
		resp *http.Response
	)

	if resp, err = c.httpClient.Get(listUpdatesURL); err != nil {
		return
	}
	b, err = getBody(resp)
	return
}

// downloadZip Загружает zip-файл из web-справочника.
func (c Client) downloadZip(u string) (b []byte, err error) {
	var resp *http.Response
	if resp, err = c.httpClient.Get(u); err != nil {
		return
	}

	b, err = getBody(resp)
	return
}

// unzipPIndex распаковывает индексы из zip-файла.
func (c Client) unzipPIndex(file []byte) (indexes []PIndx, err error) {
	if file, err = c.unzipDbf(file); err != nil {
		return
	}

	var table *godbf.DbfTable
	if table, err = godbf.NewFromByteArray(file, fileEncoding); err != nil {
		return
	}
	indexes, err = dbfToPIndx(table)
	return
}

// unzipNPIndx распаковывает индексы из zip-файла.
func (c Client) unzipNPIndx(file []byte) (indexes []NPIndx, lastMod time.Time, err error) {
	file, err = c.unzipDbf(file)
	if err != nil {
		return
	}

	var table *godbf.DbfTable
	if table, err = godbf.NewFromByteArray(file, fileEncoding); err != nil {
		return
	}

	if indexes, err = dbfToNPIndx(table); err != nil {
		return
	}

	for _, i := range indexes {
		if i.UpdatedAt.After(lastMod) {
			lastMod = i.UpdatedAt
		}
	}
	return
}

// unzipDbf dbf-файл из zip-файла, который содержит dbf-файл с именем `PIndx[N].dbf`, где N - целое число.
func (c Client) unzipDbf(b []byte) (unzipBytes []byte, err error) {
	var zipReader *zip.Reader
	if zipReader, err = zip.NewReader(bytes.NewReader(b), int64(len(b))); err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`(?si)^(PIndx|NPIndx)\d*\.dbf$`)
	for _, zipFile := range zipReader.File {
		if !re.MatchString(zipFile.Name) {
			continue
		}

		if unzipBytes, err = readZipFile(zipFile); err != nil {
			return
		}
		break
	}

	return
}
