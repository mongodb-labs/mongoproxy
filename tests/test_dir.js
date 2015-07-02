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

testFiles = fs.readdirSync(process.argv[2])
i = 0;

// hack to make sure that the server is up before running tests.
setTimeout(function() {

	var id = setInterval(function() {
		if (!testFiles[i]) {
			clearInterval(id)
			return;
		}
		var ext = testFiles[i].split('.').pop();
		if (ext != 'js') {
			i++;
			return;
		}

		console.log("Testing: " + testFiles[i])
		var shell = childProcess.spawnSync('mongo', [
			process.argv[2] + '/' + testFiles[i],
			'--port=8000'
		], {
			stdio: 'inherit'
		})
		i++;

		if (i >= testFiles.length) {
			clearInterval(id);
			server.kill();
		}
	}, 100)

}, 2000)
