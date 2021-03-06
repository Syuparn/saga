package saga

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/xerrors"
)

func TestNew(t *testing.T) {
	s := New()

	require.Equal(t, &Saga{}, s, "s must be zero value of *Saga.")
}

func TestSagaRun(t *testing.T) {
	tests := []struct {
		name           string
		saga           *Saga
		f              func() error
		expectedErrors []error
	}{
		{
			"function runs and succeeds",
			&Saga{},
			func() error { return nil },
			[]error{},
		},
		{
			"function runs and fails",
			&Saga{},
			func() error { return xerrors.Errorf("failed") },
			[]error{xerrors.Errorf("failed")},
		},
		{
			"if error is already raised, run does nothing",
			&Saga{
				errors: []error{xerrors.Errorf("previous error")},
			},
			func() error {
				t.Fatalf("this must not be called!")
				return xerrors.Errorf("failed")
			},
			[]error{xerrors.Errorf("previous error")},
		},
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			tt.saga.Run(tt.f)
			actual := tt.saga.errors

			require.Equal(t, len(tt.expectedErrors), len(actual), "number of errors are different")
			for i, e := range tt.expectedErrors {
				require.EqualError(t, actual[i], e.Error())
			}
		})
	}
}

func TestSagaAddCompensation(t *testing.T) {
	compensations := []func() error{
		func() error { return xerrors.Errorf("zero") },
		func() error { return xerrors.Errorf("one") },
		func() error { return xerrors.Errorf("two") },
	}

	tests := []struct {
		name         string
		saga         *Saga
		compensation Compensation
		expected     []Compensation
	}{
		{
			"saga without any compensations",
			&Saga{},
			compensations[0],
			[]Compensation{compensations[0]},
		},
		{
			"saga with a compensation",
			&Saga{
				compensations: []Compensation{compensations[0]},
			},
			compensations[1],
			[]Compensation{compensations[0], compensations[1]},
		},
		{
			"saga with multiple compensations",
			&Saga{
				compensations: []Compensation{compensations[0], compensations[1]},
			},
			compensations[2],
			[]Compensation{compensations[0], compensations[1], compensations[2]},
		},
		{
			"if any error is already raised, compensation is not added",
			&Saga{errors: []error{xerrors.Errorf("error")}},
			compensations[0],
			[]Compensation{},
		},
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			tt.saga.AddCompensation(tt.compensation)
			actual := tt.saga.compensations

			// NOTE: require.Equal cannot compare functions directly
			require.Equal(t, len(tt.expected), len(actual), "number of compensations are different")

			for i, e := range tt.expected {
				require.Equal(t, reflect.ValueOf(e).Pointer(), reflect.ValueOf(actual[i]).Pointer(),
					"compensations[%d] is diffrent", i)
			}
		})
	}
}

func TestSagaCompensate(t *testing.T) {
	tests := []struct {
		name           string
		saga           *Saga
		expectedErrors []error
	}{
		{
			"no compensating transactions are set",
			&Saga{},
			nil,
		},
		{
			"no errors have been raised",
			&Saga{
				compensations: []Compensation{
					func() error { return xerrors.Errorf("a") },
				},
			},
			nil,
		},
		{
			"error has been raised and compensation transaction succeeded",
			&Saga{
				compensations: []Compensation{
					func() error { return nil },
				},
				errors: []error{
					xerrors.Errorf("a"),
				},
			},
			[]error{xerrors.Errorf("a")},
		},
		{
			"error has been raised and compensation transaction raised another error",
			&Saga{
				compensations: []Compensation{
					func() error { return xerrors.Errorf("comp0") },
				},
				errors: []error{
					xerrors.Errorf("e"),
				},
			},
			[]error{
				xerrors.Errorf("e"),
				xerrors.Errorf("compensating transactions [0] failed: %w", xerrors.Errorf("comp0")),
			},
		},
		{
			"compensation transactions are executed in reversed order",
			&Saga{
				compensations: []Compensation{
					func() error { return xerrors.Errorf("comp0") },
					func() error { return xerrors.Errorf("comp1") },
				},
				errors: []error{
					xerrors.Errorf("e"),
				},
			},
			[]error{
				xerrors.Errorf("e"),
				xerrors.Errorf("compensating transactions [1] failed: %w", xerrors.Errorf("comp1")),
				xerrors.Errorf("compensating transactions [0] failed: %w", xerrors.Errorf("comp0")),
			},
		},
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			tt.saga.Compensate()
			actual := tt.saga.errors

			require.Equal(t, len(tt.expectedErrors), len(actual), "number of errors are different")

			for i, e := range tt.expectedErrors {
				require.EqualError(t, actual[i], e.Error())
			}
		})
	}
}

func TestSagaErrors(t *testing.T) {
	s := &Saga{
		errors: []error{
			xerrors.Errorf("a"),
			xerrors.Errorf("b"),
		},
	}

	require.Equal(t, s.errors, s.Errors())
}

func TestSagaError_WithError(t *testing.T) {
	s := &Saga{
		errors: []error{
			xerrors.Errorf("a"),
			xerrors.Errorf("b"),
		},
	}

	require.NotNil(t, s.Error())
}

func TestSagaError_WithoutError(t *testing.T) {
	s := &Saga{}

	require.Nil(t, s.Error())
}

func TestSagaHasError(t *testing.T) {
	tests := []struct {
		name     string
		saga     *Saga
		expected bool
	}{
		{
			"saga without errors",
			&Saga{},
			false,
		},
		{
			"saga with an error",
			&Saga{
				errors: []error{
					xerrors.Errorf("a"),
				},
			},
			true,
		},
		{
			"saga with errors",
			&Saga{
				errors: []error{
					xerrors.Errorf("a"),
					xerrors.Errorf("b"),
				},
			},
			true,
		},
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			actual := tt.saga.HasError()
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name           string
		saga           *Saga
		f              func() error
		expectedErrors []error
	}{
		{
			"function runs and succeeds",
			&Saga{},
			func() error { return nil },
			[]error{},
		},
		{
			"function runs and fails",
			&Saga{},
			func() error { return xerrors.Errorf("failed") },
			[]error{xerrors.Errorf("failed")},
		},
		{
			"if error is already raised, run does nothing",
			&Saga{
				errors: []error{xerrors.Errorf("previous error")},
			},
			func() error {
				t.Fatalf("this must not be called!")
				return xerrors.Errorf("failed")
			},
			[]error{xerrors.Errorf("previous error")},
		},
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			Run(tt.saga, tt.f)
			actual := tt.saga.errors

			require.Equal(t, len(tt.expectedErrors), len(actual), "number of errors are different")
			for i, e := range tt.expectedErrors {
				require.EqualError(t, actual[i], e.Error())
			}
		})
	}
}

func TestMake(t *testing.T) {
	tests := []struct {
		name           string
		saga           *Saga
		f              func() (string, error)
		expected       string
		expectedErrors []error
	}{
		{
			"function runs and succeeds",
			&Saga{},
			func() (string, error) { return "success", nil },
			"success",
			[]error{},
		},
		{
			"function runs and fails",
			&Saga{},
			func() (string, error) { return "", xerrors.Errorf("failed") },
			"",
			[]error{xerrors.Errorf("failed")},
		},
		{
			"if error is already raised, run does nothing",
			&Saga{
				errors: []error{xerrors.Errorf("previous error")},
			},
			func() (string, error) {
				t.Fatalf("this must not be called!")
				return "", xerrors.Errorf("failed")
			},
			"",
			[]error{xerrors.Errorf("previous error")},
		},
	}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			actual := Make(tt.saga, tt.f)
			actualErrors := tt.saga.errors

			require.Equal(t, tt.expected, actual)

			require.Equal(t, len(tt.expectedErrors), len(actualErrors), "number of errors are different")
			for i, e := range tt.expectedErrors {
				require.EqualError(t, actualErrors[i], e.Error())
			}
		})
	}
}
