package pindxru

import (
	"errors"
	"time"
)

// PIndx structure.
type PIndx struct {
	// Почтовый индекс объекта почтовой связи в соответствии с действующей системой индексации
	Index string
	// Наименование объекта почтовой связи
	OpsName string
	// Тип объекта почтовой связи
	OpsType string
	// Индекс вышестоящего по иерархии подчиненности объекта почтовой связи
	OpsSub string
	// Наименование области, края, республики, в которой находится объект почтовой связи
	Region string
	// Наименование автономной области, в которой находится объект почтовой связи
	Autonomy string
	// Наименование района, в котором находится объект почтовой связи
	Area string
	// Наименование населенного пункта, в котором находится объект почтовой связи
	City string
	// Наименование подчиненного населенного пункта, в котором находится объект почтовой связи
	SubCity string
	// Дата актуализации информации об объекте почтовой связи
	UpdatedAt time.Time
	// Почтовый индекс объект почтовой связи до ввода действующей системы индексации
	OldIndex string
	// Код региона
	RegionCode int
}

func createPIndx(data []string) (p PIndx, err error) {
	if len(data) != 11 {
		err = errors.New("Полей должно быть 11. ")
		return
	}

	var updatedAt time.Time
	if updatedAt, err = time.Parse("20060102", data[9]); err != nil {
		return
	}

	p = PIndx{
		Index:     data[0],
		OpsName:   data[1],
		OpsType:   data[2],
		OpsSub:    data[3],
		Region:    data[4],
		Autonomy:  data[5],
		Area:      data[6],
		City:      data[7],
		SubCity:   data[8],
		UpdatedAt: updatedAt,
		OldIndex:  data[10],
	}
	p.RegionCode = Regions.FindCode(p.Region, p.Autonomy)
	return
}

// NPIndx structure для обновления.
type NPIndx struct {
	// Новый почтовый индекс объекта почтовой связи в соответствии с действующей системой индексации
	NewIndex string
	// Почтовый индекс объекта почтовой связи в соответствии с действующей системой индексации
	Index string
	// Наименование объекта почтовой связи
	OpsName string
	// Тип объекта почтовой связи
	OpsType string
	//Индекс вышестоящего по иерархии подчиненности объекта почтовой связи
	OpsSub string
	// Наименование области, края, республики, в которой находится объект почтовой связи
	Region string
	// Наименование автономной области, в которой находится объект почтовой связи
	Autonomy string
	// Наименование района, в котором находится объект почтовой связи
	Area string
	// Наименование населенного пункта, в котором находится объект почтовой связи
	City string
	// Наименование подчиненного населенного пункта, в котором находится объект почтовой связи
	SubCity string
	// Дата актуализации информации об объекте почтовой связи
	UpdatedAt time.Time
	// Почтовый индекс объект почтовой связи до ввода действующей системы индексации
	OldIndex string
	// Код региона
	RegionCode int
}

func createNPIndx(data []string) (p NPIndx, err error) {
	if len(data) != 12 {
		err = errors.New("Полей должно быть 12. ")
		return
	}

	var updatedAt time.Time
	if updatedAt, err = time.Parse("20060102", data[10]); err != nil {
		return
	}

	p = NPIndx{
		Index:     data[0],
		NewIndex:  data[1],
		OpsName:   data[2],
		OpsType:   data[3],
		OpsSub:    data[4],
		Region:    data[5],
		Autonomy:  data[6],
		Area:      data[7],
		City:      data[8],
		SubCity:   data[9],
		UpdatedAt: updatedAt,
		OldIndex:  data[11],
	}
	p.RegionCode = Regions.FindCode(p.Region, p.Autonomy)
	return
}

// Updates structure. Информация о частичном обновлении.
type Updates struct {
	Date          time.Time
	Url           string
	NumberRecords int
}
