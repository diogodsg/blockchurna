package blockchain

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)


type Block struct {
	Id           string  `json:"id"`
	Index        int     `json:"index"`
	Timestamp    int64   `json:"timestamp"`
	Payload      Payload `json:"payload"`
	PreviousNode string  `json:"previous_node"`
}

type Payload struct {
	Presences 	[]Presence	`json:"presences"`
	Votes 		[]Vote     	`json:"votes"`
	Signature 	string		`json:"signature"`
}

type Presence struct {
	UserId		string	`json:"user_id"`
	Timestamp	string	`json:"timestamp"`
	Signature	string	`json:"signature"`
}

type Vote struct {
	Position	string	`json:"position"`
	Candidate	string 	`json:"candidate"`
	Hash		string	`json:"hash"`
}

func NewBlock(index int, payload Payload, previousNode string) *Block {
	block := &Block{
		Index:     index,
		Timestamp: time.Now().Unix(),
		Payload: payload,
		PreviousNode: previousNode,
	}
	block.Id = block.calculateHash()
	fmt.Println("block Id: ", block.Id)

	return block
}

func IsValidChain(blocks []*Block) bool {
	for i := 1; i < len(blocks); i++ {
        currentBlock := blocks[i]
        previousBlock := blocks[i-1]

        if currentBlock.PreviousNode != previousBlock.Id {
            return false
        }

        if currentBlock.calculateHash() != currentBlock.Id {
            return false
        }
    }
    return true
}

func (b *Block) calculateHash() string {
    payloadBytes, err := json.Marshal(b.Payload)
    if err != nil {
        log.Fatal("Error marshalling Payload: ", err)
    }

    record := fmt.Sprintf("%d%d%s%s", b.Index, b.Timestamp, b.PreviousNode, string(payloadBytes))

    hash := sha256.New()
    hash.Write([]byte(record))

    return fmt.Sprintf("%x", hash.Sum(nil))
}

func SerializePayload(payload *Payload) (string, error) {
	payloadCopy := *payload
	payloadCopy.Signature = "" // Exclude the signature field
	data, err := json.Marshal(payloadCopy)
	if err != nil {
		return "", fmt.Errorf("failed to serialize payload: %v", err)
	}
	return string(data), nil}

func ValidatePayload(payload Payload) error {
	for _, presence := range payload.Presences {
		key, err := LoadKeyFromJSON(presence.UserId)
		if err != nil {
			fmt.Printf("Error loading keys: %v\n", err)
			return err
		}

		err = validatePresence(presence, key)

		if err != nil {
			fmt.Printf("Invalid Signature for Presence %s: %v\n", presence.UserId, err)
			return err
		}
		fmt.Printf("Valid Signature for Presence %s\n", presence.UserId)
	}
	key, err := LoadKeyFromJSON("ballot")
	
	if err != nil {
		fmt.Printf("Error loading keys: %v\n", err)
		return err
	}
	data, err := SerializePayload(&payload)
	if err != nil {
		fmt.Printf("Error Serializing: %v\n", err)
		return err
	}

	blockData := strings.ReplaceAll(data, ",\"signature\":\"\"", "")

	fmt.Println(blockData)

	valid, err := verifySignature(key, blockData, payload.Signature)

	if err != nil || !valid {
		fmt.Printf("Error Verifying: %v\n", err)
		return err
	}

	return nil
}

func validatePresence(presence Presence, publicKey *rsa.PublicKey) error {
	dataToSign := presence.UserId + presence.Timestamp
	valid, err := verifySignature(publicKey, dataToSign, presence.Signature)
	 
	if err != nil || !valid  {
		fmt.Printf("Err: %v\n", err)
		return err
	}
	
	return nil
}

func verifySignature(rsaPublicKey *rsa.PublicKey, data string, signature string) (bool, error) {
	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return false, fmt.Errorf("failed to decode hex-encoded signature: %v", err)
	}

	// Hash the data
	hashed := sha256.Sum256([]byte(data))

	// Verify the signature
	err = rsa.VerifyPSS(rsaPublicKey, crypto.SHA256, hashed[:], signatureBytes, &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
	})
	if err != nil {
		return false, fmt.Errorf("signature verification failed: %v", err)
	}

	return true, nil
}