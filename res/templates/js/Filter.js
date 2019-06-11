class Filter {
	static get url() { return "https://us-central1-nervenorms.cloudfunctions.net/" }

	constructor(callback) {
		this.callback = callback
		this.queryString = Filter._queryString

		document.querySelector("form").addEventListener("submit", (event) => {
			ExVars.clearScores()
			this.update().fetchNorms().fetchOutliers()
			event.preventDefault()
		})
	}

	update() {
		this.queryString = Filter._queryString
		return this
	}

	setParticipant(name) {
		this.name = name
		return this
	}

	static get _queryString() {
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
					Filter._setAgeOptions(filter, entry[1])
					break
			}
		}

		return "?" +
			Object.keys(filter).map(function(key) {
				return encodeURIComponent(key) + "=" +
					encodeURIComponent(filter[key]);
			}).join("&");
	}

	static _setAgeOptions(filter, opts) {
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

	fetchOutliers() {
		fetch(Filter.url + "outliers" + this.queryString + "&name=" + this.name)
			.then(response => {
				return response.json()
			})
			.then(scores => {
				this.scores = scores
				ExVars.updateScores(this.name, scores)
			})
		return this
	}

	fetchNorms() {
		fetch(Filter.url + "norms" + this.queryString)
			.then(response => {
				return response.json()
			})
			.then(norms => {
				this.callback(norms)
			})
		return this
	}
}
