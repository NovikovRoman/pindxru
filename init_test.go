package pindxru

import (
	"log"
	"time"
)

var (
	cTest        *Client
	testPackages []Package
)

const (
	testdata    = "testdata"
	testZipFile = "test.zip"
	testDbfFile = "test.dbf"
)

func init() {
	var err error

	cTest = NewClient(nil)

	d := time.Date(2018, 01, 01, 0, 0, 0, 0, time.UTC)
	if testPackages, err = cTest.GetPackages(&d); err != nil {
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
