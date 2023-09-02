package transactions

import (
	"github.com/stretchr/testify/mock"
)

type MockUtils struct {
	mock.Mock
}

func (m *MockUtils) ValidateFloatPrecision(amount float64) error {
	args := m.Called(amount)
	return args.Error(0)
}
