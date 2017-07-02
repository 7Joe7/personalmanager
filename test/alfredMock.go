package test

import (
	"github.com/stretchr/testify/mock"
)

func NewAlfredMock() *alfredMock {
	return &alfredMock{}
}

type alfredMock struct {
	mock.Mock
}

func (m *alfredMock) PrintEntities(entities interface{}) {
	m.Called(entities)
}

func (m *alfredMock) PrintResult(result string) {
	m.Called(result)
}
