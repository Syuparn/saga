package main

import (
	"time"

	"github.com/syuparn/saga"
)

type bookingApp struct {
	flightBookingService *LazyFlightBookingService
	hotelBookingService  *LazyHotelBookingService
}

func NewBookingApp(fs FlightBookingService, hs HotelBookingService) *bookingApp {
	return &bookingApp{
		flightBookingService: NewLazyFlightBookingService(fs),
		hotelBookingService:  NewLazyHotelBookingService(hs),
	}
}

type Bookings struct {
	outboundTicket *FlightTicket
	inboundTicket  *FlightTicket
	room           *HotelRoom
}

func (b *bookingApp) Book() (*Bookings, error) {
	sg := saga.New()

	outboundTicket := saga.Make(sg, b.flightBookingService.Book("Tokyo", "Seoul", mustParseTime("2022/01/01 10:00")))
	sg.AddCompensation(b.flightBookingService.Cancel(outboundTicket))

	inboundTicket := saga.Make(sg, b.flightBookingService.Book("Seoul", "Tokyo", mustParseTime("2022/01/02 21:00")))
	sg.AddCompensation(b.flightBookingService.Cancel(inboundTicket))

	room := saga.Make(sg, b.hotelBookingService.Book(mustParseTime("2022/01/01 19:00")))

	sg.Compensate()

	if sg.HasError() {
		return nil, sg.Error()
	}

	return &Bookings{
		outboundTicket: outboundTicket,
		inboundTicket:  inboundTicket,
		room:           room,
	}, nil
}

func mustParseTime(s string) time.Time {
	t, err := time.Parse("2006/01/02 15:04", s)
	if err != nil {
		panic(err)
	}

	return t
}
