package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"text/template"
	"time"
)

type StatusData struct {
	Status struct {
		Water int `json:"water"`
		Wind  int `json:"wind"`
	} `json:"status"`
}

var tStatusWind string
var tStatusWater string

func main() {
	go AutoReloadJSON()
	fmt.Println("jalankan autoreloadjson")

	http.HandleFunc("/", AutoloadtoWeb)
	http.ListenAndServe(":8080", nil)
	fmt.Println("working")

}

func AutoReloadJSON() {
	for {
		fmt.Println("generate")
		min := 1
		max := 21
		wind := rand.Intn(max-min) + min
		water := rand.Intn(max-min) + min

		data := StatusData{}
		data.Status.Wind = wind
		data.Status.Water = water

		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Fatal("Error occurate while marshaling data status")
		}

		err = ioutil.WriteFile("data.json", jsonData, 0644)

		if err != nil {
			log.Fatal("Error while writing data to file")
		}
		time.Sleep(15 * time.Second)

	}
}

func AutoloadtoWeb(w http.ResponseWriter, r *http.Request) {

	jsonFile, err := os.Open("data.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Succes open file")
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var data StatusData

	json.Unmarshal(byteValue, &data)

	fmt.Println(data.Status.Water)
	fmt.Println(data.Status.Wind)

	var datastatus = map[string]string{}
	switch {
	case data.Status.Water <= 5:
		tStatusWater = "Aman"
	case data.Status.Water > 5 && data.Status.Water <= 8:
		tStatusWater = "Siaga"
	case data.Status.Water > 8:
		tStatusWater = "Bahaya"
	}

	switch {
	case data.Status.Water <= 5:
		tStatusWater = "Aman"
	case data.Status.Water > 5 && data.Status.Water <= 8:
		tStatusWater = "Siaga"
	case data.Status.Water > 8:
		tStatusWater = "Bahaya"
	}

	switch {
	case data.Status.Wind <= 6:
		tStatusWind = "Aman"
	case data.Status.Wind > 6 && data.Status.Wind <= 15:
		tStatusWind = "Siaga"
	case data.Status.Wind > 15:
		tStatusWind = "Bahaya"
	}

	datastatus["statuswater"] = tStatusWater
	datastatus["statuswind"] = tStatusWind
	fmt.Println(tStatusWater)
	fmt.Println(tStatusWind)
	tpl, err := template.ParseFiles("main.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tpl.Execute(w, datastatus)

	defer jsonFile.Close()
}
