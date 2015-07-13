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
	}

}
