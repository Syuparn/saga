package main

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/xerrors"
)

type FlightTicket struct {
	from string
	to   string
	date time.Time
}

func NewFlightTicket(from, to string, date time.Time) *FlightTicket {
	return &FlightTicket{
		from: from,
		to:   to,
		date: date,
	}
}

func (t *FlightTicket) String() string {
	return fmt.Sprintf("flight ticket: from %s to %s, %s", t.from, t.to, t.date.Format("01/02 15:04"))
}

type FlightBookingService interface {
	Book(from, to string, date time.Time) (*FlightTicket, error)
	Cancel(ticket *FlightTicket) error
}

type flightBookingService struct{}

func (s *flightBookingService) Book(from, to string, date time.Time) (*FlightTicket, error) {
	ticket := NewFlightTicket(from, to, date)

	if rand.Intn(2) == 0 {
		fmt.Printf("failed to book flight...%s\n", ticket)
		return nil, xerrors.Errorf("failed to book flight...: %s", ticket)
	}

	fmt.Printf("succeed to book flight!: %s\n", ticket)
	return ticket, nil
}

func (s *flightBookingService) Cancel(ticket *FlightTicket) error {
	if rand.Intn(4) == 0 {
		fmt.Printf("failed to cancel flight...: %s\n", ticket)
		return xerrors.Errorf("failed to cancel flight...: %s", ticket)
	}

	fmt.Printf("succeed to cancel ticket: %s\n", ticket)
	return nil
}
