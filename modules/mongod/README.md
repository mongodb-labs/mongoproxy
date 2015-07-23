# Mongod Module

A proxying module for MongoProxy, that sends the request to the configured mongod instance(s).

## Configuration

The configuration defines the server(s) that the module connects to. It has the following fields:

	{
		addresses: (array of strings) - contains addresses of servers to connect to. If no port is provided, will default to 27017.
		direct: (optional boolean) - determines whether to establish connections only with the specified server, or to obtain cluster information to connect with other servers.
		timeout: (optional integer) - the amount of time to wait for the server(s) to respond on connecting, in nanoseconds, before returning an error. If set to 0, then there is no timeout. Defaults to 10 seconds.
		auth: (optional object) {
			database: (string) - the default database that will be connected to for authentication
			username: (string)
			password: (string)
		}
	}

## Example

	{
		"addresses": [
			"localhost:27017"
		]
	}
