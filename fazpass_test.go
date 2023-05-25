package fazpass

import (
	"errors"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInitialize(t *testing.T) {
	f := Default()
	t.Run("Private key not found", func(t *testing.T) {
		_, err := Initialize(f, "key", "key.pub", "MERCHANT_KEY", "http://localhost:8080")
		assert.Equal(t, err.Error(), "file not found")
	})

	t.Run("Public key not found", func(t *testing.T) {
		f := Default()
		_, err := Initialize(f, "key.priv", "", "MERCHANT_KEY", "http://localhost:8080")
		assert.Equal(t, err.Error(), "file not found")
	})
}

func TestCheck(t *testing.T) {
	t.Run("Initialize failed", func(t *testing.T) {
		f := Default()
		_, err := Initialize(f, "key", "key.pub", "MERCHANT_KEY", "http://localhost:8080")
		assert.Equal(t, err.Error(), "file not found")
	})
	t.Run("Struct not valid", func(t *testing.T) {
		f := Default()
		fazpass, _ := Initialize(f, "key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")
		_, err := fazpass.Check("", "085811752000", "KOALA_PANDA")
		assert.Equal(t, err.Error(), "parameter cannot be empty")
	})

	t.Run("Wraping failed", func(t *testing.T) {
		f := new(FlowMock)
		fazpass, _ := Initialize(f, "key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")

		f.On("WrappingData", mock.Anything, mock.Anything).Return([]byte(""), errors.New("wraping failed"))
		_, err := fazpass.Check("anvarisy@gmail.com", "085811752000", "KOALA_PANDA")
		assert.Equal(t, err.Error(), "wraping failed")
	})
	t.Run("Sending data failed", func(t *testing.T) {
		f := new(FlowMock)
		fazpass, _ := Initialize(f, "key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")

		f.On("WrappingData", mock.Anything, mock.Anything).Return([]byte(""), nil)
		resp, _ := httpmock.NewJsonResponse(200, &Transmission{})
		f.On("SendingData", mock.Anything, mock.Anything, mock.Anything).Return(resp, errors.New("connection internet failed"))
		_, err := fazpass.Check("anvarisy@gmail.com", "085811752000", "KOALA_PANDA")
		assert.Equal(t, err.Error(), "connection internet failed")
	})
	t.Run("Extracting data failed", func(t *testing.T) {
		f := new(FlowMock)
		fazpass, _ := Initialize(f, "key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")

		f.On("WrappingData", mock.Anything, mock.Anything).Return([]byte(""), nil)
		resp, _ := httpmock.NewJsonResponse(200, &Transmission{})
		f.On("SendingData", mock.Anything, mock.Anything, mock.Anything).Return(resp, nil)
		f.On("ExtractingData", mock.Anything, mock.Anything, mock.Anything).Return(&Data{}, errors.New("extracting failed"))
		_, err := fazpass.Check("anvarisy@gmail.com", "085811752000", "KOALA_PANDA")
		assert.Equal(t, err.Error(), "extracting failed")
	})
	t.Run("Process complete", func(t *testing.T) {
		f := new(FlowMock)
		fazpass, _ := Initialize(f, "key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")

		f.On("WrappingData", mock.Anything, mock.Anything).Return([]byte(""), nil)
		resp, _ := httpmock.NewJsonResponse(200, &Transmission{})
		f.On("SendingData", mock.Anything, mock.Anything, mock.Anything).Return(resp, nil)
		f.On("ExtractingData", mock.Anything, mock.Anything, mock.Anything).Return(&Data{SessionId: "1"}, nil)
		d, _ := fazpass.Check("anvarisy@gmail.com", "085811752000", "KOALA_PANDA")
		assert.Equal(t, d.SessionId, "1")
	})
}

func TestEnroll(t *testing.T) {
	t.Run("Initialize failed", func(t *testing.T) {
		f := Default()
		_, err := Initialize(f, "key", "key.pub", "MERCHANT_KEY", "http://localhost:8080")
		assert.Equal(t, err.Error(), "file not found")
	})
	t.Run("Struct not valid", func(t *testing.T) {
		f := Default()
		fazpass, _ := Initialize(f, "key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")
		_, err := fazpass.EnrollDevice("", "085811752000", "KOALA_PANDA")
		assert.Equal(t, err.Error(), "parameter cannot be empty")
	})

	t.Run("Wraping failed", func(t *testing.T) {
		f := new(FlowMock)
		fazpass, _ := Initialize(f, "key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")

		f.On("WrappingData", mock.Anything, mock.Anything).Return([]byte(""), errors.New("wraping failed"))
		_, err := fazpass.EnrollDevice("anvarisy@gmail.com", "085811752000", "KOALA_PANDA")
		assert.Equal(t, err.Error(), "wraping failed")
	})
	t.Run("Sending data failed", func(t *testing.T) {
		f := new(FlowMock)
		fazpass, _ := Initialize(f, "key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")

		f.On("WrappingData", mock.Anything, mock.Anything).Return([]byte(""), nil)
		resp, _ := httpmock.NewJsonResponse(200, &Transmission{})
		f.On("SendingData", mock.Anything, mock.Anything, mock.Anything).Return(resp, errors.New("connection internet failed"))
		_, err := fazpass.EnrollDevice("anvarisy@gmail.com", "085811752000", "KOALA_PANDA")
		assert.Equal(t, err.Error(), "connection internet failed")
	})
	t.Run("Extracting data failed", func(t *testing.T) {
		f := new(FlowMock)
		fazpass, _ := Initialize(f, "key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")

		f.On("WrappingData", mock.Anything, mock.Anything).Return([]byte(""), nil)
		resp, _ := httpmock.NewJsonResponse(200, &Transmission{})
		f.On("SendingData", mock.Anything, mock.Anything, mock.Anything).Return(resp, nil)
		f.On("ExtractingData", mock.Anything, mock.Anything, mock.Anything).Return(&Data{}, errors.New("extracting failed"))
		_, err := fazpass.EnrollDevice("anvarisy@gmail.com", "085811752000", "KOALA_PANDA")
		assert.Equal(t, err.Error(), "extracting failed")
	})
	t.Run("Process complete", func(t *testing.T) {
		f := new(FlowMock)
		fazpass, _ := Initialize(f, "key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")

		f.On("WrappingData", mock.Anything, mock.Anything).Return([]byte(""), nil)
		resp, _ := httpmock.NewJsonResponse(200, &Transmission{})
		f.On("SendingData", mock.Anything, mock.Anything, mock.Anything).Return(resp, nil)
		f.On("ExtractingData", mock.Anything, mock.Anything, mock.Anything).Return(&Data{SessionId: "1"}, nil)
		d, _ := fazpass.EnrollDevice("anvarisy@gmail.com", "085811752000", "KOALA_PANDA")
		assert.Equal(t, d.SessionId, "1")
	})
}

func TestValidate(t *testing.T) {
	t.Run("Initialize failed", func(t *testing.T) {
		f := Default()
		_, err := Initialize(f, "key", "key.pub", "MERCHANT_KEY", "http://localhost:8080")
		assert.Equal(t, err.Error(), "file not found")
	})
	t.Run("Struct not valid", func(t *testing.T) {
		f := Default()
		fazpass, _ := Initialize(f, "key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")
		_, err := fazpass.ValidateDevice("", "KOALA_PANDA")
		assert.Equal(t, err.Error(), "parameter cannot be empty")
	})

	t.Run("Wraping failed", func(t *testing.T) {
		f := new(FlowMock)
		fazpass, _ := Initialize(f, "key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")

		f.On("WrappingData", mock.Anything, mock.Anything).Return([]byte(""), errors.New("wraping failed"))
		_, err := fazpass.ValidateDevice("FAZPASS_ID", "KOALA_PANDA")
		assert.Equal(t, err.Error(), "wraping failed")
	})
	t.Run("Sending data failed", func(t *testing.T) {
		f := new(FlowMock)
		fazpass, _ := Initialize(f, "key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")

		f.On("WrappingData", mock.Anything, mock.Anything).Return([]byte(""), nil)
		resp, _ := httpmock.NewJsonResponse(200, &Transmission{})
		f.On("SendingData", mock.Anything, mock.Anything, mock.Anything).Return(resp, errors.New("connection internet failed"))
		_, err := fazpass.ValidateDevice("FAZPASS_ID", "KOALA_PANDA")
		assert.Equal(t, err.Error(), "connection internet failed")
	})
	t.Run("Extracting data failed", func(t *testing.T) {
		f := new(FlowMock)
		fazpass, _ := Initialize(f, "key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")

		f.On("WrappingData", mock.Anything, mock.Anything).Return([]byte(""), nil)
		resp, _ := httpmock.NewJsonResponse(200, &Transmission{})
		f.On("SendingData", mock.Anything, mock.Anything, mock.Anything).Return(resp, nil)
		f.On("ExtractingData", mock.Anything, mock.Anything, mock.Anything).Return(&Data{}, errors.New("extracting failed"))
		_, err := fazpass.ValidateDevice("FAZPASS_ID", "KOALA_PANDA")
		assert.Equal(t, err.Error(), "extracting failed")
	})
	t.Run("Process complete", func(t *testing.T) {
		f := new(FlowMock)
		fazpass, _ := Initialize(f, "key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")

		f.On("WrappingData", mock.Anything, mock.Anything).Return([]byte(""), nil)
		resp, _ := httpmock.NewJsonResponse(200, &Transmission{})
		f.On("SendingData", mock.Anything, mock.Anything, mock.Anything).Return(resp, nil)
		f.On("ExtractingData", mock.Anything, mock.Anything, mock.Anything).Return(&Data{SessionId: "1"}, nil)
		d, _ := fazpass.ValidateDevice("FAZPASS_ID", "KOALA_PANDA")
		assert.Equal(t, d.SessionId, "1")
	})
}

func TestRemove(t *testing.T) {
	t.Run("Initialize failed", func(t *testing.T) {
		f := Default()
		_, err := Initialize(f, "key", "key.pub", "MERCHANT_KEY", "http://localhost:8080")
		assert.Equal(t, err.Error(), "file not found")
	})
	t.Run("Struct not valid", func(t *testing.T) {
		f := Default()
		fazpass, _ := Initialize(f, "key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")
		_, err := fazpass.RemoveDevice("", "KOALA_PANDA")
		assert.Equal(t, err.Error(), "parameter cannot be empty")
	})

	t.Run("Wraping failed", func(t *testing.T) {
		f := new(FlowMock)
		fazpass, _ := Initialize(f, "key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")

		f.On("WrappingData", mock.Anything, mock.Anything).Return([]byte(""), errors.New("wraping failed"))
		_, err := fazpass.RemoveDevice("FAZPASS_ID", "KOALA_PANDA")
		assert.Equal(t, err.Error(), "wraping failed")
	})
	t.Run("Sending data failed", func(t *testing.T) {
		f := new(FlowMock)
		fazpass, _ := Initialize(f, "key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")

		f.On("WrappingData", mock.Anything, mock.Anything).Return([]byte(""), nil)
		resp, _ := httpmock.NewJsonResponse(200, &Transmission{})
		f.On("SendingData", mock.Anything, mock.Anything, mock.Anything).Return(resp, errors.New("connection internet failed"))
		_, err := fazpass.RemoveDevice("FAZPASS_ID", "KOALA_PANDA")
		assert.Equal(t, err.Error(), "connection internet failed")
	})
	t.Run("Extracting data failed", func(t *testing.T) {
		f := new(FlowMock)
		fazpass, _ := Initialize(f, "key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")

		f.On("WrappingData", mock.Anything, mock.Anything).Return([]byte(""), nil)
		resp, _ := httpmock.NewJsonResponse(200, &Transmission{})
		f.On("SendingData", mock.Anything, mock.Anything, mock.Anything).Return(resp, nil)
		f.On("ExtractingData", mock.Anything, mock.Anything, mock.Anything).Return(&Data{}, errors.New("extracting failed"))
		_, err := fazpass.RemoveDevice("FAZPASS_ID", "KOALA_PANDA")
		assert.Equal(t, err.Error(), "extracting failed")
	})
	t.Run("Process complete", func(t *testing.T) {
		f := new(FlowMock)
		fazpass, _ := Initialize(f, "key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")

		f.On("WrappingData", mock.Anything, mock.Anything).Return([]byte(""), nil)
		resp, _ := httpmock.NewJsonResponse(200, &Transmission{})
		f.On("SendingData", mock.Anything, mock.Anything, mock.Anything).Return(resp, nil)
		f.On("ExtractingData", mock.Anything, mock.Anything, mock.Anything).Return(&Data{SessionId: "1"}, nil)
		d, _ := fazpass.RemoveDevice("FAZPASS_ID", "KOALA_PANDA")
		assert.Equal(t, d.SessionId, "1")
	})
}
