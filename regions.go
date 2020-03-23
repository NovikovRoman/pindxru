package pindxru

import "strings"

type regions []struct {
	Name     string
	Code     int
	Autonomy bool
}

func (r regions) GetCode(region string) (code int, autonomy bool) {
	if region == "" {
		return
	}
	region = strings.ToUpper(region)
	for _, item := range r {
		if item.Name == region {
			return item.Code, item.Autonomy
		}
	}
	return
}

// GetName возвращает название региона по его коду
func (r regions) GetName(code int) (name string) {
	if code <= 0 {
		return
	}
	for _, item := range r {
		if item.Code == code {
			return item.Name
		}
	}
	return
}

func (r regions) FindCode(region string, autonomy string) int {
	if region == "" {
		region = autonomy
	}
	code, _ := Regions.GetCode(region)
	return code
}

// Список регионов
var Regions = regions{
	{Name: "МОСКВА", Code: 77,},
	{Name: "МОСКОВСКАЯ ОБЛАСТЬ", Code: 50,},
	{Name: "АДЫГЕЯ РЕСПУБЛИКА", Code: 1},
	{Name: "АЛТАЙ РЕСПУБЛИКА", Code: 4,},
	{Name: "АЛТАЙСКИЙ КРАЙ", Code: 22,},
	{Name: "АМУРСКАЯ ОБЛАСТЬ", Code: 28,},
	{Name: "АРХАНГЕЛЬСКАЯ ОБЛАСТЬ", Code: 29,},
	{Name: "АСТРАХАНСКАЯ ОБЛАСТЬ", Code: 30,},
	{Name: "БАШКОРТОСТАН РЕСПУБЛИКА", Code: 2,},
	{Name: "БЕЛГОРОДСКАЯ ОБЛАСТЬ", Code: 31,},
	{Name: "БРЯНСКАЯ ОБЛАСТЬ", Code: 32,},
	{Name: "БУРЯТИЯ РЕСПУБЛИКА", Code: 3,},
	{Name: "ВЛАДИМИРСКАЯ ОБЛАСТЬ", Code: 33,},
	{Name: "ВОЛГОГРАДСКАЯ ОБЛАСТЬ", Code: 34,},
	{Name: "ВОЛОГОДСКАЯ ОБЛАСТЬ", Code: 35,},
	{Name: "ВОРОНЕЖСКАЯ ОБЛАСТЬ", Code: 36,},
	{Name: "ДАГЕСТАН РЕСПУБЛИКА", Code: 5,},
	{Name: "ЗАБАЙКАЛЬСКИЙ КРАЙ", Code: 75,},
	{Name: "ИВАНОВСКАЯ ОБЛАСТЬ", Code: 37,},
	{Name: "ИНГУШЕТИЯ РЕСПУБЛИКА", Code: 6,},
	{Name: "ИРКУТСКАЯ ОБЛАСТЬ", Code: 38,},
	{Name: "КАБАРДИНО-БАЛКАРСКАЯ РЕСПУБЛИКА", Code: 7,},
	{Name: "КАЛИНИНГРАДСКАЯ ОБЛАСТЬ", Code: 39,},
	{Name: "КАЛМЫКИЯ РЕСПУБЛИКА", Code: 8,},
	{Name: "КАЛУЖСКАЯ ОБЛАСТЬ", Code: 40,},
	{Name: "КАМЧАТСКИЙ КРАЙ", Code: 41,},
	{Name: "КАРАЧАЕВО-ЧЕРКЕССКАЯ РЕСПУБЛИКА", Code: 9,},
	{Name: "КАРЕЛИЯ РЕСПУБЛИКА", Code: 10,},
	{Name: "КЕМЕРОВСКАЯ ОБЛАСТЬ", Code: 42,},
	{Name: "КИРОВСКАЯ ОБЛАСТЬ", Code: 43,},
	{Name: "КОМИ РЕСПУБЛИКА", Code: 11,},
	{Name: "КОСТРОМСКАЯ ОБЛАСТЬ", Code: 44,},
	{Name: "КРАСНОДАРСКИЙ КРАЙ", Code: 23,},
	{Name: "КРАСНОЯРСКИЙ КРАЙ", Code: 24,},
	{Name: "КРЫМ РЕСПУБЛИКА", Code: 82,},
	{Name: "КУРГАНСКАЯ ОБЛАСТЬ", Code: 45,},
	{Name: "КУРСКАЯ ОБЛАСТЬ", Code: 46,},
	{Name: "ЛЕНИНГРАДСКАЯ ОБЛАСТЬ", Code: 47,},
	{Name: "ЛИПЕЦКАЯ ОБЛАСТЬ", Code: 48,},
	{Name: "МАГАДАНСКАЯ ОБЛАСТЬ", Code: 49,},
	{Name: "МАРИЙ ЭЛ РЕСПУБЛИКА", Code: 12,},
	{Name: "МОРДОВИЯ РЕСПУБЛИКА", Code: 13,},
	{Name: "МУРМАНСКАЯ ОБЛАСТЬ", Code: 51,},
	{Name: "НИЖЕГОРОДСКАЯ ОБЛАСТЬ", Code: 52,},
	{Name: "НОВГОРОДСКАЯ ОБЛАСТЬ", Code: 53,},
	{Name: "НОВОСИБИРСКАЯ ОБЛАСТЬ", Code: 54,},
	{Name: "ОМСКАЯ ОБЛАСТЬ", Code: 55,},
	{Name: "ОРЕНБУРГСКАЯ ОБЛАСТЬ", Code: 56,},
	{Name: "ОРЛОВСКАЯ ОБЛАСТЬ", Code: 57,},
	{Name: "ПЕНЗЕНСКАЯ ОБЛАСТЬ", Code: 58,},
	{Name: "ПЕРМСКИЙ КРАЙ", Code: 59,},
	{Name: "ПРИМОРСКИЙ КРАЙ", Code: 25,},
	{Name: "ПСКОВСКАЯ ОБЛАСТЬ", Code: 60,},
	{Name: "РОСТОВСКАЯ ОБЛАСТЬ", Code: 61,},
	{Name: "РЯЗАНСКАЯ ОБЛАСТЬ", Code: 62,},
	{Name: "САМАРСКАЯ ОБЛАСТЬ", Code: 63,},
	{Name: "САНКТ-ПЕТЕРБУРГ", Code: 78,},
	{Name: "САРАТОВСКАЯ ОБЛАСТЬ", Code: 64,},
	{Name: "САХА (ЯКУТИЯ) РЕСПУБЛИКА", Code: 14,},
	{Name: "САХАЛИНСКАЯ ОБЛАСТЬ", Code: 65,},
	{Name: "СВЕРДЛОВСКАЯ ОБЛАСТЬ", Code: 66,},
	{Name: "СЕВАСТОПОЛЬ", Code: 92,},
	{Name: "СЕВЕРНАЯ ОСЕТИЯ - АЛАНИЯ РЕСПУБЛИКА", Code: 15,},
	{Name: "СМОЛЕНСКАЯ ОБЛАСТЬ", Code: 67,},
	{Name: "СТАВРОПОЛЬСКИЙ КРАЙ", Code: 26,},
	{Name: "ТАМБОВСКАЯ ОБЛАСТЬ", Code: 68,},
	{Name: "ТАТАРСТАН РЕСПУБЛИКА", Code: 16,},
	{Name: "ТВЕРСКАЯ ОБЛАСТЬ", Code: 69,},
	{Name: "ТОМСКАЯ ОБЛАСТЬ", Code: 70,},
	{Name: "ТУЛЬСКАЯ ОБЛАСТЬ", Code: 71,},
	{Name: "ТЫВА РЕСПУБЛИКА", Code: 17,},
	{Name: "ТЮМЕНСКАЯ ОБЛАСТЬ", Code: 72,},
	{Name: "УДМУРТСКАЯ РЕСПУБЛИКА", Code: 18,},
	{Name: "ХАБАРОВСКИЙ КРАЙ", Code: 27,},
	{Name: "ХАКАСИЯ РЕСПУБЛИКА", Code: 19,},
	{Name: "ЧЕЛЯБИНСКАЯ ОБЛАСТЬ", Code: 74,},
	{Name: "ЧЕЧЕНСКАЯ РЕСПУБЛИКА", Code: 20,},
	{Name: "ЧУВАШИЯ РЕСПУБЛИКА", Code: 21,},
	{Name: "ЯРОСЛАВСКАЯ ОБЛАСТЬ", Code: 76,},
	{Name: "УЛЬЯНОВСКАЯ ОБЛАСТЬ", Code: 73,},
	{Name: "НЕНЕЦКИЙ АВТОНОМНЫЙ ОКРУГ", Code: 83, Autonomy: true,},
	{Name: "ХАНТЫ-МАНСИЙСКИЙ-ЮГРА АВТОНОМНЫЙ ОКРУГ", Code: 86, Autonomy: true,},
	{Name: "ЯМАЛО-НЕНЕЦКИЙ АВТОНОМНЫЙ ОКРУГ", Code: 89, Autonomy: true,},
	{Name: "ЕВРЕЙСКАЯ АВТОНОМНАЯ ОБЛАСТЬ", Code: 79, Autonomy: true,},
	{Name: "ЧУКОТСКИЙ АВТОНОМНЫЙ ОКРУГ", Code: 87, Autonomy: true,},
}
