package apiclient

// Broker describes the interactions with remote service brokers or similar
type Broker interface {
	Catalog()
}
