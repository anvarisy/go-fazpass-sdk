package fazpass

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type FlowMock struct {
	mock.Mock
}

// extractingData implements FlowInterface
func (fm *FlowMock) extractingData(f *Fazpass, response *http.Response, data *Data) (*Data, error) {
	ret := fm.Called(f, response, data)
	return ret.Get(0).(*Data), ret.Error(1)
}

// sendingData implements FlowInterface
func (fm *FlowMock) sendingData(f *Fazpass, wrappedMessage []byte, urls string) (*http.Response, error) {
	ret := fm.Called(f, wrappedMessage, urls)
	return ret.Get(0).(*http.Response), ret.Error(1)
}

// wrapingData implements FlowInterface
func (fm *FlowMock) wrapingData(f *Fazpass, model interface{}) ([]byte, error) {
	ret := fm.Called(f, model)
	return ret.Get(0).([]byte), ret.Error(1)
}

/*
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
*/
func Init(f *Flow) FlowInterface {
	return &FlowMock{}
}
