'use strict';

module.exports = {

	// AJAX request for a metric with the rule at index i, the given time granularity,
	// the time of the beginning of the data, and the time of the end of the data.
	// The callback is called on success, and error is called on error.
	getTabularMetric: function(i, timeGranularity, startTime, endTime, callback, error) {
		
		// TODO: Find a better way to reference rules rather than with index. Currently, rules
		// do not have a unique ID, so the easiest way to get the unique identifier is with their
		// position in the array.
		
		$.ajax('/tabular/' + i + '/' + timeGranularity + '/' + startTime.toISOString() + '/' +
			endTime.toISOString(), {
				success: callback,
				error: error
			})
	},

	// AJAX request for a metric with the rule at index i, the given time granularity,
	// the time of the beginning of the data, the time of the end of the data, and the value
	// of the field to return.
	// The callback is called on success, and error is called on error.
	getTabularMetricValue: function(i, timeGranularity, startTime, endTime, value, callback, error) {
		
		// TODO: Find a better way to reference rules rather than with index. Currently, rules
		// do not have a unique ID, so the easiest way to get the unique identifier is with their
		// position in the array.

		$.ajax('/tabular/' + i + '/' + timeGranularity + '/' + startTime.toISOString() + '/' +
			endTime.toISOString() + '/' + value, {
				success: callback,
				error: error
			})
	},

	// AJAX request to retrieve metadata for a rule index and a time granularity. 
	// The callback is called on success, and error is called on error.
	getMetadata: function(i, timeGranularity, callback, error) {
		$.ajax('/metadata/' + i + '/' + timeGranularity, {
				success: callback,
				error: error
			})
	},

	// AJAX request to save a configuration document.
	// The callback is called on success, and error is called on error.
	postConfiguration: function(configJSON, callback, error) {
		$.ajax({
			type: "POST",
			url: '/config',
			data: JSON.stringify(configJSON),
			success: callback,
			error: error
		})
	}

}
