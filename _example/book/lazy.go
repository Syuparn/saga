// TODO: replace lazy structs with generated codes by thunk once https://github.com/golang/go/issues/46041 is fixed

package main

import "time"

type LazyHotelBookingService struct {
	inner HotelBookingService
}

func NewLazyHotelBookingService(inner HotelBookingService) *LazyHotelBookingService {
	return &LazyHotelBookingService{inner: inner}
}

func (s *LazyHotelBookingService) Book(date time.Time) func() (*HotelRoom, error) {
	return func() (*HotelRoom, error) {
		return s.inner.Book(date)
	}
}

func (s *LazyHotelBookingService) Cancel(r *HotelRoom) func() error {
	return func() error {
		return s.inner.Cancel(r)
	}
}

type LazyFlightBookingService struct {
	inner FlightBookingService
}

func NewLazyFlightBookingService(inner FlightBookingService) *LazyFlightBookingService {
	return &LazyFlightBookingService{inner: inner}
}

func (s *LazyFlightBookingService) Book(from, to string, date time.Time) func() (*FlightTicket, error) {
	return func() (*FlightTicket, error) {
		return s.inner.Book(from, to, date)
	}
}

func (s *LazyFlightBookingService) Cancel(t *FlightTicket) func() error {
	return func() error {
		return s.inner.Cancel(t)
	}
}
