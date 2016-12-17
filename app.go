package main

import (
	"io/ioutil"
	"strings"
	"regexp"
	"strconv"
	"os/exec"
	"time"
)

var co2Regex = regexp.MustCompile(`.*\[VALUE\].*CO2\:\s*(\d+)`)
var co2Value = 0

func main()  {

	go func () {
		threshold := 1800 // see http://www.raumluft.org/natuerliche-mechanische-lueftung/co2-als-lueftungsindikator/
		for range time.Tick(1 * time.Second){
			if co2Value > threshold {
				piipCmd := "/home/canbus/wiringPi/piip.sh"
				exec.Command(piipCmd).Output()
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
			co2Value, _ = strconv.Atoi(co2ValueString)
			println(co2Value)
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