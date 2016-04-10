# SQLProxy Module

A proxying module for SQLProxy, that sends the request to the configured SQLProxy instance(s).

## Configuration

The configuration defines the server(s) that the module connects to. It has the following fields:

	{
		addresses: contains addresses of the SQLProxy
		schema: contains location of SQLProxy schema file
	}

## Example
		{
			"name": "sqlproxy",
			"config": {
				"address": "127.0.0.1:27017",
				"schema": "testdata/blackbox.yml",
			},
		}
