# PIndxru

> Библиотека для получения обновлений почтовых индексов России из «[Эталонного справочника почтовых индексов объектов почтовой связи](https://vinfo.russianpost.ru/database/ops.html)»

## Содержание

* [PIndxru](#pindxru)
  * [Начало работы](#начало-работы)
  * [Примеры использования](#примеры-использования)
  * [Тесты](#тесты)

## Начало работы

```shell
go get github.com/NovikovRoman/pindxru
```

## Примеры использования

### Получить все индексы:
```go
…
pindxruClient := pindxru.NewClient(nil)
indexes, lastMod, err := pindxruClient.All()
if err != nil {
    log.Fatalln(err)
}

fmt.Println(lastMod, len(indexes), indexes[0].Index, indexes[0].Region)
…
```

### Получить индексы, если были изменения после определенной даты:
```go
…
pindxruClient := pindxru.NewClient(nil)
lastMod := time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)
indexes, lastMod, err := pindxruClient.AllUpdated(lastMod)
if err != nil {
    log.Fatalln(err)
}

if len(indexes) > 0 {
    fmt.Println(lastMod, len(indexes))

} else {
    fmt.Println("Изменений не было.")
}
…
```

### Получить список изменений с определенной даты и получить индексы изменения:
```go
pindxruClient := pindxru.NewClient(nil)

lastMod := time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)
updates, err := pindxruClient.Updates(&lastMod)
if err != nil {
    log.Fatalln(err)
}

if len(updates) > 0 {
    indexes, lastMod, err := pindxruClient.GetNPIndxes(updates[0].Url)
    if err != nil {
        log.Fatalln(err)
    }

    fmt.Println(lastMod, len(indexes))

} else {
    fmt.Println("Изменений не было.")
}
```

## Тесты

```shell
go test -v -race
```