# saga

package `saga` is a tiny library to help Golang compensating transaction.

In combination with code generator [thunk](https://github.com/Syuparn/thunk),
you can implement compensating transactions without nested `if err != nil` structures.

# Usage

Example:

```go
func (b *bookingApp) Book() (*Bookings, error) {
	// create new saga
	sg := saga.New()

	// run some process
	outboundTicket := saga.Make(sg, b.flightBookingService.Book("Tokyo", "Seoul", mustParseTime("2022/01/01 10:00")))
	// add compensation transaction, which will run if any following process failed
	sg.AddCompensation(b.flightBookingService.Cancel(outboundTicket))

	inboundTicket := saga.Make(sg, b.flightBookingService.Book("Seoul", "Tokyo", mustParseTime("2022/01/02 21:00")))
	sg.AddCompensation(b.flightBookingService.Cancel(inboundTicket))

	room := saga.Make(sg, b.hotelBookingService.Book(mustParseTime("2022/01/01 19:00")))

	// run all compensating transactions
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
```

```bash
# happy path
$ ./book
succeed to book flight!: flight ticket: from Tokyo to Seoul, 01/01 10:00
succeed to book flight!: flight ticket: from Seoul to Tokyo, 01/02 21:00
succeed to book hotel!: hotel room: number 699, 01/01
ready for traveling!!

# error triggers compensating transactions
$ ./book
succeed to book flight!: flight ticket: from Tokyo to Seoul, 01/01 10:00
failed to book flight...flight ticket: from Seoul to Tokyo, 01/02 21:00
succeed to cancel ticket: flight ticket: from Tokyo to Seoul, 01/01 10:00
1 error occurred:
        * failed to book flight...: flight ticket: from Seoul to Tokyo, 01/02 21:00

$ ./book
succeed to book flight!: flight ticket: from Tokyo to Seoul, 01/01 10:00
succeed to book flight!: flight ticket: from Seoul to Tokyo, 01/02 21:00
failed to book hotel...: hotel room: number 755, 01/01
succeed to cancel ticket: flight ticket: from Seoul to Tokyo, 01/02 21:00
succeed to cancel ticket: flight ticket: from Tokyo to Seoul, 01/01 10:00
1 error occurred:
        * failed to book hotel...: hotel room: number 755, 01/01
```

Try [example](https://github.com/Syuparn/saga/tree/main/_examples/book) and check how it works!
