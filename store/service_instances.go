package config

import (
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
	ServiceInstances []*FSServiceInstance `yaml:"service_instances"`
}

type FSServiceInstance struct {
	ID        string             `yaml:"id"`
	Name      string             `yaml:"name"`
	ServiceID string             `yaml:"service_id"`
	PlanID    string             `yaml:"plan_id"`
	BrokerURL string             `yaml:"broker_url"`
	Bindings  []fsServiceBinding `yaml:"bindings"`
	CreatedAt time.Time          `yaml:"created_at"`
}

// ServiceBinding represents a binding with credentials
type fsServiceBinding struct {
	ID          string      `yaml:"id"`
	Name        string      `yaml:"name"`
	Credentials interface{} `yaml:"credentials"`
	CreatedAt   time.Time   `yaml:"created_at"`
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
func (c FSConfig) ProvisionNewServiceInstance(id, name, serviceID, planID, brokerURL string) {
	_, inst := c.findOrCreateServiceInstanceByIDOrName(id, name)
	inst.ServiceID = serviceID
	inst.PlanID = planID
	inst.BrokerURL = brokerURL
	c.Save()
}

// BindServiceInstance records a new bindingID
func (c FSConfig) BindServiceInstance(instanceID, bindingID, name string, credentials interface{}) {
	_, inst := c.findOrCreateServiceInstance(instanceID)
	binding := fsServiceBinding{
		ID:          bindingID,
		Name:        name,
		Credentials: credentials,
		CreatedAt:   time.Now(),
	}
	inst.Bindings = append(inst.Bindings, binding)
	c.Save()
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
func (c FSConfig) ServiceInstances() FSServiceInstances {
	return c.schema
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
