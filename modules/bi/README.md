# BI Module

A business intelligence module for MongoProxy. Passes the request through the module pipeline, then analyzes `insert` requests to see if they match one of its configured rules, to collect metrics.

## Usage

	name: bi

## Configuration

The configuration defines rules and behavior on collecting data from requests, and where to store the data afterwards.

It contains two fields: a `connection` and `rules`. The connection field is an object:

	connection: {
		addresses: (array of strings) - contains addresses of servers to connect to. If no port is provided, will default to 27017.
		direct: (optional boolean) - determines whether to establish connections only with the specified server, or to obtain cluster information to connect with other servers.
		timeout: (optional integer) - the amount of time to wait for the server(s) to respond on connecting, in nanoseconds, before returning an error. If set to 0, then there is no timeout. Defaults to 10 seconds.
		auth: (optional object) {
			database: (string) - the default database that will be connected to for authentication
			username: (string)
			password: (string)
		}
	}

The `rules` field is an array of objects, with each object in the array having the following fields:

	{
		origin: (string) - the origin collection name to intercept from the original user operation.
		prefix: (string) - a prefix for the namespace where the metric documents will be stored. The metric documents may be stored over multiple collections, but each of those collections will begin with this prefix.
		timeGranularity: (an array of strings) - determines the time granularities in which the collected data is stored.
		valueField: (string) - the field that is analyzed in the document.
		timeField: (optional string) - a field with a date that determines the timestamp of the particular document. If no timeField is present, a timestamp will be generated based on the time the request was processed by the BI Module.
	}

#### Time Granularity

The possible time granularities (for the `timeGranularity` field in a rule) are as follows:

* `M` - monthly
* `D` - daily
* `h` - hourly
* `m` - minutely
* `s` - secondly

### Example Configuration

	{
	    connection: {
	        addresses: ["localhost"],
	        database: "test"
	    }
	    rules: [ 
	        {
	            origin: "db.foo",
	            prefix: "db.foo-metrics"
	            timeGranularity: ["M", "D", "h", "m", "s"],
	            valueField: "users",
	            timeField: "created"
	        }
	    ]
	}
