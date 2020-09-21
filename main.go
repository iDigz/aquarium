package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Workdays map[int][]map[string]map[string]int
type Holidays map[int][]map[string]map[string]int
type Pins []int

func Find(slice []bool, val bool) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func main() {
	filename_wd, _ := filepath.Abs("./config/workdays.yaml")
	yamlFile_wd, err := ioutil.ReadFile(filename_wd)
	if err != nil {
		panic(err)
	}

	filename_hd, _ := filepath.Abs("./config/holidays.yaml")
	yamlFile_hd, err := ioutil.ReadFile(filename_hd)
	if err != nil {
		panic(err)
	}

	filename_pin, _ := filepath.Abs("./config/pins.yaml")
	yamlFile_pin, err := ioutil.ReadFile(filename_pin)
	if err != nil {
		panic(err)
	}

	var config_wd Workdays
	var config_hd Holidays
	var config_pin Pins

	err = yaml.Unmarshal(yamlFile_wd, &config_wd)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile_hd, &config_hd)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile_pin, &config_pin)
	if err != nil {
		panic(err)
	}

	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmap gpio memory when done
	//defer rpio.Close()

	for pin := 0; pin < len(config_pin); pin++ {
		rpio.Pin(config_pin[pin]).Output()
	}
Start:
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
	now := time.Now()
	if int(now.Weekday()) < 6 {
		fmt.Println("Config for Weekdays")
		for i := 0; i < len(config_wd); i++ {
			check := make([]bool, 6)
			for j := 0; j < len(config_wd[i]); j++ {
				var starttime [7]time.Time
				var stoptime [7]time.Time
				starttime[j] = time.Date(now.Year(), now.Month(), now.Day(), config_wd[i][j]["start"]["H"], config_wd[i][j]["start"]["M"], 0, 0, time.Local)
				stoptime[j] = time.Date(now.Year(), now.Month(), now.Day(), config_wd[i][j]["stop"]["H"], config_wd[i][j]["stop"]["M"], 0, 0, time.Local)
				check[j] = now.After(starttime[j]) && now.Before(stoptime[j])
				fmt.Println("Cicle", j, "--->", starttime[j].Hour(), starttime[j].Minute(), "-", stoptime[j].Hour(), stoptime[j].Minute(), "--->", check[j])
			}
			if Find(check, true) {
				fmt.Println("Pin set in LOW")
				rpio.Pin(config_pin[i]).Low()
			} else {
				fmt.Println("Pin set in HIGH")
				rpio.Pin(config_pin[i]).High()
			}
		}
	} else {
		fmt.Println("Config for Holidays")
		for i := 0; i < len(config_hd); i++ {
			check := make([]bool, 6)
			for j := 0; j < len(config_hd[i]); j++ {
				var starttime [7]time.Time
				var stoptime [7]time.Time
				starttime[j] = time.Date(now.Year(), now.Month(), now.Day(), config_hd[i][j]["start"]["H"], config_hd[i][j]["start"]["M"], 0, 0, time.Local)
				stoptime[j] = time.Date(now.Year(), now.Month(), now.Day(), config_hd[i][j]["stop"]["H"], config_hd[i][j]["stop"]["M"], 0, 0, time.Local)
				check[j] = now.After(starttime[j]) && now.Before(stoptime[j])
				fmt.Println("Cicle", j, "--->", starttime[j].Hour(), starttime[j].Minute(), "-", stoptime[j].Hour(), stoptime[j].Minute(), "--->", check[j])
			}
			if Find(check, true) {
				fmt.Println("Pin set in LOW")
				rpio.Pin(config_pin[i]).Low()
			} else {
				fmt.Println("Pin set in HIGH")
				rpio.Pin(config_pin[i]).High()
			}
		}
	}

	fmt.Println("---------------------------------------")
	fmt.Println("Time now: ", now.Local())
	time.Sleep(1 * time.Second)
	goto Start
}
