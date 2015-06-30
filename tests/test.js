// janky script to run the server and x clients

var childProcess = require('child_process');
var fs = require('fs');
var path = process.cwd();

var server = childProcess.spawn('go', [
	'run',
	'./tests/main/mongod-server.go',
	'-port=' + 8000,
	'-logLevel=' + 1
], {
	stdio: 'inherit'
});

testFiles = fs.readdirSync(process.argv[2])
i = 0;

setTimeout(function() {

	var id = setInterval(function() {

		console.log("Testing: " + testFiles[i])
		var shell = childProcess.spawnSync('mongo', [
			process.argv[2] + '/' + testFiles[i],
			'--port=8000'
		], {
			stdio: 'inherit'
		})
		i++;

		if (i >= testFiles.length) {
			clearInterval(id)
		}
	}, 100)

}, 2000)
