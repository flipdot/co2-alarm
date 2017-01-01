package main

import (
	"io/ioutil"
	"strings"
	"regexp"
	"strconv"
	"os/exec"
	"fmt"
	"net/http"
	"bytes"
	"time"
)

var co2Regex = regexp.MustCompile(`.*\[VALUE\].*CO2\:\s*(\d+)`)
var co2Value = 0.0

func main()  {

	/*
	go func () {
		bps := calcBps(co2Value)
		for range time.Tick(time.Duration(bps) * time.Second){
			if co2Value > normalCo2 {
				piipCmd := "/home/canbus/wiringPi/piip.sh"
				exec.Command(piipCmd).Output()
			}
		}
	}()
	*/

	go func () {
		client := &http.Client{}
		for range time.Tick(1 * time.Minute){
			url := fmt.Sprintf("https://api.flipdot.org/sensors/co2/chill/%.2f/ppm", co2Value)
			req, err := http.NewRequest("POST", url, nil)
			if err != nil {
				fmt.Println(err)
				return
			}

			res, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
				return
			}

			if res.StatusCode != 200 {
				buf := new(bytes.Buffer)
				buf.ReadFrom(res.Body)
				fmt.Println(buf)
			}
		}
	}()

	for {
		co2Output := readOutput()
		lines := strings.Split(co2Output, "\n")
		for _, lin  := range lines {
			groups := co2Regex.FindStringSubmatch(lin)
			noMatch := len(groups) < 1
			if noMatch {
				continue
			}

			co2ValueString := groups[1]
			co2ValueInt, _ := strconv.Atoi(co2ValueString)
			co2Value = float64(co2ValueInt)
		}
	}
}


func readOutput() string {
	useFile := false
	if useFile {
		dat, _ := ioutil.ReadFile("sample/out")
		return string(dat)
	} else {
		out, _ := exec.Command("/home/pi/co2-monitor/monitor").Output()
		return string(out)
	}
}