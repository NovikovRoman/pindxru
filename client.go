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

type Client struct {
	*http.Client
}

// New create new pindxru Client.
func NewClient(transport *http.Transport) *Client {
	c := &http.Client{}
	if transport != nil {
		c.Transport = transport
	}

	return &Client{
		Client: c,
	}
}

// GetLastModified возвращает дату последнего обновления из web-справочника.
func (c Client) GetLastModified() (lastMod time.Time, err error) {
	var resp *http.Response
	if resp, err = c.Head(fullZipURL); err != nil {
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

// GetPostIndexes возвращает все почтовые индексы из web-справочника.
func (c Client) All() (indexes []PIndx, lastMod time.Time, err error) {
	var b []byte
	if b, lastMod, err = c.downloadZip(fullZipURL); err != nil {
		return
	}

	indexes, err = c.unzipPIndex(b)
	return
}

// AllUpdates возвращает все почтовые индексы из web-справочника, если есть обновления.
func (c Client) AllUpdates(lastModified time.Time) (indexes []PIndx, lastMod time.Time, err error) {
	var ok bool
	if ok, err = c.isUpdates(lastModified); err != nil || !ok {
		return
	}
	return c.All()
}

// ZipAll загружает zip-файл со всеми почтовыми индексами.
func (c Client) ZipAll(filename string, perm os.FileMode) (lastMod time.Time, err error) {
	var b []byte
	if b, lastMod, err = c.downloadZip(fullZipURL); err != nil {
		return
	}
	err = ioutil.WriteFile(filename, b, perm)
	return
}

// ZipAllUpdates загружает zip-файл со всеми почтовыми индексами, если есть обновления.
func (c Client) ZipAllUpdates(fname string, perm os.FileMode, lastMod time.Time) (lm time.Time, ok bool, err error) {
	if ok, err = c.isUpdates(lastMod); err != nil || !ok {
		return
	}
	lm, err = c.ZipAll(fname, perm)
	ok = err == nil
	return
}

// DbfAll загружает dbf-файл со всеми почтовыми индексами.
func (c Client) DbfAll(fname string, perm os.FileMode) (lastMod time.Time, err error) {
	var b []byte
	if b, lastMod, err = c.downloadZip(fullZipURL); err != nil {
		return
	}

	if b, err = c.unzipDbf(b); err != nil {
		return
	}

	err = ioutil.WriteFile(fname, b, perm)
	return
}

// DbfAllUpdates загружает dbf-файл со всеми почтовыми индексами, если есть обновления.
func (c Client) DbfAllUpdates(fname string, perm os.FileMode, lastMod time.Time) (lm time.Time, ok bool, err error) {
	if ok, err = c.isUpdates(lastMod); err != nil || !ok {
		return
	}
	lm, err = c.DbfAll(fname, perm)
	ok = err == nil
	return
}

// Updates возвращает список дат обновлений начиная от даты >= lastModified.
func (c Client) Updates(lastModified *time.Time) (updates []Updates, err error) {
	resp, err := c.Get(listUpdatesURL)
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

	updates = []Updates{}
	for _, i := range l {
		if lastModified != nil && i.Date.Before(*lastModified) {
			continue
		}

		updates = append(updates, i)
	}

	return
}

// GetNPIndxes возвращает изменения, если есть
func (c Client) GetNPIndxes(u string) (indexes []NPIndx, lastMod time.Time, err error) {
	var b []byte
	if b, lastMod, err = c.downloadZip(u); err != nil {
		return
	}

	indexes, err = c.unzipNPIndx(b)
	return indexes, lastMod, err
}

// ZipNPIndx загружает zip-файл с изменнными почтовыми индексами.
func (c Client) ZipNPIndx(u string, filename string, perm os.FileMode) (lastMod time.Time, err error) {
	var b []byte
	if b, lastMod, err = c.downloadZip(u); err != nil {
		return
	}
	err = ioutil.WriteFile(filename, b, perm)
	return
}

// DbfNPIndx загружает zip-файл с изменнными почтовыми индексами.
func (c Client) DbfNPIndx(u string, filename string, perm os.FileMode) (lastMod time.Time, err error) {
	var b []byte
	if b, lastMod, err = c.downloadZip(u); err != nil {
		return
	}

	if b, err = c.unzipDbf(b); err != nil {
		return
	}

	err = ioutil.WriteFile(filename, b, perm)
	return
}

// isUpdates есть ли обновление.
func (c Client) isUpdates(lastModified time.Time) (ok bool, err error) {
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
	resp, err := c.Get(u)
	if err != nil {
		return
	}

	lastMod, err = time.Parse(time.RFC1123, resp.Header.Get("Last-Modified"))
	b, err = getBody(resp)

	return b, lastMod, nil
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
