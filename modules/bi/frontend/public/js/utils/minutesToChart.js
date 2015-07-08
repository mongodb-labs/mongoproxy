var _ = require('lodash');

module.exports = function(input, startTime, number) {
	
	// test for minute first
	if (!input || !input.length) {
		return [];
	}

	_.sortBy(input, function(n) {
  		return n.start;
	});

	var start = input[0].start;

	// get the difference in minutes
	var diffInMinutes = startTime.diff(start, 'minutes')

	// convert the input into a flat array (with 0 as all the blanks)
	var dataArray = [];
	var expectedTime = start;

	for (var i = 0; i < input.length; i++){
		var currentTime = input[i].start;
		for (; currentTime < expectedTime; currentTime.minutes(currentTime.minutes() + 1)) {
			dataArray.push(0)
		}
		for (var j = 0; j < 60; j ++) {
			dataArray.push(input[i].minute[j] || 0)
		}
		expectedTime.hours(expectedTime.hours() + 1)
	}

	// figure out the start time
	// return a slice of the array

	return dataArray.slice(diffInMinutes, diffInMinutes + number);
}
