package main

import (
	controller "cowin-alert/controllers"
	"cowin-alert/database"
	"cowin-alert/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	BASE_COWIN_URL = "https://cdn-api.co-vin.in/api/v2/appointment/sessions/public/findByDistrict"
	PKD = 308
)

func main() {
	database.Connection()
	today_date := time.Now().Format("02-01-2006")
	tomorrowDate := time.Now().Add(time.Hour * 24).Format("02-01-2006")
	dayAfterTomorrow := time.Now().Add(time.Hour * 48).Format("02-01-2006")

	queryParamsToday := fmt.Sprintf("?district_id=%d&date=%s",PKD,today_date)
	queryParamsTomorrow := fmt.Sprintf("?district_id=%d&date=%s",PKD,tomorrowDate)
	queryParamsDayAfter := fmt.Sprintf("?district_id=%d&date=%s",PKD,dayAfterTomorrow)
	finalURL := BASE_COWIN_URL + queryParamsToday
	FetchData(finalURL)
	finalURL = BASE_COWIN_URL + queryParamsTomorrow
	FetchData(finalURL)
	finalURL = BASE_COWIN_URL + queryParamsDayAfter
	FetchData(finalURL)
}

func FetchData(finalURL string) {
	res,err := http.Get(finalURL)
	if err != nil {
		log.Fatal(err)
	}
	
	body,err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	bodyStrings := string(body)

	var jsonResponse models.JSONResponse
	err = json.Unmarshal([]byte(bodyStrings), &jsonResponse)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(finalURL)
	if len(jsonResponse.Sessions) > 0 {
		controller.AddCenterToDB(jsonResponse.Sessions)
		controller.ExtractDetails(jsonResponse.Sessions)
	}

}