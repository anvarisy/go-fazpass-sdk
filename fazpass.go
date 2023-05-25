package fazpass

import (
	"crypto/rsa"
	"errors"
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
	Flow        FlowInterface
}

func Initialize(flow FlowInterface, privatePath string, publicPath string, merchantKey string, url string) (FazpassInterface, error) {
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
	f.Flow = flow
	return f, err
}

func (f *Fazpass) Check(email string, phone string, encData string) (*Data, error) {
	check := &CheckRequest{
		Email: email,
		Phone: phone,
		Data:  encData,
	}
	data := &Data{}
	validator := govalidator.New()
	err := validator.Struct(check)
	if err != nil {
		return data, errors.New("parameter cannot be empty")
	}

	wrappedMessage, err := f.Flow.wrapingData(f, check)
	if err != nil {
		return data, err
	}
	response, err := f.Flow.sendingData(f, wrappedMessage, "/check")
	if err != nil {
		return data, err
	}
	data, err = f.Flow.extractingData(f, response, data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (f *Fazpass) EnrollDevice(email string, phone string, encData string) (*Data, error) {
	enroll := &EnrollRequest{
		Email: email,
		Phone: phone,
		Data:  encData,
	}
	data := &Data{}
	validator := govalidator.New()
	err := validator.Struct(enroll)
	if err != nil {
		return data, errors.New("parameter cannot be empty")
	}

	wrappedMessage, err := f.Flow.wrapingData(f, enroll)
	if err != nil {
		return data, err
	}
	response, err := f.Flow.sendingData(f, wrappedMessage, "/enroll")
	if err != nil {
		return data, err
	}
	data, err = f.Flow.extractingData(f, response, data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (f *Fazpass) ValidateDevice(fazpassId string, encData string) (*Data, error) {
	validate := &ValidateRequest{
		FazpassId: fazpassId,
		Data:      encData,
	}
	data := &Data{}
	validator := govalidator.New()
	err := validator.Struct(validate)
	if err != nil {
		return data, errors.New("parameter cannot be empty")
	}

	wrappedMessage, err := f.Flow.wrapingData(f, validate)
	if err != nil {
		return data, err
	}
	response, err := f.Flow.sendingData(f, wrappedMessage, "/validate")
	if err != nil {
		return data, err
	}
	data, err = f.Flow.extractingData(f, response, data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func (f *Fazpass) RemoveDevice(fazpassId string, encData string) (*Data, error) {
	remove := &RemoveRequest{
		FazpassId: fazpassId,
		Data:      encData,
	}
	data := &Data{}
	validator := govalidator.New()
	err := validator.Struct(remove)
	if err != nil {
		return data, errors.New("parameter cannot be empty")
	}

	wrappedMessage, err := f.Flow.wrapingData(f, remove)
	if err != nil {
		return data, err
	}
	response, err := f.Flow.sendingData(f, wrappedMessage, "/remove")
	if err != nil {
		return data, err
	}
	data, err = f.Flow.extractingData(f, response, data)
	if err != nil {
		return data, err
	}
	return data, nil
}
