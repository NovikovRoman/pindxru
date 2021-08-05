package pindxru

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"github.com/LindsayBradford/go-dbf/godbf"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"
)

const (
	rootURL = "https://www.pochta.ru"
	// fullZipURL     = rootURL + "/documents/10231/6755366698/PIndx.zip/912561e7-9221-49d7-9fc9-5673eb04cd40" //"/PIndx.zip"*/
	listUpdatesURL = rootURL + "/database/ops"
	fileEncoding   = "cp866"
)

// Client structure.
type Client struct {
	httpClient *http.Client
	transport  *http.Transport
	page       []byte
	fullZipURL string
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

// GetLastModified возвращает дату последнего обновления из web-справочника.
func (c Client) GetLastModified() (lastMod time.Time, err error) {
	var packages []Package

	if packages, err = c.GetPackages(nil); err != nil || len(packages) == 0 {
		return
	}

	lastMod = packages[len(packages)-1].Date
	return
}

// Indexes возвращает все почтовые индексы из web-справочника.
func (c Client) Indexes(lastModified *time.Time) (indexes []PIndx, lastMod time.Time, err error) {
	var (
		b  []byte
		ok bool
	)

	if lastModified != nil {
		if ok, err = c.hasUpdates(*lastModified); err != nil || !ok {
			return
		}
	}

	if err = c.getFullZipURL(); err != nil {
		return
	}

	if b, lastMod, err = c.downloadZip(c.fullZipURL); err != nil {
		return
	}

	indexes, err = c.unzipPIndex(b)
	return
}

// IndexesZip загружает zip-файл со всеми почтовыми индексами.
func (c Client) IndexesZip(fname string, perm os.FileMode, lastMod *time.Time) (modify time.Time, ok bool, err error) {
	if lastMod != nil {
		if ok, err = c.hasUpdates(*lastMod); err != nil || !ok {
			return
		}
	}

	if err = c.getFullZipURL(); err != nil {
		return
	}

	ok = true
	var b []byte
	if b, modify, err = c.downloadZip(c.fullZipURL); err != nil {
		return
	}
	err = ioutil.WriteFile(fname, b, perm)
	return
}

// IndexesDbf загружает dbf-файл со всеми почтовыми индексами.
func (c Client) IndexesDbf(fname string, perm os.FileMode, lastMod *time.Time) (modify time.Time, ok bool, err error) {
	if lastMod != nil {
		if ok, err = c.hasUpdates(*lastMod); err != nil || !ok {
			return
		}
	}

	if err = c.getFullZipURL(); err != nil {
		return
	}

	ok = true
	var b []byte
	if b, modify, err = c.downloadZip(c.fullZipURL); err != nil {
		return
	}

	if b, err = c.unzipDbf(b); err != nil {
		return
	}

	err = ioutil.WriteFile(fname, b, perm)
	return
}

// GetPackages возвращает список обновлений начиная от даты >= lastModified.
func (c Client) GetPackages(lastModified *time.Time) (packages []Package, err error) {
	var (
		l []Package
	)

	if err = c.loadPage(); err != nil {
		return
	}

	if l, err = getListUpdates(c.page); err != nil {
		return
	}

	packages = []Package{}
	for _, i := range l {
		if lastModified != nil && i.Date.Before(*lastModified) {
			continue
		}

		packages = append(packages, i)
	}

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

func (c *Client) loadPage() (err error) {
	var (
		resp *http.Response
	)

	if len(c.page) == 0 {
		if resp, err = c.httpClient.Get(listUpdatesURL); err != nil {
			return
		}

		if c.page, err = getBody(resp); err != nil {
			return
		}
	}

	return
}

func (c *Client) getFullZipURL() (err error) {
	if c.fullZipURL != "" {
		return
	}

	if err = c.loadPage(); err != nil {
		return
	}

	m := regexp.MustCompile(`(?si)<a\s+href="([^"]+)"[^>]*>\s*Эталонный\s+справочник\s+почтовых\s+индексов`).
		FindAllSubmatch(c.page, 1)
	if len(m) < 1 {
		err = errors.New("Не найдена ссылка на общий zip-файл. ")
		return
	}

	c.fullZipURL = fmt.Sprintf("%s%s", rootURL, m[0][1])
	return
}

// hasUpdates есть ли обновление.
func (c Client) hasUpdates(lastModified time.Time) (ok bool, err error) {
	var lastMod time.Time
	lastMod, err = c.GetLastModified()
	if err != nil {
		return
	}

	ok = !(lastModified.Equal(lastMod) || lastModified.After(lastMod))
	return
}

// downloadZip загружает zip-файл из web-справочника.
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
