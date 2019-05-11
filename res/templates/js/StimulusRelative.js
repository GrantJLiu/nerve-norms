class StimulusRelative extends Chart {
	constructor(plots) {
		super([0, 200], [0, 100])
		this.participant = this.calculateData(plots.sr.data)
		this.xName = 'x'
		this.yName = 'y'
		this.xSDName = 'SD'
		this.ySDName = undefined
		this.xMeanName = 'mean'
		this.yMeanName = undefined
	}

	calculateData(data) {
		const stimFor50PercentMax = data[24].valueX // Could also be extracted from excitability variables
		const meanStimFor50PercentMax = data[24].meanX

		return data.map((d, i) => {
			return {
				'y': (i + 1) * 2,
				// Normalize each element
				'x': d.valueX / stimFor50PercentMax * 100,
				'SD': d.SDX / d.meanX * 100,
				'mean': d.meanX / meanStimFor50PercentMax * 100,
			}
		})
	}

	get name() { return "Relative Stimulus Response" }
	get xLabel() { return "Stimulus (% Mean Threshold)" }
	get yLabel() { return "Peak Response (% Max)" }

	updatePlots(plots) {
		this.participant = this.calculateData(plots.sr.data)
		this.animateXYLine(this.participant, "srel")
	}

	updateNorms(norms) {
		this.animateNorms(this.participant, "srel")
	}

	drawLines(svg) {
		this.createXYLine(this.participant, "srel")
		this.createNorms(this.participant, "srel")
		this.animateXYLine(this.participant, "srel")
		this.animateNorms(this.participant, "srel")
	}
}
