var g = require('./granularities')

/* 
Given some input metric documents (JSON array), a time granularity (string), a startTime
(moment.JS object), and an integer range, return the data in a continuous-time
array, where range is the number of ticks of that time granularity.

Input metric documents have the following format:
{
	valueField: string,
	start: Moment.JS object - the time that this input started recording data,
	total: number,
	<month|day|hour|minute|second>: {
		1: number
		2: number
		...
	} - an object with data at the time granularity mentioned in the field.
}

The data from the various input documents are unrolled into a flat array, with the 
first value beginning at startTime, and the last value at startTime + range ticks
of the time granularity.

The data is then returned as a JSON object in the following format:
{
	data: [number] - the flattened array of the data,
	time: [string] - formatted timestamps corresponding to the data
}
*/
var dataToTabular = function(input, granularity, startTime, range) {

	// create the array for all the timestamps
	var slicedTime = [];
	var cTime = startTime.clone().startOf(g.getProperGranularity(granularity));

	for (var i = 0; i < range; i++) {
		slicedTime.push(cTime.format(timeFormat));
		cTime.add(1, g.getProperGranularity(granularity));
	}

	// initialize the data array
	var dataArray = [];
	for (var i = 0; i < range; i++) {
		dataArray[i] = 0;

	}

	// we have nothing. Return them as is.
	if (!input || !input.length) {

		return {
			data: dataArray,
			time: slicedTime
		}
	}

	for (var i = 0; i < input.length; i++) {

		var ticks = g.getNumInNextGranularity(granularity);

		for (var j = 0; j < ticks; j++) {
			var val = input[i][g.getProperGranularity(granularity)][j] || 0;
			var currentTime = input[i].start.clone().add(j, g.getProperGranularity(granularity));
			var index = currentTime.diff(startTime, g.getProperGranularity(granularity));

			if (index >= 0 && index < range) {
				dataArray[index] = val;
			}

		}

	}

	return {
		data: dataArray,
		time: slicedTime
	}

}
