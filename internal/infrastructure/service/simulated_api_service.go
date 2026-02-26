package service

import (
	"math/rand"
	"net/http"
)

type SimulatedApiService struct{}

func NewSimulatedApiService() *SimulatedApiService {
	return &SimulatedApiService{}
}

func (sae *SimulatedApiService) RandomizeHTTPStatus() int {
	statuses := []int{
		http.StatusOK,
		http.StatusCreated,
		http.StatusBadRequest,
		http.StatusInternalServerError,
	}

	randomStatus := statuses[rand.Intn(len(statuses))]

	return randomStatus

}
