package apiclient

// Broker describes the interactions with remote service brokers or similar
type Broker interface {
	Catalog()
	ProvisionAndBind(serviceID, planID string)
	Bind(serviceID, planID, instanceID, bindingID string)
	Unbind(serviceID, planID, instanceID, bindingID string)
	Deprovision(serviceID, planID, instanceID string)
}
