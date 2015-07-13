function addLabelsToRules(rules) {
	if (!rules || !rules.length) {
		return rules;
	}

	// add extra fields to the rules that are useful for the frontend
	for (var i = 0; i < rules.length; i++) {
		rules[i].label = rules[i].ValueField;
		rules[i].value = i;
		rules[i].index = i;
	}

	return rules;
}

module.exports = addLabelsToRules;
