class RecoveryCycle extends Chart {
	constructor(data) {
		super([1, 200], [-50, 110], Chart.scaleType.LOG)
		this.data = data
	}

	get name() { return "Recovery Cycle" }
	get xLabel() { return "Threshold Change (%)" }
	get yLabel() { return "Interstimulus Interval (ms)" }

	drawLines(svg) {
		this.animateXYLineWithMean(this.data)
		this.drawHorizontalLine(this.linesLayer, 0)
	}
}