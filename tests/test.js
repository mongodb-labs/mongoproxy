// node script to run the proxy server and Javascript integration tests.
// Usage: node test <directory with tests>

var childProcess = require('child_process');
var fs = require('fs');
var path = process.cwd();


var build = childProcess.spawnSync('go', [
	'build',
	'-o',
	__dirname + '/out',
	__dirname + '/main/test_server.go'
], {
	stdio: 'inherit'
});

var server = childProcess.spawn(__dirname + '/out', [
	'-port=' + 8000,
	'-logLevel=' + 1
], {
	stdio: 'inherit'
});

testFile = process.argv[2]
i = 0;

// hack to make sure that the server is up before running tests.
setTimeout(function() {

	console.log("Testing: " + testFile)
	var shell = childProcess.spawnSync('mongo', [
		testFile,
		'--port=8000'
	], {
		stdio: 'inherit'
	})
	server.kill();
}, 2000)
