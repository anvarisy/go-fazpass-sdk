package fazpass

import (
	"crypto/rsa"
	"net/http"

	"github.com/stretchr/testify/mock"
)

type FlowMock struct {
	mock.Mock
}

// extractingData implements FlowInterface
func (fm *FlowMock) ExtractingData(privKey *rsa.PrivateKey, response *http.Response, data *Data) (*Data, error) {
	ret := fm.Called(privKey, response, data)
	return ret.Get(0).(*Data), ret.Error(1)
}

// sendingData implements FlowInterface
func (fm *FlowMock) SendingData(baseUrl string, wrappedMessage []byte, merchantKey string) (*http.Response, error) {
	ret := fm.Called(baseUrl, wrappedMessage, merchantKey)
	return ret.Get(0).(*http.Response), ret.Error(1)
}

// wrapingData implements FlowInterface
func (fm *FlowMock) WrappingData(pubKey *rsa.PublicKey, model interface{}) ([]byte, error) {
	ret := fm.Called(pubKey, model)
	return ret.Get(0).([]byte), ret.Error(1)
}

func Init(f *Flow) FlowInterface {
	return &FlowMock{}
}
