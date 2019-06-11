class ChartFactory {
	constructor(participants) {
		this.dataManager = new DataManager(participants, () => { return this.plots })

		this.plots = {
			"recoveryCycle": null,
			"thresholdElectrotonus": null,
			"chargeDuration": null,
			"thresholdIV": null,
			"stimulusResponse": null,
			"stimulusResponseRelative": null,
		}

		Object.keys(this.plots).forEach(key => {
			this.plots[key] = this.build(key)
			this.plots[key].draw(d3.select("#" + key + " svg"), true)
			this.plots[key].setDelayTime(Chart.fastDelay) // After initial setup, remove the delay
		})

		// Now set all excitability variables
		ExVars.updateValues(this.dataManager.data)
	}

	/**
	 * Detect if the current active responsive breakpoint in Bootstrap is not XS
	 * Modified from code by farside {@link https://stackoverflow.com/users/4354249/farside}
	 */
	static bootstrapSizeIsXS() {
		const $el = $("<div>")
		$el.appendTo($("body"))
		$el.addClass("d-sm-none")
		const isBiggerThanXS = $el.is(":hidden")
		$el.remove()
		return !isBiggerThanXS
	};

	drawModal(typeStr) {
		if (ChartFactory.bootstrapSizeIsXS()) {
			// Screen is too small to display this properly
			return
		}
		const chart = this.build(typeStr)
		chart.setDelayTime(Chart.fastDelay).setTransitionTime(Chart.fastTransition)

		document.getElementById('modal-title').innerHTML = chart.name
		d3.selectAll("#modal svg > *").remove()

		chart.draw(d3.select('#modal svg'))
		$('#modal').modal('toggle')
	}

	build(typeStr) {
		switch (typeStr) {
			case "recoveryCycle":
				return new RecoveryCycle(this.dataManager.data, this.dataManager.norms)
			case "thresholdElectrotonus":
				return new ThresholdElectrotonus(this.dataManager.data, this.dataManager.norms)
			case "chargeDuration":
				return new ChargeDuration(this.dataManager.data, this.dataManager.norms)
			case "thresholdIV":
				return new ThresholdIV(this.dataManager.data, this.dataManager.norms)
			case "stimulusResponse":
				return new StimulusResponse(this.dataManager.data, this.dataManager.norms)
			case "stimulusResponseRelative":
				return new StimulusRelative(this.dataManager.data, this.dataManager.norms)
		}
	}
}
