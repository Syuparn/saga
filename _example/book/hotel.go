package main

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/xerrors"
)

type HotelRoom struct {
	number int
	date   time.Time
}

func NewHotelRoom(date time.Time) *HotelRoom {
	return &HotelRoom{
		number: rand.Intn(9)*100 + rand.Intn(99) + 101, // 101~999
		date:   date,
	}
}

func (r *HotelRoom) String() string {
	return fmt.Sprintf("hotel room: number %d, %s", r.number, r.date.Format("01/02"))
}

type HotelBookingService interface {
	Book(date time.Time) (*HotelRoom, error)
	Cancel(*HotelRoom) error
}

type hotelBookingService struct{}

func (s *hotelBookingService) Book(date time.Time) (*HotelRoom, error) {
	room := NewHotelRoom(date)

	if rand.Intn(2) == 0 {
		fmt.Printf("failed to book hotel...: %s\n", room)
		return nil, xerrors.Errorf("failed to book hotel...: %s", room)
	}

	fmt.Printf("succeed to book hotel!: %s\n", room)
	return room, nil
}

func (s *hotelBookingService) Cancel(room *HotelRoom) error {
	if rand.Intn(4) == 0 {
		fmt.Printf("failed to cancel hotel...: %s\n", room)
		return xerrors.Errorf("failed to cancel hotel...: %s", room)
	}

	fmt.Printf("succeed to cancel hotel\n: %s", room)
	return nil
}
