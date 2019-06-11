class Filter {
	static get url() { return "https://us-central1-nervenorms.cloudfunctions.net/" }

	static asQueryString() {
		var data = new FormData(document.querySelector("form"))
		const filter = { "species": "human", "nerve": "median" }
		for (const entry of data) {
			switch (entry[0]) {
				case "sex-options":
					filter.sex = entry[1]
					break
				case "country-options":
					filter.country = entry[1]
					break
				case "age-options":
					Filter.setAgeOptions(filter, entry[1])
					break
			}
		}

		return "?" +
			Object.keys(filter).map(function(key) {
				return encodeURIComponent(key) + "=" +
					encodeURIComponent(filter[key]);
			}).join("&");
	}

	static setAgeOptions(filter, opts) {
		switch (opts) {
			case "any":
				filter.minAge = 0
				filter.maxAge = 200
				break
			case "-30":
				filter.minAge = 0
				filter.maxAge = 30
				break
			case "31-40":
				filter.minAge = 31
				filter.maxAge = 40
				break
			case "41-50":
				filter.minAge = 41
				filter.maxAge = 50
				break
			case "51-60":
				filter.minAge = 51
				filter.maxAge = 60
				break
			case "61-":
				filter.minAge = 61
				filter.maxAge = 200
				break
		}
	}

	static updateOutliers(name, queryString) {
		if (queryString === undefined) {
			queryString = Filter.asQueryString()
		}

		fetch(Filter.url + "outliers" + queryString + "&name=" + name)
			.then(response => {
				return response.json()
			})
			.then(scores => {
				this.scores = scores
				ExVars.updateScores(name, scores)
			})
	}

	static updateNorms(queryString, callback) {
		if (queryString === undefined) {
			queryString = Filter.asQueryString()
		}

		fetch(Filter.url + "norms" + queryString)
			.then(response => {
				return response.json()
			})
			.then(norms => {
				callback(norms)
			})
	}
}
