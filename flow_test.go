package fazpass

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/anvarisy/go-fazpass-sdk/utils"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestWrappingData(t *testing.T) {
	t.Run("Wrapping data", func(t *testing.T) {
		f := Default()
		pub, _ := os.ReadFile("key.pub")
		pubKey, _ := utils.BytesToPublicKey(pub)
		_, err := f.WrappingData(pubKey, &Data{})
		assert.Equal(t, err, nil)
	})
}

func TestSendingData(t *testing.T) {
	t.Run("Sending data error", func(t *testing.T) {
		f := Default()
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		baseUrl := "http://localhost/check"
		httpmock.RegisterResponder("POST", baseUrl,
			httpmock.NewErrorResponder(errors.New("failed")))
		marshalled, _ := json.Marshal(&Data{})
		_, err := f.SendingData(baseUrl, marshalled, "MERCHANT_KEY")
		assert.NotEqual(t, err, nil)
	})
	t.Run("Sending data success", func(t *testing.T) {
		f := Default()
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		baseUrl := "http://localhost/check"
		httpmock.RegisterResponder("POST", baseUrl,
			httpmock.NewStringResponder(200, `{"id": 1, "session_id": ""}`))
		marshalled, _ := json.Marshal(&Data{})
		_, err := f.SendingData(baseUrl, marshalled, "MERCHANT_KEY")
		assert.Equal(t, err, nil)
	})
}

func TestExtractingData(t *testing.T) {
	t.Run("Extracting data", func(t *testing.T) {
		var (
			marshalled []byte
			encrypted  []byte
		)
		f := Default()
		pub, _ := os.ReadFile("key.pub")
		pubKey, _ := utils.BytesToPublicKey(pub)
		priv, _ := os.ReadFile("key.priv")
		privKey, _ := utils.BytesToPrivateKey(priv)
		marshalled, _ = json.Marshal(&Data{})
		encrypted, _ = utils.EncryptWithPublicKey(marshalled, pubKey)
		message := base64.StdEncoding.EncodeToString([]byte(encrypted))
		resp, _ := httpmock.NewJsonResponse(200, &Transmission{message})
		data, _ := f.ExtractingData(privKey, resp, &Data{})
		assert.Equal(t, data.SessionId, "")
	})
}
