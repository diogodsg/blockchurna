package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type AggregatedData struct {
	City  string          `json:"city"`
	State string          `json:"state"`
	Votes []CandidateVote `json:"votes"`
	Section string 		  `json:"section"`
}

type CandidateVote struct {
	Candidate string `json:"candidate"`
	Votes     int    `json:"votes"`
}

func (bc *Blockchain) AggregateVotes() []*AggregatedData {
	aggregationMap := make(map[string]map[string]int)

	for _, block := range bc.Blocks {
		city := block.Payload.City
		state := block.Payload.State
		section := block.Payload.Section + " - " + block.Payload.Zone
		key := fmt.Sprintf("%s|%s|%s", city, state, section)

		if _, exists := aggregationMap[key]; !exists {
			aggregationMap[key] = make(map[string]int)
		}

		for _, vote := range block.Payload.Votes {
			aggregationMap[key][vote.Candidate]++
		}
	}

	var aggregatedData []*AggregatedData
	for key, candidateMap := range aggregationMap {
		parts := strings.Split(key, "|")
		city := parts[0]
		state := parts[1]
		section := parts[2]

		var votes []CandidateVote
		for candidate, count := range candidateMap {
			votes = append(votes, CandidateVote{
				Candidate: candidate,
				Votes:     count,
			})
		}

		aggregatedData = append(aggregatedData, &AggregatedData{
			City:  city,
			State: state,
			Votes: votes,
			Section: section,
		})
	}

	return aggregatedData
}

func (bc *Blockchain) VerifyVote(voterId string, userPin string, tsePin string) (*Vote, error) {
	hash := generateHash(voterId+userPin+tsePin)
	fmt.Println(voterId+userPin+tsePin)
	for _, block := range bc.Blocks {
		for _, presence := range block.Payload.Presences {
			if presence.UserId == voterId {
				for _, vote := range block.Payload.Votes {
					if vote.Hash == hash {
						return &vote, nil
					}
				}
			}
		}		
	}

	return nil, errors.New("vote not found")
}

func generateHash(data string) string {
	// Create a new SHA-256 hash object
	hash := sha256.New()

	// Write the data into the hash object
	hash.Write([]byte(data))

	// Get the resulting hash as a slice of bytes
	hashBytes := hash.Sum(nil)

	// Convert the hash bytes to a hexadecimal string and return
	return hex.EncodeToString(hashBytes)
}