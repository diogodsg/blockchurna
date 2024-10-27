package p2p

import (
	"crypto/rand"
	"os"

	"github.com/libp2p/go-libp2p/core/crypto"
)

func savePrivateKeyToFile(privateKey crypto.PrivKey, filename string) error {
    keyBytes, err := crypto.MarshalPrivateKey(privateKey)
    if err != nil {
        return err
    }
    return os.WriteFile(filename, keyBytes, 0600)
}

func loadPrivateKeyFromFile(filename string) (crypto.PrivKey, error) {
    keyBytes, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return crypto.UnmarshalPrivateKey(keyBytes)
}

func generateKeyPair() (crypto.PrivKey, crypto.PubKey, error) {
    return crypto.GenerateKeyPairWithReader(crypto.Ed25519, 2048, rand.Reader)
}

func GetNodeKey(port string) crypto.PrivKey {
	privateKeyFile := port + ".key"

    var priv crypto.PrivKey

    // Try loading existing keypair
    if _, err := os.Stat(privateKeyFile); os.IsNotExist(err) {
        priv, _, err = generateKeyPair() // generate new keypair
        if err != nil {
            panic(err)
        }

        // Save keypair to file
        if err := savePrivateKeyToFile(priv, privateKeyFile); err != nil {
            panic(err)
        }

        println("New keypair generated and saved.")
    } else {
        priv, err = loadPrivateKeyFromFile(privateKeyFile)
        if err != nil {
            panic(err)
        }
        println("Loaded existing keypair from file.")
    }

	return priv
}