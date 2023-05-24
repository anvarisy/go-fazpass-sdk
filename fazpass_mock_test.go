package fazpass

import (
	"github.com/stretchr/testify/mock"
)

type FazpassMock struct {
	mock.Mock
}

// Check implements FazpassInterface
func (fm *FazpassMock) Check(email string, phone string, encData string) (*Data, error) {
	ret := fm.Called(email, phone, encData)
	return ret.Get(0).(*Data), ret.Error(1)
}

// EnrollDevice implements FazpassInterface
func (fm *FazpassMock) EnrollDevice(email string, phone string, encData string) (*Data, error) {
	ret := fm.Called(email, phone, encData)
	return ret.Get(0).(*Data), ret.Error(1)
}

// RemoveDevice implements FazpassInterface
func (fm *FazpassMock) RemoveDevice(fazpassId string, encData string) (*Data, error) {
	ret := fm.Called(fazpassId, encData)
	return ret.Get(0).(*Data), ret.Error(1)
}

// ValidateDevice implements FazpassInterface
func (fm *FazpassMock) ValidateDevice(fazpassId string, encData string) (*Data, error) {
	ret := fm.Called(fazpassId, encData)
	return ret.Get(0).(*Data), ret.Error(1)
}

func Init(f *Fazpass) FazpassInterface {
	return &FazpassMock{}
}
