package service

import (
	"errors"
	"math/rand"
	"net/http"
	"time"
)

type SimulatedApiService struct{}

func NewSimulatedApiService() *SimulatedApiService {
	return &SimulatedApiService{}
}

func (sae *SimulatedApiService) RandomizeHTTPStatus() (int, error) {
	statuses := []int{
		http.StatusOK,
		http.StatusCreated,
		http.StatusBadRequest,
		http.StatusInternalServerError,
	}

	randomStatus := statuses[rand.Intn(len(statuses))]
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		if randomStatus == http.StatusOK || randomStatus == http.StatusCreated {
			return randomStatus, nil
		}
		time.Sleep(time.Duration(i) * time.Second)
	}
	return randomStatus, errors.New("an error occurred while trying to send")
}
