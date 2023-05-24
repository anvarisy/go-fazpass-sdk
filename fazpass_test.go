package fazpass

import (
	"crypto/rsa"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInitialize(t *testing.T) {
	t.Run("Private key not found", func(t *testing.T) {
		_, err := Initialize("key", "key.pub", "MERCHANT_KEY", "http://localhost:8080")
		assert.Equal(t, err.Error(), "file not found")
	})

	t.Run("Public key not found", func(t *testing.T) {
		_, err := Initialize("key.priv", "", "MERCHANT_KEY", "http://localhost:8080")
		assert.Equal(t, err.Error(), "file not found")
	})
}

func TestCheck(t *testing.T) {
	t.Run("Initialize failed", func(t *testing.T) {
		_, err := Initialize("key", "key.pub", "MERCHANT_KEY", "http://localhost:8080")
		assert.Equal(t, err.Error(), "file not found")
	})
	t.Run("Check failed", func(t *testing.T) {
		f, _ := Initialize("key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")
		_, err := f.Check("", "085811752000", "KOALA_PANDA")
		assert.NotEqual(t, err, nil)
	})
	t.Run("Check success", func(t *testing.T) {
		f, _ := Initialize("key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")
		var fm = new(FazpassMock)
		var privKey *rsa.PrivateKey
		priv, _ := os.ReadFile("key.priv")
		privKey, _ = bytesToPrivateKey(priv)
		fm.On("Check", mock.Anything, mock.Anything, mock.Anything).Return(&Data{}, nil)
		data, _ := f.Check("anvarisy@gmail.com", "085811752000", "KOALA_PANDA")
		dataMarshal, _ := json.Marshal(data)
		_, err := decryptWithPrivateKey(dataMarshal, privKey)
		assert.NotEqual(t, err, nil)
	})

}

func TestEnroll(t *testing.T) {
	t.Run("Initialize failed", func(t *testing.T) {
		_, err := Initialize("key", "key.pub", "MERCHANT_KEY", "http://localhost:8080")
		assert.Equal(t, err.Error(), "file not found")
	})
	t.Run("Enroll failed", func(t *testing.T) {
		f, _ := Initialize("key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")
		_, err := f.EnrollDevice("", "085811752000", "KOALA_PANDA")
		assert.NotEqual(t, err, nil)
	})
	t.Run("Enroll success", func(t *testing.T) {
		f, _ := Initialize("key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")
		var fm = new(FazpassMock)
		fm.On("EnrollDevice", mock.Anything, mock.Anything, mock.Anything).Return(&Data{}, nil)
		_, err := f.EnrollDevice("anvarisy@gmail.com", "085811752000", "KOALA_PANDA")
		assert.Equal(t, err, nil)
	})
}

func TestValidate(t *testing.T) {
	t.Run("Initialize failed", func(t *testing.T) {
		_, err := Initialize("key", "key.pub", "MERCHANT_KEY", "http://localhost:8080")
		assert.Equal(t, err.Error(), "file not found")
	})
	t.Run("Validate failed", func(t *testing.T) {
		f, _ := Initialize("key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")
		_, err := f.ValidateDevice("", "ENC_DATA")
		assert.NotEqual(t, err, nil)
	})
	t.Run("Validate success", func(t *testing.T) {
		f, _ := Initialize("key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")
		var fm = new(FazpassMock)
		fm.On("ValidateDevice", mock.Anything, mock.Anything).Return(&Data{}, nil)
		_, err := f.ValidateDevice("FAZPASS_ID", "ENC_DATA")
		assert.Equal(t, err, nil)
	})
}

func TestRemove(t *testing.T) {
	t.Run("Initialize failed", func(t *testing.T) {
		_, err := Initialize("key", "key.pub", "MERCHANT_KEY", "http://localhost:8080")
		assert.Equal(t, err.Error(), "file not found")
	})
	t.Run("Remove failed", func(t *testing.T) {
		f, _ := Initialize("key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")
		_, err := f.RemoveDevice("", "ENC_DATA")
		assert.NotEqual(t, err, nil)
	})
	t.Run("Remove success", func(t *testing.T) {
		f, _ := Initialize("key.priv", "key.pub", "MERCHANT_KEY", "http://localhost:8080")
		var fm = new(FazpassMock)
		fm.On("RemoveDevice", mock.Anything, mock.Anything).Return(&Data{}, nil)
		_, err := f.RemoveDevice("FAZPASS_ID", "ENC_DATA")
		assert.Equal(t, err, nil)
	})
}

/*
	func TestGenerateKeyPair(t *testing.T) {
		t.Run("No Error", func(t *testing.T) {
			_, err := generateKeyPair()
			assert.Equal(t, nil, err)
		})
	}

	func TestEncrypt(t *testing.T) {
		data, _ := os.ReadFile("data.txt")
		key, _ := os.ReadFile("key.pub")
		t.Run("Check Encrypt", func(t *testing.T) {
			pubKey, _ := bytesToPublicKey(key)
			encrypted := encryptWithPublicKey(data, pubKey)
			f, _ := os.Create("en.txt")
			defer f.Close()
			v := base64.StdEncoding.EncodeToString([]byte(encrypted))
			f.WriteString(v)
		})
	}

	func TestDecrypt(t *testing.T) {
		enc, _ := os.ReadFile("en.txt")
		v, _ := base64.StdEncoding.DecodeString(string(enc))
		key, _ := os.ReadFile("key.priv")
		log.Println(string(key))
		t.Run("Check Decrypt", func(t *testing.T) {
			privKey, _ := bytesToPrivateKey(key)
			decrypted := decryptWithPrivateKey(v, privKey)
			assert.Equal(t, "koala panda", string(decrypted))
		})
	}
*/
