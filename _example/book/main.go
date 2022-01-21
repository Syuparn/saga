package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	app := NewBookingApp(&flightBookingService{}, &hotelBookingService{})

	_, err := app.Book()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	fmt.Println("ready for traveling!!")
}
