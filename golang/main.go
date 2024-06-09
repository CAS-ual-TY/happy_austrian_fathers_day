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

var filename = "../player_numbers.csv"
var attempts = 3

type ServerProperties struct {
	Numplayers string `json:"numplayers"`
}

type Server struct {
	Properties ServerProperties `json:"properties"`
}

type Response struct {
	Servers []Server `json:"servers"`
}

func fetch() (int, time.Time, []byte, error) {
	t := time.Now()
	resp, err := http.Get("https://servers.realitymod.com/api/ServerInfo")
	if err != nil {
		return 0, t, nil, err
	}
	defer resp.Body.Close()
	bytes, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		return 0, t, nil, err2
	}
	var response Response
	err3 := json.Unmarshal(bytes, &response)
	if err3 != nil {
		fmt.Println(string(bytes))
		return 0, t, bytes, err3
	}
	totalPlayers := 0
	for _, s := range response.Servers {
		numplayers, err4 := strconv.Atoi(s.Properties.Numplayers)
		if err4 != nil {
			return 0, t, bytes, err4
		}
		totalPlayers += numplayers
	}
	return totalPlayers, t, bytes, nil
}

func main() {
	var totalPlayers int
	var t time.Time
	var bytes []byte
	var err error = nil

	for i := range attempts {
		fmt.Printf("Attempt %d\n", (i + 1))
		totalPlayers, t, bytes, err = fetch()

		if err != nil {
			if i == attempts-1 {
				fmt.Printf("%d failed attempts... exiting without success\n", attempts)
				if bytes != nil {
					fmt.Println("Decoded bytes:")
					fmt.Println(string(bytes))
				}
				if e, ok := err.(*json.SyntaxError); ok {
					fmt.Printf("Syntax error at offset %d\n", e.Offset)
					if bytes != nil {
						fmt.Println("If CAS_ual_TY isnt completely retarded, this should be somewhere here:")
						off := int(e.Offset)
						minOff := max(0, off-off%8-8*3)
						maxOff := min(len(bytes), off-off%8+8*4)
						fmt.Printf("[%d, %d]\n", minOff, maxOff)
						fmt.Println(string(bytes[minOff:maxOff]))
					}
				}

				panic(err)
			}
		} else {
			fmt.Println("... success!")
			break
		}

		time.Sleep(time.Second * 31)
	}

	if err != nil {
		time.Sleep(time.Second * 10)

	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	newline := fmt.Sprintf("%s,%d\n", t.Round(time.Second).UTC().Format("2006-01-02 15:04:05"), totalPlayers)

	if _, err = f.WriteString(newline); err != nil {
		panic(err)
	}
}
