package blockchain

import (
	"fmt"
	"strings"
)

type AggregatedData struct {
	City  string          `json:"city"`
	State string          `json:"state"`
	Votes []CandidateVote `json:"votes"`
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
		key := fmt.Sprintf("%s|%s", city, state)

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
		})
	}

	return aggregatedData
}