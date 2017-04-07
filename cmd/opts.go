package cmd

// BrokerOpts describes subset of flags/options for selecting target service broker API
type BrokerOpts struct {
	URLOpt          string `long:"url"           description:"Open Service Broker URL"                env:"EDEN_BROKER_URL" required:"true"`
	ClientOpt       string `long:"client"        description:"Override username or UAA client"        env:"EDEN_BROKER_CLIENT" required:"true"`
	ClientSecretOpt string `long:"client-secret" description:"Override password or UAA client secret" env:"EDEN_BROKER_CLIENT_SECRET" required:"true"`
}

// EdenOpts describes the flags/options for the CLI
type EdenOpts struct {
	// Slice of bool will append 'true' each time the option
	// is encountered (can be set multiple times, like -vvv)
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information" env:"EDEN_TRACE"`

	InstanceName string `short:"i" long:"instance" description:"Service instance name/ID" env:"EDEN_SERVICE"`

	Broker BrokerOpts `group:"Broker Options"`

	Catalog   CatalogOpts   `command:"catalog" alias:"c" alias:"inventory" alias:"inv" description:"Show available service catalog"`
	Provision ProvisionOpts `command:"provision" alias:"p" description:"Create new service instance"`
	Bind      BindOpts      `command:"bind" alias:"b" description:"Generate credentials for service instance"`
	Unbind    UnbindOpts    `command:"unbind" alias:"u" description:"Remove credentials for service instance"`
}

// Opts carries all the user provided options (from flags or env vars)
var Opts EdenOpts
