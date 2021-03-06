// Code generated by thunk; DO NOT EDIT.
package main

import (
	"time"
)

type LazyFlightBookingService struct {
	inner FlightBookingService
}

func (l *LazyFlightBookingService) Book(from string, to string, date time.Time) func() (*FlightTicket, error) {
	return func() (*FlightTicket, error) {
		return l.inner.Book(from, to, date)
	}
}

func (l *LazyFlightBookingService) Cancel(ticket *FlightTicket) func() error {
	return func() error {
		return l.inner.Cancel(ticket)
	}
}

func NewLazyFlightBookingService(inner FlightBookingService) *LazyFlightBookingService {
	return &LazyFlightBookingService{inner: inner}
}

type LazyHotelBookingService struct {
	inner HotelBookingService
}

func (l *LazyHotelBookingService) Book(date time.Time) func() (*HotelRoom, error) {
	return func() (*HotelRoom, error) {
		return l.inner.Book(date)
	}
}

func (l *LazyHotelBookingService) Cancel(arg0 *HotelRoom) func() error {
	return func() error {
		return l.inner.Cancel(arg0)
	}
}

func NewLazyHotelBookingService(inner HotelBookingService) *LazyHotelBookingService {
	return &LazyHotelBookingService{inner: inner}
}
