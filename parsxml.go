package main

import (
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Text    string   `xml:",chardata"`
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valute  []Val    `xml:"Valute"`
}

type Val struct {
	Text     string `xml:",chardata"`
	ID       string `xml:"ID,attr"`
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  string `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}

//Date
func Date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day,
		0, 0, 0, 0, time.UTC)
}

func main() {
	fday := 1
	month := 10
	dataMap := make(map[string]ValCurs, 100)

	for i := 1; i <= 90; i++ {

		var s ValCurs
		url := "https://www.cbr.ru/scripts/XML_daily_eng.asp?date_req="
		date := "/2020"

		//cобираем url
		url = url + strconv.Itoa(fday) + "/" + strconv.Itoa(month) + date

		//проверка на колличество дней в месяце
		t := Date(2020, month, 0)
		if fday == t.Day() {
			fday = 1
			month++
		}

		//делаем запрос в цб
		res, err := http.Get(url)
		if err != nil { ///
			log.Fatal("can't get '%s': %v", url, err)
		}
		defer res.Body.Close()

		// меняем содировку
		dec := xml.NewDecoder(res.Body)
		dec.CharsetReader = charset.NewReaderLabel

		//заполняем структуру
		if err = dec.Decode(&s); err != nil {
			log.Fatal("can't decode xml body: %v", err)
		}
		//заполняем мапу
		dataMap[s.Date] = s
		fday++
		time.Sleep(time.Second)
	}

	var avgValueRuble, maxValueRuble, minavgValueRuble float64
	var nameMaxValueRuble, nameMinValueRuble, dateMin, dateMax string
	var cnt int

	minavgValueRuble = 100000000000000000.0
	fday = 1
	month = 10

	for y := 1; y <= 90; y++ {

		//собираем ключ
		date := strconv.Itoa(fday) + "." + strconv.Itoa(month) + ".2020"

		t := Date(2020, month, 0)
		if fday == t.Day() {
			fday = 1
			month++
		}

		for i, k := range dataMap[date].Valute {
			tmp, _ := strconv.ParseFloat(strings.Replace(k.Value, ",", ".", 1), 32)
			avgValueRuble += tmp

			if maxValueRuble <= tmp {
				maxValueRuble = tmp
				nameMaxValueRuble = k.Name
				dateMax = date
			}

			if minavgValueRuble >= tmp {
				minavgValueRuble = tmp
				nameMinValueRuble = k.Name
				dateMin = date
			}
			cnt += i
		}
		fday++
	}
	fmt.Printf("AvgValueRuble: %f\n"+
		"MaxValueRuble: %v %f %v\n"+
		"MinavgValueRuble: %v %f %v\n",
		avgValueRuble/(float64(cnt)),
		nameMaxValueRuble, maxValueRuble, dateMax,
		nameMinValueRuble, minavgValueRuble, dateMin,
	)

}
