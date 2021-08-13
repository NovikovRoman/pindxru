package pindxru

import (
	"time"
)

type ReferenceRow struct {
	Date   time.Time
	Number string
	Update ReferenceFile
	Full   ReferenceFile
}

type ReferenceFile struct {
	Url     string
	Records int
}

type ReferenceRows []ReferenceRow

// GetLastModified Возвращает дату последнего обновления из web-справочника.
func (r ReferenceRows) GetLastModified() (lastMod time.Time, err error) {
	if len(r) == 0 {
		return
	}

	lastMod = r[len(r)-1].Date
	return
}

// LastRow Возвращает последнюю строку.
func (r ReferenceRows) LastRow() (referenceRow *ReferenceRow, err error) {
	if len(r) == 0 {
		return
	}

	referenceRow = &r[len(r)-1]
	return
}

// GetUpdatePackages Возвращает список обновлений начиная от даты >= lastModified.
func (r ReferenceRows) GetUpdatePackages(lastModified *time.Time) (packages []Package, err error) {
	packages = []Package{}

	for _, rr := range r {
		if lastModified != nil && rr.Date.Before(*lastModified) {
			continue
		}

		packages = append(packages, Package{
			Date:          rr.Date,
			Url:           rr.Update.Url,
			NumberRecords: rr.Update.Records,
			Indexes:       []NPIndx{},
		})
	}

	return
}

// hasUpdates Есть ли обновление.
func (r ReferenceRows) hasUpdates(lastModified time.Time) (ok bool, err error) {
	var lastMod time.Time

	if lastMod, err = r.GetLastModified(); err != nil {
		return
	}

	ok = !(lastModified.Equal(lastMod) || lastModified.After(lastMod))
	return
}
