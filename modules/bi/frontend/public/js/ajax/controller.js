


function ISODateString(d){
	if (d.toDate) {
		d = d.toDate();
	}
 	function pad(n){
 		return n<10 ? '0'+n : n
 	}
 	return d.getUTCFullYear()+'-'
      + pad(d.getUTCMonth()+1)+'-'
      + pad(d.getUTCDate())+'T'
      + pad(d.getUTCHours())+':'
      + pad(d.getUTCMinutes())+':'
      + pad(d.getUTCSeconds())+'Z';
  	}

module.exports = {

	getMetric: function(i, timeGranularity, valueField, startTime, endTime, cb) {
		$.ajax('/data/'+i+'/'+timeGranularity+'/'+valueField+'/' + ISODateString(startTime) + '/' + ISODateString(endTime), {
			success: cb,
			error: function() {
				// bad things happen here
			}
		})
	}

}
