package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

var filename = "./player_numbers.csv"

type ServerProperties struct {
	Numplayers string `json:"numplayers"`
}

type Server struct {
	Properties ServerProperties `json:"properties"`
}

type Response struct {
	Servers []Server `json:"servers"`
}

func fetch() (int, time.Time, error) {
	t := time.Now()
	resp, err := http.Get("https://servers.realitymod.com/api/ServerInfo")
	if err != nil {
		return 0, t, err
	}
	defer resp.Body.Close()
	bytes, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		return 0, t, err2
	}
	var response Response
	err3 := json.Unmarshal(bytes, &response)
	if err3 != nil {
		return 0, t, err3
	}
	totalPlayers := 0
	for _, s := range response.Servers {
		numplayers, err4 := strconv.Atoi(s.Properties.Numplayers)
		if err4 != nil {
			return 0, t, err4
		}
		totalPlayers += numplayers
	}
	return totalPlayers, t, nil
}

func main() {
	var totalPlayers int
	var t time.Time
	var err error = nil

	for i := range 3 {
		fmt.Printf("Attempt: %d\n", (i + 1))
		totalPlayers, t, err = fetch()

		if err != nil {
			if i == 3-1 {
				fmt.Println("3 failed attempts... exiting without success")
				panic(err)
			} else {
				time.Sleep(time.Second * 10)
			}
		}
	}

	if err != nil {
		time.Sleep(time.Second * 10)

	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	newline := fmt.Sprintf("%s,%d\n", t.Round(time.Second).UTC(), totalPlayers)

	if _, err = f.WriteString(newline); err != nil {
		panic(err)
	}
}
