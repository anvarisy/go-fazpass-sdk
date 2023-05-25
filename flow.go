package fazpass

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Default() FlowInterface {
	return &Flow{}
}

type Flow struct {
}

type FlowInterface interface {
	wrapingData(f *Fazpass, model interface{}) ([]byte, error)
	sendingData(f *Fazpass, wrappedMessage []byte, urls string) (*http.Response, error)
	extractingData(f *Fazpass, response *http.Response, data *Data) (*Data, error)
}

func (flow *Flow) wrapingData(f *Fazpass, model interface{}) ([]byte, error) {
	var (
		err        error
		marshalled []byte
		encrypted  []byte
		wrapedData []byte
	)
	marshalled, err = json.Marshal(model)
	encrypted, err = encryptWithPublicKey(marshalled, f.PublicKey)
	message := base64.StdEncoding.EncodeToString([]byte(encrypted))
	wrapedData, err = json.Marshal(&Transmission{Message: message})
	return wrapedData, err
}

func (flow *Flow) sendingData(f *Fazpass, wrappedMessage []byte, urls string) (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest("POST", f.BaseUrl+urls, bytes.NewReader(wrappedMessage))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+f.MerchantKey)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (flow *Flow) extractingData(f *Fazpass, response *http.Response, data *Data) (*Data, error) {
	defer response.Body.Close()
	messageBodyResponse, _ := io.ReadAll(response.Body)
	transmissionResponse := &Transmission{}
	json.Unmarshal(messageBodyResponse, transmissionResponse)
	decryptMessage, err := base64.StdEncoding.DecodeString(string(transmissionResponse.Message))

	decrypted, _ := decryptWithPrivateKey(decryptMessage, f.PrivateKey)
	json.Unmarshal(decrypted, data)
	fmt.Print(data.SessionId)
	if err != nil {
		return data, err
	}

	return data, nil
}
