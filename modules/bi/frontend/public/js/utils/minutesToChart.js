var _ = require('lodash');
var moment = require('moment');

var Controller = require('../ajax/controller')

var minutesToChart = function(input, startTime, number) {

	// test for minute first
	if (!input || !input.length) {
		return [];
		var d = [];
		for (var i = 0; i < number; i++) {
			d[i] = 0;
		}
		return d;
	}

	_.sortBy(input, function(n) {
		return n.start;
	});

	var dataArray = [];

	var start = input[0].start.clone();
	var prePadding = input[0].start.clone();
	while (prePadding.isAfter(startTime)) {
		prePadding.hour(prePadding.hour() - 1);
		for (var i = 0; i < number; i++) {
			dataArray.push(0);
		}
	}

	// get the difference in minutes
	var diffInMinutes = startTime.diff(prePadding, 'minutes')

	// convert the input into a flat array (with 0 as all the blanks)

	var expectedTime = start;

	for (var i = 0; i < input.length; i++) {
		var currentTime = input[i].start;
		for (; currentTime.isBefore(expectedTime); currentTime.minutes(currentTime.minutes() + 1)) {
			dataArray.push(0)
		}
		for (var j = 0; j < 60; j++) {
			dataArray.push(input[i].minute[j] || 0)
		}
		expectedTime.hours(expectedTime.hours() + 1)
	}

	// figure out the start time
	// return a slice of the array

	console.log(diffInMinutes);
	var slicedArray = dataArray.slice(diffInMinutes + 2, diffInMinutes + 2 + number);
	for (var i = slicedArray.length; i < number; i++) {
		slicedArray.unshift(0);
	}
	return slicedArray;
}

function getMinutes(rules, index, callback) {
	var endTime = moment();
	var startTime = moment();
	var rule = rules[index];

	// TODO: eventually be able to expand the range
	var range = 60;
	startTime.minutes(endTime.minutes() - range);

	Controller.getMetric(index, "m", rule.ValueField, startTime, endTime,
		function(data) {
			if (!data) {
				console.log(data);
				var d = [];
				for (var i = 0; i < range; i++) {
					d[i] = 0;
				}
				callback(d);
				return;
			}
			for (var i = 0; i < data.length; i++) {
				data[i].start = moment(data[i].start)
			}

			dataArray = minutesToChart(data, startTime, 60);
			callback(dataArray);
		});
}

module.exports = getMinutes;
