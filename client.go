package pindxru

import (
	"archive/zip"
	"bytes"
	"errors"
	"github.com/LindsayBradford/go-dbf/godbf"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

const (
	rootURL        = "https://www.pochta.ru"
	listUpdatesURL = rootURL + "/database/ops"
	fileEncoding   = "cp866"
)

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
	content := regexp.MustCompile(`(?si)"(Эталонный\\x20справочник\\x20почтовых.+?)"`).FindSubmatch(b)
	if len(content) == 0 {
		err = errors.New("Не найден контент. ")
		return
	}

	matches := regexp.MustCompile(`(?si)(\\x([0-9a-f]{2}))`).FindAllStringSubmatch(string(content[1]), -1)
	processed := map[string]bool{}
	for _, m := range matches {
		var (
			i int64
			e error
		)
		if _, ok := processed[m[1]]; ok {
			continue
		}

		processed[m[2]] = true

		switch m[2] {
		case "a0":
			i = 32

		case "ab", "bb":
			i = 34

		default:
			i, e = strconv.ParseInt(m[2], 16, 16)
			if e != nil {
				log.Fatalln(e)
			}
		}

		content[1] = regexp.MustCompile(`(?si)(\\x`+string(m[2])+")").ReplaceAll(content[1], []byte{uint8(i)})
	}

	re := regexp.MustCompile(`(?si)\|\s*(\d{2}\.\d{2}\.\d{4})\s*\|\s*(\d+)\s*\|\s*\[NPIndx.+?]\((.+?)\).+?(\d+)\s+запи.+?\s*\|\s*\[PIndx.+?]\((.+?)\).+?(\d+)\s+запи.+?\s*\|\s`)
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
	if b, lastMod, err = c.downloadZip(lastRow.Full.Url); err != nil {
		return
	}

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
		err = ioutil.WriteFile(fname, b, perm)
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
		err = ioutil.WriteFile(fname, b, perm)
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
	b, modify, err = c.downloadZip(lastRow.Full.Url)
	return
}

// GetPackageIndexes получает изменения.
func (c Client) GetPackageIndexes(pack *Package) (lastMod time.Time, err error) {
	var b []byte
	if b, lastMod, err = c.downloadZip(pack.Url); err != nil {
		return
	}

	pack.Indexes, err = c.unzipNPIndx(b)
	return lastMod, err
}

// PackageZip загружает zip-файл пакета изменений.
func (c Client) PackageZip(pack Package, filename string, perm os.FileMode) (lastMod time.Time, err error) {
	var b []byte
	if b, lastMod, err = c.downloadZip(pack.Url); err != nil {
		return
	}
	err = ioutil.WriteFile(filename, b, perm)
	return
}

// PackageDbf загружает dbf-файл пакета изменений.
func (c Client) PackageDbf(pack Package, filename string, perm os.FileMode) (lastMod time.Time, err error) {
	var b []byte
	if b, lastMod, err = c.downloadZip(pack.Url); err != nil {
		return
	}

	if b, err = c.unzipDbf(b); err != nil {
		return
	}

	err = ioutil.WriteFile(filename, b, perm)
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
func (c Client) downloadZip(u string) (b []byte, lastMod time.Time, err error) {
	resp, err := c.httpClient.Get(u)
	if err != nil {
		return
	}

	lastMod, err = time.Parse(time.RFC1123, resp.Header.Get("Last-Modified"))
	b, err = getBody(resp)
	return
}

// unzipPIndex распаковывает индексы из zip-файла.
func (c Client) unzipPIndex(file []byte) (indexes []PIndx, err error) {
	file, err = c.unzipDbf(file)
	if err != nil {
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
func (c Client) unzipNPIndx(file []byte) (indexes []NPIndx, err error) {
	file, err = c.unzipDbf(file)
	if err != nil {
		return
	}

	var table *godbf.DbfTable
	if table, err = godbf.NewFromByteArray(file, fileEncoding); err != nil {
		return
	}
	indexes, err = dbfToNPIndx(table)
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

		unzipBytes, err = readZipFile(zipFile)
		if err != nil {
			return nil, err
		}
		break
	}

	return unzipBytes, nil
}
