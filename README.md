# eden

Interact with any Open Service Broker API to discover/provision/bind/unbind/deprovision hundreds of different services.

## Usage

Use environment variables to target an Open Service Broker API:

```
export EDEN_BROKER_URL=https://mybroker.com
export EDEN_BROKER_CLIENT=username
export EDEN_BROKER_CLIENT_SECRET=password
```

To see the available services and plans:

```
eden catalog
```

To create (`provision`) a new service instance, and to generate a set of access credentials (`bind`):

```
export EDEN_INSTANCE=my-db-name
eden provision -s servicename -p planname
eden bind
```

Instead of env vars, you can use CLI flags. See `eden -h` for more details.

## Install

```
go get -u github.com/starkandwayne/eden
```
