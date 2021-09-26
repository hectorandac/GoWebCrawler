package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var reader = bufio.NewReader(os.Stdin)

type WeatherStatus struct {
	station            string
	weather            string
	visibility         string
	cloud              string
	temperature        string
	windDirection      string
	windSpeed          string
	humidity           string
	dailyPresipitation string
	airPressure        string
}

func (w WeatherStatus) String() string {
	text := fmt.Sprintf("%s:\n", w.station)
	if w.weather != "" {
		text += fmt.Sprintf("\tWeather: %s\n", w.weather)
	}
	if w.visibility != "" {
		text += fmt.Sprintf("\tVisibility: %skm\n", w.visibility)
	}
	if w.cloud != "" {
		text += fmt.Sprintf("\tCloud: %s/10\n", w.cloud)
	}
	if w.temperature != "" {
		text += fmt.Sprintf("\tTemperature: %s°C\n", w.temperature)
	}
	if w.windDirection != "" {
		text += fmt.Sprintf("\tWind direction: %s\n", w.windDirection)
	}
	if w.windSpeed != "" {
		text += fmt.Sprintf("\tWind speed: %sm/s\n", w.windSpeed)
	}
	if w.humidity != "" {
		text += fmt.Sprintf("\tHumidity: %s%%\n", w.humidity)
	}
	if w.dailyPresipitation != "" {
		text += fmt.Sprintf("\tPresipitation: %smm\n", w.dailyPresipitation)
	}
	if w.airPressure != "" {
		text += fmt.Sprintf("\tPressure: %shPa\n", w.airPressure)
	}

	return text + "----------------------------------------------"
}

func main() {

	weathers := []WeatherStatus{}

	htmlData := fetchBody("https://web.kma.go.kr/eng/weather/forecast/current_korea.jsp")
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlData))

	if err != nil {
		panic(err)
	}

	doc.Find("table").Each(func(index int, tablehtml *goquery.Selection) {
		tablehtml.Find("tr").Each(func(indextr int, rowhtml *goquery.Selection) {
			var row []string
			rowhtml.Find("td").Each(func(i int, s *goquery.Selection) {
				row = append(row, s.Text())
			})
			if len(row) == 10 {
				weathers = append(weathers, WeatherStatus{row[0], row[1], row[2], row[3], row[4], row[5], row[6], row[7], row[8], row[9]})
			}
		})
	})

	fmt.Println("####### number of records = ", len(weathers))
	fmt.Println("Enter the station name")

	for {
		desiredStation := strings.ToLower(captureUserInput())
		for i := range weathers {
			if strings.Contains(strings.ToLower(weathers[i].station), desiredStation) {
				fmt.Println(weathers[i])
			}
		}
	}
}

func captureUserInput() string {
	fmt.Print("-> ")
	text, _ := reader.ReadString('\n')
	return strings.Replace(text, "\n", "", -1)
}

func fetchBody(url string) string {
	fmt.Println("Simple web crawler | Getting temperature ...")
	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s\n", html)
}
