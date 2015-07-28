# BI Module Frontend

A web application to visualize the data collected from the BI Module.

## Building and running

All operations assume that you are in the `modules/bi/frontend` directory, relative to the base path of the project (the directory of this README).

The application is written in Go and Javascript. It requires `node` and `npm` to be installed.

To grab dependencies:

	npm install

To run (requires mongod to be running on localhost:27017):

	npm start
