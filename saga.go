package saga

import (
	"golang.org/x/xerrors"

	"github.com/hashicorp/go-multierror"
)

// Compensation is a type alias of compensating transaction function.
type Compensation = func() error

// Saga controls compensating transactions and keeps errors raised in the middle of exection.
type Saga struct {
	errors        []error
	compensations []Compensation
}

// Run runs the recieved function. Raised error is kept inside.
// If any error has already raised in saga, it does nothing.
func (s *Saga) Run(f func() error) {
	if s.HasError() {
		return
	}

	if err := f(); err != nil {
		s.errors = append(s.errors, err)
	}
}

// AddCompensation adds a compensating transaction to the saga.
func (s *Saga) AddCompensation(c Compensation) {
	s.compensations = append(s.compensations, c)
}

// Compensate executes compensating transactions.
// If no errors have been raised so far, it does nothing.
// The compensation transactions are executed in the reversed order of addition.
func (s *Saga) Compensate() {
	// if no errors have been raised, do nothing
	if !s.HasError() {
		return
	}

	for i := len(s.compensations) - 1; i >= 0; i-- {
		c := s.compensations[i]

		if err := c(); err != nil {
			s.errors = append(s.errors, xerrors.Errorf("compensating transactions [%d] failed: %w", i, err))
		}
	}
}

// Errors returns all errors raised during the saga, including compensating transaction errors.
func (s *Saga) Errors() []error {
	return s.errors
}

// Error returns an error raised during the saga. compensating transaction errors are wrapped inside.
func (s *Saga) Error() error {
	if !s.HasError() {
		return nil
	}

	return multierror.Append(s.errors[0], s.errors[1:]...)
}

// HasError returns whether error(s) is raised in the saga.
func (s *Saga) HasError() bool {
	return len(s.errors) > 0
}

// New creates a new saga.
func New() *Saga {
	return &Saga{}
}

// Run runs the recieved function in saga. Raised error is kept inside saga.
// If any error has already raised in saga, it does nothing.
func Run(s *Saga, f func() error) {
	s.Run(f)
}
