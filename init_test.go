package pindxru

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	cTest             *Client
	testPackages      []Package
	testReferenceRows ReferenceRows
)

const (
	testdata    = "testdata"
	testZipFile = "test.zip"
	testDbfFile = "test.dbf"
)

func init() {
	var (
		transport *http.Transport
		err       error
	)

	if transport, err = getTransportTest(); err != nil {
		log.Fatalln(err)
	}

	cTest = NewClient(transport)

	if testReferenceRows, err = cTest.GetReferenceRows(); err != nil {
		log.Fatalln(err)
	}

	d := time.Date(2018, 01, 01, 0, 0, 0, 0, time.UTC)
	if testPackages, err = testReferenceRows.GetUpdatePackages(&d); err != nil {
		log.Fatalln(err)
	}

	if len(testPackages) == 0 {
		log.Fatalln("Нет пакетов изменений.")
	}

	for _, p := range testPackages {
		if p.Date.IsZero() || p.NumberRecords <= 0 {
			log.Fatalln(p, "Ошибочные данные в пакетах.")
		}
	}
}

func getTransportTest() (transport *http.Transport, err error) {
	var (
		u     *url.URL
		proxy string
	)

	if proxy = os.Getenv("PROXY"); proxy == "" {
		return
	}

	if u, err = url.Parse(proxy); err != nil {
		return
	}

	transport = &http.Transport{
		Proxy: http.ProxyURL(u),
	}

	return
}
