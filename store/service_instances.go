package config

import (
	"encoding/json"
	"os"
	"time"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"gopkg.in/yaml.v2"
)

type FSConfig struct {
	path string
	fs   boshsys.FileSystem

	schema FSServiceInstances
}

type FSServiceInstances struct {
	ServiceInstances []*FSServiceInstance `yaml:"service_instances" json:"service_instances"`
}

type FSServiceInstance struct {
	ID          string             `yaml:"id"           json:"id"`
	Name        string             `yaml:"name"         json:"name"`
	ServiceID   string             `yaml:"service_id"   json:"service_id"`
	ServiceName string             `yaml:"service_name" json:"service_name"`
	PlanID      string             `yaml:"plan_id"      json:"plan_id"`
	PlanName    string             `yaml:"plan_name"    json:"plan_name"`
	BrokerURL   string             `yaml:"broker_url"   json:"broker_url"`
	Bindings    []fsServiceBinding `yaml:"bindings"     json:"bindings"`
	CreatedAt   time.Time          `yaml:"created_at"   json:"created_at"`
}

// ServiceBinding represents a binding with credentials
type fsServiceBinding struct {
	ID          string    `yaml:"id"             json:"id"`
	Name        string    `yaml:"name"           json:"name"`
	Credentials string    `yaml:"credentials"    json:"credentials"`
	CreatedAt   time.Time `yaml:"created_at"     json:"created_at"`
}

func NewFSConfigFromPath(path string, fs boshsys.FileSystem) (FSConfig, error) {
	var schema FSServiceInstances

	absPath, err := fs.ExpandPath(path)
	if err != nil {
		return FSConfig{}, err
	}

	if fs.FileExists(absPath) {
		bytes, err := fs.ReadFile(absPath)
		if err != nil {
			return FSConfig{}, bosherr.WrapErrorf(err, "Reading config '%s'", absPath)
		}

		err = yaml.Unmarshal(bytes, &schema)
		if err != nil {
			return FSConfig{}, bosherr.WrapError(err, "Unmarshalling config")
		}
	}
	return FSConfig{path: absPath, fs: fs, schema: schema}, nil
}

// ProvisionNewServiceInstance initialize new FSServiceInstance
func (c FSConfig) ProvisionNewServiceInstance(id, name, serviceID, serviceName, planID, planName, brokerURL string) {
	_, inst := c.findOrCreateServiceInstanceByIDOrName(id, name)
	inst.ServiceID = serviceID
	inst.ServiceName = serviceName
	inst.PlanID = planID
	inst.PlanName = planName
	inst.BrokerURL = brokerURL
	c.Save()
}

// FindServiceInstance returns a copy of a service instance record
func (c FSConfig) FindServiceInstance(idOrName string) FSServiceInstance {
	_, inst := c.findOrCreateServiceInstance(idOrName)
	return *inst
}

// RenameServiceInstance updates the .Name of a service instance
func (c FSConfig) RenameServiceInstance(idOrName, newName string) {
	_, inst := c.findOrCreateServiceInstance(idOrName)
	inst.Name = newName
	c.Save()
}

// BindServiceInstance records a new bindingID
func (c FSConfig) BindServiceInstance(instanceID, bindingID, name string, rawCredentials interface{}) (err error) {
	_, inst := c.findOrCreateServiceInstance(instanceID)

	credentialsStr, err := json.Marshal(rawCredentials)
	if err != nil {
		return bosherr.WrapError(err, "Marshalling raw credentials")
	}

	binding := fsServiceBinding{
		ID:          bindingID,
		Name:        name,
		Credentials: string(credentialsStr),
		CreatedAt:   time.Now(),
	}
	inst.Bindings = append(inst.Bindings, binding)
	return c.Save()
}

// UnbindServiceInstance removes record of a binding
func (c FSConfig) UnbindServiceInstance(instanceID, bindingNameOrID string) {
	_, inst := c.findOrCreateServiceInstance(instanceID)
	bindings := []fsServiceBinding{}
	for _, binding := range inst.Bindings {
		if binding.ID != bindingNameOrID && binding.Name != bindingNameOrID {
			bindings = append(bindings, binding)
		}
	}
	inst.Bindings = bindings
	c.Save()
}

// DeprovisionServiceInstance removes record of an instance
func (c FSConfig) DeprovisionServiceInstance(instanceNameOrID string) {
	instances := []*FSServiceInstance{}
	for _, instance := range c.schema.ServiceInstances {
		if instance.ID != instanceNameOrID && instance.Name != instanceNameOrID {
			instances = append(instances, instance)
		}
	}
	c.schema.ServiceInstances = instances
	c.Save()
}

// ServiceInstances returns the list of service instances created locally
func (c FSConfig) ServiceInstances() []*FSServiceInstance {
	return c.schema.ServiceInstances
}

// Save configuration/data to file
func (c FSConfig) Save() error {
	bytes, err := yaml.Marshal(c.schema)
	if err != nil {
		return bosherr.WrapError(err, "Marshalling config")
	}

	err = c.fs.WriteFile(c.path, bytes)
	if err != nil {
		return bosherr.WrapErrorf(err, "Writing config '%s'", c.path)
	}

	os.Chmod(c.path, 0600)

	return nil
}

func (c *FSConfig) findOrCreateServiceInstance(idOrName string) (int, *FSServiceInstance) {
	if idOrName != "" {
		for i, instance := range c.schema.ServiceInstances {
			if idOrName == instance.ID || idOrName == instance.Name {
				return i, instance
			}
		}
	}

	return c.appendNewServiceInstanceWithID(idOrName)
}

func (c *FSConfig) findOrCreateServiceInstanceByIDOrName(id, name string) (int, *FSServiceInstance) {
	for i, instance := range c.schema.ServiceInstances {
		if id == instance.ID || (name != "" && name == instance.Name) {
			return i, instance
		}
	}

	i, instance := c.appendNewServiceInstanceWithID(id)
	instance.Name = name
	return i, instance
}

func (c *FSConfig) appendNewServiceInstanceWithID(id string) (int, *FSServiceInstance) {
	instance := &FSServiceInstance{ID: id, CreatedAt: time.Now()}
	c.schema.ServiceInstances = append(c.schema.ServiceInstances, instance)
	index := len(c.schema.ServiceInstances) - 1
	return index, instance
}

func (c FSConfig) deepCopy() FSConfig {
	bytes, err := yaml.Marshal(c.schema)
	if err != nil {
		panic("serializing config schema")
	}

	var schema FSServiceInstances

	err = yaml.Unmarshal(bytes, &schema)
	if err != nil {
		panic("deserializing config schema")
	}

	return FSConfig{path: c.path, fs: c.fs, schema: schema}
}

// Credentials fixes any map[interface{}]interface{} into map[string]interface{}
// as expected by JSON marshalling
func (b fsServiceBinding) CredentialsJSON() (out map[string]interface{}, err error) {

	err = json.Unmarshal([]byte(b.Credentials), &out)
	if err != nil {
		return nil, bosherr.WrapError(err, "Unmarshalling raw credentials")
	}
	return
}

func deepCopy(rawInput interface{}) (output map[string]interface{}) {
	output = map[string]interface{}{}
	switch input := rawInput.(type) {
	case map[string]interface{}:
		for key, rawValue := range input {
			switch value := rawValue.(type) {
			case map[interface{}]interface{}:
				output[key] = deepCopy(value)
			default:
				output[key] = value
			}
		}
	case map[interface{}]interface{}:
		for key, rawValue := range input {
			switch value := rawValue.(type) {
			case map[interface{}]interface{}:
				output[key.(string)] = deepCopy(value)
			default:
				output[key.(string)] = value
			}
		}
	}
	return
}
