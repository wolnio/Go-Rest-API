package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"log"
	"context"
	"time"
)

type randomData struct {
	Stddev float64
	Data   []int
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
}

/* Calculate standard deviation
		 */
func calcStdDev(dataArray []int) float64{

		 var deviation float64
		sumVal := 0
		for _, item := range dataArray {
			sumVal += item
		}

		var meanVal float64 = float64(sumVal) / float64(len(dataArray))

		var xVal float64 = 0
		for _, item := range dataArray {
			xVal += math.Pow(float64(item)-meanVal, 2)
		}

		deviation = math.Sqrt(xVal / float64(len(dataArray)))
		return deviation
}

func getRandomInt(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	mapQuery, err := url.ParseQuery(r.URL.RawQuery)

	queryReq, err := strconv.Atoi(mapQuery["requests"][0])
	if err != nil {
		fmt.Fprintf(w, "Incorrect type of query parameters (requests)!\n")
	}

	//this var is never used but I map it and convert it to int anyway just to check if it's correct type
	queryLen, err := strconv.Atoi(mapQuery["length"][0])
	if err != nil {
		fmt.Fprintf(w, "Incorrect type of query parameters (length)!\n")
	}
	_ = queryLen

	randomDataSet := []randomData{}

	for i := 0; i < queryReq; i++ {
		
		var dataTemp []int

		endpoint := "https://www.random.org/integers/?min=1&max=10&col=1&base=10&format=plain&rnd=new"

		timeoutContext, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
		defer cancel()

		requ, err := http.NewRequestWithContext(timeoutContext, "GET", endpoint, nil)
		if err != nil {
			fmt.Println("No response from request\n")
		}

		/* Adds query parameter to the request URL - number of requested random integers
		*/
		q := requ.URL.Query()
		q.Add("num",mapQuery["length"][0])
		requ.URL.RawQuery = q.Encode()

		resp, err := http.DefaultClient.Do(requ)
		if err != nil {
			panic(err)
		}

		/*	Read given API response
		*/
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body) // response body is []byte
		if err != nil {
			fmt.Println(err)
		}

		/* Convert given body which is []byte type to []int type
		 */
		lines := strings.Split(string(body), "\n")
		for i := range lines {
			line := strings.TrimSpace(lines[i])
			if line == "" {
				continue
			}

			numb, err := strconv.Atoi(line)
			if err != nil {
				fmt.Println(err)
				fmt.Println("here")
				return
			}

			dataTemp = append(dataTemp, numb)
		}

		randomDataSet = append(randomDataSet, randomData{Stddev: calcStdDev(dataTemp), Data: dataTemp})
	}

	var dataAll []int
	for _, item := range randomDataSet{
		dataAll=append(dataAll, item.Data...)
	}

	randomDataSet = append(randomDataSet, randomData{Stddev: calcStdDev(dataAll), Data: dataAll})

	/* Represent data as JSON
	 */
	js, err := json.MarshalIndent(randomDataSet, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
	fmt.Println("Endpoint Hit: getRandomInt")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/random/mean", getRandomInt)

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

func main() {
	handleRequests()
}