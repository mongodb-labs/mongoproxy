function addLabelsToRules(rules) {
	if (!rules || !rules.length) {
		return rules;
	}
	for (var i = 0; i < rules.length; i++) {
		rules[i].label = rules[i].ValueField;
		rules[i].value = i;
	}

	return rules;
}

module.exports = addLabelsToRules;
