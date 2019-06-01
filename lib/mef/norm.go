package mef

import "gogs.bellstone.ca/james/jitter/lib/mem"

type Norm struct {
	CDNorm    NormTable            `json:"cd"`
	RCNorm    NormTable            `json:"rc"`
	SRXNorm   NormTable            `json:"srx"`
	SRYNorm   NormTable            `json:"sry"`
	SRelXNorm NormTable            `json:"srelx"`
	IVNorm    NormTable            `json:"iv"`
	TENorm    map[string]NormTable `json:"te"`
}

func (mef *Mef) Norm() Norm {
	mef = mef.FilteredMef()

	norm := Norm{
		CDNorm:    NewNormTable(mem.CDDuration, mef, "CD", ""),
		IVNorm:    NewNormTable(mem.IVCurrent, mef, "IV", ""),
		RCNorm:    NewNormTable(mem.RCInterval, mef, "RC", ""),
		SRXNorm:   NewNormTable(nil, mef, "SR", "calculatedX"),
		SRYNorm:   NewNormTable(nil, mef, "SR", "calculatedY"),
		SRelXNorm: NewNormTable(nil, mef, "SR", "relative"),
		TENorm:    map[string]NormTable{},
	}

	for _, name := range []string{"h40", "h20", "d40", "d20"} {
		norm.TENorm[name] = NewNormTable(mem.TEDelay(name), mef, "TE", name)
	}

	return norm
}
