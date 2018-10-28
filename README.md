# eden

Interact with any Open Service Broker API to discover/provision/bind/unbind/deprovision hundreds of different services.

* [Concourse CI](https://ci-ohio.starkandwayne.com/teams/cfcommunity/pipelines/eden)
* Pull requests will be automatically compiled and tested (see `test-pr` job)
* Discussions and CI notifications at [#eden channel](https://openservicebrokerapi.slack.com/messages/C6Y5A2N8Z/) on http://slack.openservicebrokerapi.org/

## Installation

For Ubuntu/Debian:

```shell
wget -q -O - https://raw.githubusercontent.com/starkandwayne/homebrew-cf/master/public.key | apt-key add -
echo "deb http://apt.starkandwayne.com stable main" | tee /etc/apt/sources.list.d/starkandwayne.list
apt-get update
apt-get install eden
```

For Mac OS using Homebrew:

```shell
brew install starkandwayne/cf/eden
```

From source using Golang:

```shell
go get -u github.com/starkandwayne/eden
```

## Usage

Use environment variables to target an Open Service Broker API:

```shell
export SB_BROKER_URL=https://mybroker.com
export SB_BROKER_USERNAME=username
export SB_BROKER_PASSWORD=password
```

To see the available services and plans:

```shell
eden catalog
```

To create (`provision`) a new service instance, and to generate a set of access credentials (`bind`):

```shell
export SB_INSTANCE=my-db-name
eden provision -s servicename -p planname
eden bind
```

To view the credentials for your binding:

```shell
eden credentials
```

To extract a single credentials, say a fully formed `uri`, you can use `eden credentials --attribute uri`:

For example, to provision a PostgreSQL service, generate bindings, and use them immediately with `psql`:

```shell
export SB_INSTANCE=pg1
eden provision -s postgresql96
eden bind
psql `eden creds -a uri`
```

### CLI flags and environment variables

In addition to using env vars, you can use CLI flags. See `eden -h` and `eden <command> -h` for more details.
