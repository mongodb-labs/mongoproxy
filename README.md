# mongoproxy

To grab dependencies:

	chmod 755 ./vendor.sh # only needs to be done once
	./vendor.sh

Note that currently, mongoproxy requires a specific fork of mgo because it requires some features not yet pulled into the main repo.

To run (requires mongod to be running on localhost:27017):

	chmod 755 ./start.sh # only needs to be done once
	./start.sh

To run tests:
	
	chmod 755 ./test.sh
	./test.sh <name of package to test>

To run integration tests:

	node tests/test <js file to test>
	node tests/test_dir <directory of files to test>

## BI Module Frontend

All operations assume that you are in the `modules/bi/frontend` directory, relative to the base path of the project. 
	
	cd modules/bi/frontend # from project directory

The frontend requires `node` and `npm` to be installed.

To grab dependencies:

	npm install

To run (requires mongod to be running on localhost:27017):

	npm start
