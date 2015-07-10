// Should try to develop here...
var _ = require('lodash');
var moment = require('moment');

var Controller = require('../ajax/controller')

var timeFormat = "YYYY-MM-DD HH:mm:ss";

// helper function to get the proper time scale for moment.JS given a time granularity.
var getProperGranularity = function(granularity) {
	switch (granularity) {
		case "M":
			return "month"
			break;
		case "D":
			return "day"
			break;
		case "h":
			return "hour"
			break;
		case "m":
			return "minute"
			break;
		case "s":
			return "second"
			break;
		default:
			return ""
	}
}

// helper function to get the next higher time scale for moment.JS given a time granularity.
var getHigherGranularity = function(granularity) {
	switch (granularity) {
		case "M":
			return "year"
			break;
		case "D":
			return "month"
			break;
		case "h":
			return "day"
			break;
		case "m":
			return "hour"
			break;
		case "s":
			return "minute"
			break;
		default:
			return ""
	}
}

// helper function to get the number of ticks of this granularity happens in the next higher time
// granularity.
var getNumInNextGranularity = function(granularity, time) {
	switch (granularity) {
		case "M":
			return 12
			break;
		case "D":
			switch (time.month) {
				case 4:
				case 6:
				case 9:
				case 11:
					return 30;
					break;
				case 2:
					if (time.isLeapYear()) {
						return 29;
					}
					return 28;
					break;
				default:
					return 31;
			}
			break;
		case "h":
			return 24
			break;
		case "m":
			return 60
			break;
		case "s":
			return 60
			break;
		default:
			return 0
	}
}

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
var metricsToChart = function(input, granularity, startTime, range) {

	// create the array for all the timestamps
	var slicedTime = [];
	var cTime = startTime.clone().startOf(getProperGranularity(granularity));

	for (var i = 0; i < range; i++) {
		slicedTime.push(cTime.format(timeFormat));
		cTime.add(1, getProperGranularity(granularity));
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

	// make sure that the inputs are in sorted order, by their start times.
	_.sortBy(input, function(n) {
		return n.start;
	});

	for (var i = 0; i < input.length; i++) {

		var ticks = getNumInNextGranularity(granularity);

		for (var j = 0; j < ticks; j++) {
			var val = input[i][getProperGranularity(granularity)][j] || 0;
			var currentTime = input[i].start.clone().add(j, getProperGranularity(granularity));
			var index = currentTime.diff(startTime, getProperGranularity(granularity));

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

function getMetrics(rules, granularity, index, callback) {
	var endTime = moment();
	var startTime = moment();
	var rule = rules[index];

	// TODO: eventually be able to expand the range
	var range = 60;
	startTime.subtract(range, getProperGranularity(granularity));

	Controller.getMetric(index, granularity, startTime, endTime,
		function(data) {
			if (!data) {
				dataObj = metricsToChart(data, granularity, startTime, range);
				callback(dataObj);
				return;
			}
			for (var i = 0; i < data.length; i++) {
				data[i].start = moment(data[i].start)
			}

			dataObj = metricsToChart(data, granularity, startTime, range);
			callback(dataObj);
		});

}

module.exports = getMetrics;
