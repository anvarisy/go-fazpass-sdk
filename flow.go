package fazpass

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/anvarisy/go-fazpass-sdk/utils"
	"io"
	"net/http"
)

func Default() FlowInterface {
	return &Flow{}
}

type Flow struct {
	client *http.Client
}

type FlowInterface interface {
	WrappingData(pubKey *rsa.PublicKey, model interface{}) ([]byte, error)
	SendingData(baseUrl string, wrappedMessage []byte, merchantKey string) (*http.Response, error)
	ExtractingData(privKey *rsa.PrivateKey, response *http.Response, data *Data) (*Data, error)
}

func (flow *Flow) WrappingData(pubKey *rsa.PublicKey, model interface{}) ([]byte, error) {
	var (
		err        error
		marshalled []byte
		encrypted  []byte
		wrapedData []byte
	)
	marshalled, err = json.Marshal(model)
	encrypted, err = utils.EncryptWithPublicKey(marshalled, pubKey)
	message := base64.StdEncoding.EncodeToString([]byte(encrypted))
	wrapedData, err = json.Marshal(&Transmission{Message: message})
	return wrapedData, err
}

func (flow *Flow) SendingData(baseUrl string, wrappedMessage []byte, merchantKey string) (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest("POST", baseUrl, bytes.NewReader(wrappedMessage))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+merchantKey)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (flow *Flow) ExtractingData(privKey *rsa.PrivateKey, response *http.Response, data *Data) (*Data, error) {
	defer response.Body.Close()
	messageBodyResponse, _ := io.ReadAll(response.Body)
	transmissionResponse := &Transmission{}
	json.Unmarshal(messageBodyResponse, transmissionResponse)
	decryptMessage, err := base64.StdEncoding.DecodeString(string(transmissionResponse.Message))

	decrypted, _ := utils.DecryptWithPrivateKey(decryptMessage, privKey)
	json.Unmarshal(decrypted, data)
	fmt.Print(data.SessionId)
	if err != nil {
		return data, err
	}

	return data, nil
}
