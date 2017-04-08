package config

import (
	"fmt"
	"time"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"gopkg.in/yaml.v2"
)

type FSConfig struct {
	path string
	fs   boshsys.FileSystem

	schema fsServiceInstances
}

type fsServiceInstances struct {
	ServiceInstances []fsServiceInstance `yaml:"service_instances"`
}

type fsServiceInstance struct {
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
	var schema fsServiceInstances

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

func (c FSConfig) ProvisionNewServiceInstance(id, name, serviceID, planID, brokerURL string) fsServiceInstance {
	_, inst := c.findOrCreateServiceInstanceByIDOrName(id, name)
	inst.ServiceID = serviceID
	inst.PlanID = planID
	inst.BrokerURL = brokerURL
	fmt.Printf("%#v\n", inst)
	fmt.Printf("%#v\n", c)
	c.Save()
	return *inst
}

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

func (c *FSConfig) findOrCreateServiceInstance(idOrName string) (int, *fsServiceInstance) {
	if idOrName != "" {
		for i, instance := range c.schema.ServiceInstances {
			if idOrName == instance.ID || idOrName == instance.Name {
				return i, &instance
			}
		}
	}

	return c.appendNewServiceInstanceWithID(idOrName)
}

func (c *FSConfig) findOrCreateServiceInstanceByIDOrName(id, name string) (int, *fsServiceInstance) {
	for i, instance := range c.schema.ServiceInstances {
		if id == instance.ID || (name != "" && name == instance.Name) {
			return i, &instance
		}
	}

	i, instance := c.appendNewServiceInstanceWithID(id)
	instance.Name = name
	return i, instance
}

func (c *FSConfig) appendNewServiceInstanceWithID(id string) (int, *fsServiceInstance) {
	instance := fsServiceInstance{ID: id, CreatedAt: time.Now()}
	c.schema.ServiceInstances = append(c.schema.ServiceInstances, instance)
	return len(c.schema.ServiceInstances) - 1, &c.schema.ServiceInstances[len(c.schema.ServiceInstances)-1]
}

func (c FSConfig) deepCopy() FSConfig {
	bytes, err := yaml.Marshal(c.schema)
	if err != nil {
		panic("serializing config schema")
	}

	var schema fsServiceInstances

	err = yaml.Unmarshal(bytes, &schema)
	if err != nil {
		panic("deserializing config schema")
	}

	return FSConfig{path: c.path, fs: c.fs, schema: schema}
}
