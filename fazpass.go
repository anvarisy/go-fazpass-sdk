package fazpass

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	govalidator "github.com/go-playground/validator/v10"
)

type FazpassInterface interface {
	// Initialize(privatePath string, publicPath string, merchantKey string, url string) (*Fazpass, error)
	Check(email string, phone string, encData string) (*Data, error)
	EnrollDevice(email string, phone string, encData string) (*Data, error)
	ValidateDevice(fazpassId string, encData string) (*Data, error)
	RemoveDevice(fazpassId string, encData string) (*Data, error)
}

type Fazpass struct {
	PrivateKey  *rsa.PrivateKey
	PublicKey   *rsa.PublicKey
	MerchantKey string
	BaseUrl     string
}

func Initialize(privatePath string, publicPath string, merchantKey string, url string) (FazpassInterface, error) {
	var err error
	var privKey *rsa.PrivateKey
	var pubKey *rsa.PublicKey
	f := &Fazpass{}
	priv, errFile := os.ReadFile(privatePath)
	if errFile != nil {
		return f, errors.New("file not found")
	}
	privKey, _ = bytesToPrivateKey(priv)
	pub, errFile := os.ReadFile(publicPath)
	if errFile != nil {
		return f, errors.New("file not found")
	}
	pubKey, _ = bytesToPublicKey(pub)
	f.BaseUrl = url
	f.MerchantKey = merchantKey
	f.PrivateKey = privKey
	f.PublicKey = pubKey
	return f, err
}

func (f *Fazpass) Check(email string, phone string, encData string) (*Data, error) {
	check := &CheckRequest{
		Email: email,
		Phone: phone,
		Data:  encData,
	}
	data, err := sendToServer(f, check, "/check")
	return data, err
}

func (f *Fazpass) EnrollDevice(email string, phone string, encData string) (*Data, error) {
	enroll := &EnrollRequest{
		Email: email,
		Phone: phone,
		Data:  encData,
	}
	data, err := sendToServer(f, enroll, "/enroll")
	return data, err
}

func (f *Fazpass) ValidateDevice(fazpassId string, encData string) (*Data, error) {
	validate := &ValidateRequest{
		FazpassId: fazpassId,
		Data:      encData,
	}
	data, err := sendToServer(f, validate, "/validate")
	return data, err
}

func (f *Fazpass) RemoveDevice(fazpassId string, encData string) (*Data, error) {
	remove := &RemoveRequest{
		FazpassId: fazpassId,
		Data:      encData,
	}
	data, err := sendToServer(f, remove, "/remove")
	return data, err
}

func sendToServer(f *Fazpass, model interface{}, urls string) (*Data, error) {
	client := &http.Client{}
	data := &Data{}
	validator := govalidator.New()
	err := validator.Struct(model)
	if err != nil {
		return data, errors.New("parameter cannot be empty")
	}
	marshalledCheck, err := json.Marshal(model)
	if err != nil {
		return data, err
	}
	encrypted := encryptWithPublicKey(marshalledCheck, f.PublicKey)
	message := base64.StdEncoding.EncodeToString([]byte(encrypted))
	marshalledMessage, _ := json.Marshal(&Transmission{Message: message})
	request, err := http.NewRequest("POST", f.BaseUrl+urls, bytes.NewReader(marshalledMessage))
	if err != nil {
		return data, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+f.MerchantKey)
	response, err := client.Do(request)
	if err != nil {
		return data, nil
	}
	defer response.Body.Close()

	messageBodyResponse, _ := io.ReadAll(response.Body)
	transmissionResponse := &Transmission{}
	json.Unmarshal(messageBodyResponse, transmissionResponse)
	decryptMessage, err := base64.StdEncoding.DecodeString(string(transmissionResponse.Message))

	decrypted := decryptWithPrivateKey(decryptMessage, f.PrivateKey)
	json.Unmarshal(decrypted, data)
	fmt.Print(data.SessionId)
	if err != nil {
		return data, err
	}

	return data, nil
}

func bytesToPrivateKey(priv []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(priv)
	b := block.Bytes
	key, err := x509.ParsePKCS1PrivateKey(b)
	return key, err
}

// BytesToPublicKey bytes to public key
func bytesToPublicKey(pub []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pub)
	b := block.Bytes
	ifc, err := x509.ParsePKIXPublicKey(b)
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		log.Println("not ok")
	}
	return key, err
}

// EncryptWithPublicKey encrypts data with public key
func encryptWithPublicKey(msg []byte, pub *rsa.PublicKey) []byte {
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pub, msg)
	if err != nil {
		log.Println(err)
	}
	return ciphertext

}

// DecryptWithPrivateKey decrypts data with private key
func decryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		log.Println(err)
	}
	return plaintext
}
