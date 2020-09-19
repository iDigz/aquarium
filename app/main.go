package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"os"
	"time"
)

var pinmap = [8]int{14, 15, 18, 23, 24, 25, 8, 7}
var workdays = [8][8]int{
	{13, 0, 15, 0, 15, 0, 19, 0},   // Основной свет
	{5, 30, 7, 0, 22, 0, 23, 30},   // Аэратор подстветка мощные диоды
	{5, 20, 13, 15, 19, 0, 21, 0},  // Помпа
	{9, 0, 13, 0, 19, 0, 20, 0},    // POWER GLO 40 Bт
	{19, 0, 20, 0, 20, 30, 22, 30}, // Сине-зелёная подсветка в крышке
	{6, 0, 9, 0, 21, 0, 23, 0},     // Подсветка в крышке ЛИНЗЫ
	{5, 25, 6, 25, 19, 0, 21, 0},   // Cиняя подстветка Тумбы
	{20, 0, 20, 30, 21, 30, 22, 0}, // Лампа Т5
}

var holidays = [8][8]int{
	{13, 0, 15, 0, 15, 0, 19, 0},     // Основной свет
	{21, 30, 22, 30, 22, 30, 23, 15}, // Аэратор подстветка мощные диоды
	{10, 30, 13, 15, 19, 0, 21, 0},   // Помпа
	{9, 0, 13, 0, 19, 0, 20, 0},      // POWER GLO 40 Bт
	{19, 0, 20, 0, 20, 30, 22, 30},   // Сине-зелёная подсветка в крышке
	{20, 0, 21, 30, 22, 30, 23, 15},  // Подсветка в крышке ЛИНЗЫ
	{19, 0, 20, 0, 20, 0, 21, 0},     // Cиняя подстветка Тумбы
	{20, 0, 20, 30, 21, 30, 22, 0},   // Лампа Т5
}

func main() {
	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Unmap gpio memory when done
	defer rpio.Close()

	for pin := 0; pin < len(pinmap); pin++ {
		rpio.Pin(pinmap[pin]).Output()
	}
Start:
	now := time.Now()
	if int(now.Weekday()) < 6 {
		fmt.Println("Starting plan for Workdays")
		for i := 0; i < len(workdays); i++ {
			time1 := time.Date(now.Year(), now.Month(), now.Day(), workdays[i][0], workdays[i][1], 0, 0, time.Local)
			time2 := time.Date(now.Year(), now.Month(), now.Day(), workdays[i][2], workdays[i][3], 0, 0, time.Local)
			time3 := time.Date(now.Year(), now.Month(), now.Day(), workdays[i][4], workdays[i][5], 0, 0, time.Local)
			time4 := time.Date(now.Year(), now.Month(), now.Day(), workdays[i][6], workdays[i][7], 0, 0, time.Local)
			check1 := now.After(time1) && now.Before(time2)
			check2 := now.After(time3) && now.Before(time4)
			if check1 || check2 {
				rpio.Pin(pinmap[i]).Low()
				fmt.Println("Pin ", i, "status ", check1 || check2)
			} else {
				rpio.Pin(pinmap[i]).High()
				fmt.Println("Pin ", i, "status ", check1 || check2)
			}
		}
	} else {
		fmt.Println("Starting plan for Holydays")
		for i := 0; i < len(holidays); i++ {
			time1 := time.Date(now.Year(), now.Month(), now.Day(), holidays[i][0], holidays[i][1], 0, 0, time.Local)
			time2 := time.Date(now.Year(), now.Month(), now.Day(), holidays[i][2], holidays[i][3], 0, 0, time.Local)
			time3 := time.Date(now.Year(), now.Month(), now.Day(), holidays[i][4], holidays[i][5], 0, 0, time.Local)
			time4 := time.Date(now.Year(), now.Month(), now.Day(), holidays[i][6], holidays[i][7], 0, 0, time.Local)
			check1 := now.After(time1) && now.Before(time2)
			check2 := now.After(time3) && now.Before(time4)
			if check1 || check2 {
				rpio.Pin(pinmap[i]).Low()
				fmt.Println("Pin ", i, "status ", check1 || check2)
			} else {
				rpio.Pin(pinmap[i]).High()
				fmt.Println("Pin ", i, "status ", check1 || check2)
			}
		}
	}
	fmt.Println("---------------------------------------")
	fmt.Println("Time now: ", now.Local())
	time.Sleep(2 * time.Second)
	goto Start
}
