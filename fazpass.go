package fazpass

import (
	"crypto/rsa"
	"errors"
	"os"

	"github.com/anvarisy/go-fazpass-sdk/utils"

	govalidator "github.com/go-playground/validator/v10"
)

type FazpassInterface interface {
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
	privKey, _ = utils.BytesToPrivateKey(priv)
	pub, errFile := os.ReadFile(publicPath)
	if errFile != nil {
		return f, errors.New("file not found")
	}
	pubKey, _ = utils.BytesToPublicKey(pub)
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

	wrappedMessage, err := f.Flow.WrappingData(f.PublicKey, check)
	if err != nil {
		return data, err
	}
	response, err := f.Flow.SendingData(f.BaseUrl+"/check", wrappedMessage, f.MerchantKey)
	if err != nil {
		return data, err
	}
	data, err = f.Flow.ExtractingData(f.PrivateKey, response, data)
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

	wrappedMessage, err := f.Flow.WrappingData(f.PublicKey, enroll)
	if err != nil {
		return data, err
	}
	response, err := f.Flow.SendingData(f.BaseUrl+"/enroll", wrappedMessage, f.MerchantKey)
	if err != nil {
		return data, err
	}
	data, err = f.Flow.ExtractingData(f.PrivateKey, response, data)
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

	wrappedMessage, err := f.Flow.WrappingData(f.PublicKey, validate)
	if err != nil {
		return data, err
	}
	response, err := f.Flow.SendingData(f.BaseUrl+"/validate", wrappedMessage, f.MerchantKey)
	if err != nil {
		return data, err
	}
	data, err = f.Flow.ExtractingData(f.PrivateKey, response, data)
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

	wrappedMessage, err := f.Flow.WrappingData(f.PublicKey, remove)
	if err != nil {
		return data, err
	}
	response, err := f.Flow.SendingData(f.BaseUrl+"/remove", wrappedMessage, f.MerchantKey)
	if err != nil {
		return data, err
	}
	data, err = f.Flow.ExtractingData(f.PrivateKey, response, data)
	if err != nil {
		return data, err
	}
	return data, nil
}
