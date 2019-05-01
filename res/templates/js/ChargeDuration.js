class ChargeDuration extends Chart {
	constructor(data) {
		super([0, 1], [0, 10])
		this.data = data
		this.xName = 'stimWidth'
	}

	get name() { return "Charge Duration" }
	get xLabel() { return "Stimulus Width (ms)" }
	get yLabel() { return "Threshold Change (mA•ms)" }

	drawLines(svg) {
		this.animateXYLineWithMean(this.data)
	}
}
