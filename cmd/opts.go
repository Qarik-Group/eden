package cmd

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	edenstore "github.com/starkandwayne/eden/store"
)

// InstanceOpts describes a target service instance
type InstanceOpts struct {
	NameOrID string `short:"i" long:"instance" description:"Service instance name/ID" env:"SB_INSTANCE"`
}

// BrokerOpts describes subset of flags/options for selecting target service broker API
type BrokerOpts struct {
	URLOpt          string `long:"url"           description:"Open Service Broker URL"                env:"SB_BROKER_URL" required:"true"`
	ClientOpt       string `long:"client"        description:"Override username or UAA client"        env:"SB_BROKER_USERNAME" required:"true"`
	ClientSecretOpt string `long:"client-secret" description:"Override password or UAA client secret" env:"SB_BROKER_PASSWORD" required:"true"`
	APIVersion      string `long:"api-version"   description:"API version request to pass to backend broker" env:"SB_BROKER_API_VERSION" default:"2.13"`
}

// EdenOpts describes the flags/options for the CLI
type EdenOpts struct {
	Version bool `short:"v" long:"version" description:"Show version"`

	// Slice of bool will append 'true' each time the option
	// is encountered (can be set multiple times, like -vvv)
	Verbose []bool `long:"verbose" description:"Show verbose debug information" env:"EDEN_TRACE"`

	ConfigPathOpt string `long:"config" description:"Config file path" env:"EDEN_CONFIG" default:"~/.eden/config"`

	Instance InstanceOpts `group:"Service Instance Options"`
	Broker   BrokerOpts   `group:"Broker Options"`

	// Broker API commands
	Catalog     CatalogOpts     `command:"catalog" alias:"cat" alias:"inventory" alias:"inv" description:"Show available service catalog"`
	Provision   ProvisionOpts   `command:"provision" alias:"p" description:"Create new service instance"`
	Bind        BindOpts        `command:"bind" alias:"b" description:"Generate credentials for service instance"`
	Unbind      UnbindOpts      `command:"unbind" alias:"u" description:"Remove credentials for service instance"`
	Deprovision DeprovisionOpts `command:"deprovision" alias:"d" description:"Destroy service instance"`

	// Local data commands
	Services    ServicesOpts    `command:"services" alias:"s" description:"List service instances (stored in config file)"`
	Credentials CredentialsOpts `command:"credentials" alias:"creds" alias:"c" description:"Display binding credentials (stored in config file)"`
	Rename      RenameOpts      `command:"rename" description:"Rename service instance (stored in config file)"`
}

// Opts carries all the user provided options (from flags or env vars)
var Opts EdenOpts

// TODO: need to move this into separate struct; bosh-cli has cmd.BasicDeps
func (opts EdenOpts) fs() boshsys.FileSystem {
	logger := boshlog.NewLogger(boshlog.LevelInfo)
	return boshsys.NewOsFileSystem(logger)
}

func (opts EdenOpts) config() edenstore.FSConfig {
	config, err := edenstore.NewFSConfigFromPath(opts.ConfigPathOpt, opts.fs())
	if err != nil {
		panic(err)
	}

	return config
}
