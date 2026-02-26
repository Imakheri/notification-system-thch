package gateway

type SimulatedApiService interface {
	RandomizeHTTPStatus() int
}
