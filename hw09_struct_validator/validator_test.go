package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

type UserRole string

type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	UnsupportedFields struct {
		BoolField  bool           `validate:"in:true"`
		FloatField float64        `validate:"min:0.5"`
		MapField   map[string]int `validate:"len:1"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "John Doe",
				Age:    30,
				Email:  "john@example.com",
				Role:   "admin",
				Phones: []string{"12345678901", "09876543210"},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "12345678901234567890123456789012345",
				Name:   "John Doe",
				Age:    17,
				Email:  "john@example.com",
				Role:   "user",
				Phones: []string{"12345678901", "09876543210"},
			},
			expectedErr: ValidationErrors{
				ValidationError{"ID", ErrValidationLen{Field: "ID", Len: 36}},
				ValidationError{"Age", ErrValidationMin{Field: "Age", Min: 18}},
				ValidationError{"Role", ErrValidationIn{Field: "Role", Set: []string{"admin", "stuff"}}},
			},
		},
		{
			in: App{
				Version: "1.0.0",
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "1.0",
			},
			expectedErr: ValidationErrors{
				ValidationError{"Version", ErrValidationLen{Field: "Version", Len: 5}},
			},
		},
		{
			in: Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 200,
				Body: "",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 201,
				Body: "",
			},
			expectedErr: ValidationErrors{
				ValidationError{"Code", ErrValidationIn{Field: "Code", Set: []string{"200", "404", "500"}}},
			},
		},
		{
			in: UnsupportedFields{
				BoolField:  true,
				FloatField: 1.0,
				MapField:   map[string]int{"one": 1},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "BoolField",
					Err:   UnsupportedFieldTypeError{Field: "BoolField", Type: "bool"},
				},
				ValidationError{
					Field: "FloatField",
					Err:   UnsupportedFieldTypeError{Field: "FloatField", Type: "float64"},
				},
				ValidationError{
					Field: "MapField",
					Err:   UnsupportedFieldTypeError{Field: "MapField", Type: "map"},
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			switch {
			case err != nil && tt.expectedErr == nil:
				t.Errorf("expected no error, but got error: %v", err)
			case err != nil && tt.expectedErr != nil:
				var ve, expectedVE ValidationErrors
				if errors.As(err, &ve) && errors.As(tt.expectedErr, &expectedVE) {
					if !validationErrorsEqual(ve, expectedVE) {
						t.Errorf("expected error %v, got %v", tt.expectedErr, err)
					}
				} else {
					t.Errorf("expected error %v, got %v", tt.expectedErr, err)
				}
			case err == nil && tt.expectedErr != nil:
				t.Errorf("expected error %v, but got no error", tt.expectedErr)
			}
		})
	}
}

func validationErrorsEqual(a, b ValidationErrors) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i].Field != b[i].Field || a[i].Err.Error() != b[i].Err.Error() {
			return false
		}
	}

	return true
}
