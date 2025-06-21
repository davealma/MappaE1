package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type ProbeResponse struct {
	Distance string	`json:"distance"`
	Time     string `json:"time"`
	DistanceQ	float64
	TimeQ		float64
}

func PostSolution(speed int) {
	body := []byte(`{
		"speed": `+ strconv.Itoa(speed) +`,
		
	}`)
	request, err := http.NewRequest("POST", "https://makers-challenge.altscore.ai/v1/s1/e1/solution", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("API-KEY", os.Getenv("API_KEY"))

	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	bodyResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("Response Status: ", string(bodyResp))
}

func GetProbeResponse(tick time.Time) ProbeResponse{
	baseUrl := os.Getenv("API_URL")
 	req, err := http.NewRequest("Get", baseUrl + "/v1/s1/e1/resources/measurement", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("API-KEY", os.Getenv("API_KEY"))
	
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("Response Status: ", string(body))
	//{"distance":"695 AU","time":"1.7160493827160495 hours"}	
	
	var probeResp ProbeResponse
	err = json.Unmarshal(body, &probeResp)

	if _, errDistance := fmt.Sscanf(probeResp.Distance, "%f AU", &probeResp.DistanceQ); errDistance != nil{
		probeResp.DistanceQ = 0
	}
	if _, errTime := fmt.Sscanf(probeResp.Time, "%f hours", &probeResp.TimeQ); errTime != nil {
		probeResp.TimeQ = 0
	}
	if err !=nil {
		panic(err)
	}	
	time.Sleep(1*time.Second)
	return probeResp
}


func main() {
	godotenv.Load()
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for t := range ticker.C {
			probeResp := GetProbeResponse(t)			
			fmt.Println("distance: ", probeResp.DistanceQ)			
			fmt.Println("time: ", probeResp.TimeQ)
			if probeResp.DistanceQ > 0 && probeResp.TimeQ > 0 {				
				velocity := probeResp.DistanceQ / probeResp.TimeQ
				fmt.Println("Speed Round", math.Round(velocity))
				ticker.Stop()
			}
		}
	}()
	for {
		time.Sleep(100 * time.Millisecond)
	}	
}