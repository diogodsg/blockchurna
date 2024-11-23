package blockchain

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
)

// KeyEntry represents the structure of the key data in the JSON file
type KeyEntry struct {
	KeyID           string `json:"key_id"`
	PublicKeyBase64 string `json:"public_key_base64"`
}

func LoadKeyFromJSON(userId string) (*rsa.PublicKey, error) {
	data, err := ioutil.ReadFile("./keys.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	var keyEntries []KeyEntry
	err = json.Unmarshal(data, &keyEntries)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	var entry KeyEntry
	for _, key := range keyEntries {
		if key.KeyID == userId {
			entry = key
		}
	}


	pemData, err := base64.StdEncoding.DecodeString(entry.PublicKeyBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 for key ID '%s': %v", entry.KeyID, err)
	}
	
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("invalid PEM block for key ID '%s'", entry.KeyID)
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key for key ID '%s': %v", entry.KeyID, err)
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("the key is not an RSA public key")
	}

	return rsaPub, nil
}