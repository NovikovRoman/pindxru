package pindxru

import (
	"archive/zip"
	"bytes"
	"github.com/LindsayBradford/go-dbf/godbf"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"
)

const (
	rootURL        = "https://vinfo.russianpost.ru/database"
	fullZipURL     = rootURL + "/PIndx.zip"
	listUpdatesURL = rootURL + "/ops.html"
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

// GetLastModified возвращает дату последнего обновления из web-справочника.
func (c Client) GetLastModified() (lastMod time.Time, err error) {
	var resp *http.Response
	if resp, err = c.httpClient.Head(fullZipURL); err != nil {
		return
	}

	defer func() {
		if derr := resp.Body.Close(); derr != nil {
			err = derr
		}
	}()

	lastMod, err = time.Parse(time.RFC1123, resp.Header.Get("Last-Modified"))
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

	if b, lastMod, err = c.downloadZip(fullZipURL); err != nil {
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

	ok = true
	var b []byte
	if b, modify, err = c.downloadZip(fullZipURL); err != nil {
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

	ok = true
	var b []byte
	if b, modify, err = c.downloadZip(fullZipURL); err != nil {
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
	resp, err := c.httpClient.Get(listUpdatesURL)
	if err != nil {
		return
	}

	body, err := getBody(resp)
	if err != nil {
		return
	}

	l, err := getListUpdates(body)
	if err != nil {
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
