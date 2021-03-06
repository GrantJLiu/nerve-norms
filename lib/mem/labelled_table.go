package mem

import (
	"encoding/json"
	"errors"
	"fmt"
)

type LabelledTable interface {
	XName() string
	YName() string
	XColumnAt(int) float64
	YColumnAt(int) float64
	WasImputedAt(int) bool
	Len() int
	IncludeOutlierScore(int) bool
}

type LabelledTableFromMem func(*Mem) LabelledTable

type LabTab struct {
	section string
	xname   string
	yname   string
	xcol    Column
	ycol    Column
	wasimp  Column

	// TE has more than one table
	tableNum int

	// Set to true if using log interpolation
	logScale bool

	// These fields can be set to run an alternative column import
	altSection    string
	altXname      string
	altYname      string
	altImportFunc func(*LabTab)

	// This can be set to do something with the data before or after imputation.
	// The passed Column is the raw xcol, and the returned Column is the updated xcol.
	preImputeAction  func(RawSection, Column) Column
	postImputeAction func()

	// precision is the float precision
	precision float64
}

func (lt LabTab) XName() string {
	return lt.xname
}

func (lt LabTab) YName() string {
	return lt.yname
}

func (lt LabTab) XColumnAt(idx int) float64 {
	return lt.xcol[idx]
}

func (lt LabTab) YColumnAt(idx int) float64 {
	return lt.ycol[idx]
}

func (lt LabTab) WasImputedAt(idx int) bool {
	return len(lt.wasimp) != 0 && lt.wasimp[idx] > 0.5
}

func (lt LabTab) Len() int {
	return len(lt.ycol)
}

func (lt LabTab) IncludeOutlierScore(idx int) bool {
	// Don't include ones that were imputed
	return !lt.WasImputedAt(idx)
}

func NewLabTab(xname, yname string, xcol, ycol, wasimp Column) LabTab {
	return LabTab{
		xname:  xname,
		yname:  yname,
		xcol:   xcol,
		ycol:   ycol,
		wasimp: wasimp,
	}
}

type emptyLT struct{}

func (lt emptyLT) XName() string                    { return "" }
func (lt emptyLT) YName() string                    { return "" }
func (lt emptyLT) XColumnAt(idx int) float64        { return 0 }
func (lt emptyLT) YColumnAt(idx int) float64        { return 0 }
func (lt emptyLT) WasImputedAt(idx int) bool        { return true }
func (lt emptyLT) Len() int                         { return 0 }
func (lt emptyLT) IncludeOutlierScore(idx int) bool { return false }

// jsonTable is used to restructure LabTab data for json.
type jsonTable struct {
	Columns []string `json:"columns"`
	Data    Table    `json:"data"`
}

func (lt LabTab) MarshalJSON() ([]byte, error) {
	jt := jsonTable{
		Columns: []string{lt.xname, lt.yname},
		Data:    []Column{lt.xcol, lt.ycol},
	}

	if lt.wasimp != nil {
		jt.Columns = append(jt.Columns, "Was Imputed")
		jt.Data = append(jt.Data, lt.wasimp)
	}

	return json.Marshal(&jt)
}

func (lt *LabTab) UnmarshalJSON(value []byte) error {
	var jt jsonTable
	err := json.Unmarshal(value, &jt)
	if err != nil {
		return err
	}

	numCol := len(jt.Columns)
	numDat := len(jt.Data)

	if numCol < 2 || numCol > 3 || numDat < 2 || numDat > 3 {
		if numCol == 0 && numDat == 0 {
			return nil // This section is nothing
		}
		return fmt.Errorf("Incorrect number of LabelledTable columns in JSON %d %d", numCol, numDat)
	}

	lt.xname = jt.Columns[0]
	lt.yname = jt.Columns[1]
	lt.xcol = jt.Data[0]
	lt.ycol = jt.Data[1]

	if numCol == 3 {
		if jt.Columns[2] != "Was Imputed" {
			return errors.New("Incorrect TablelledTable column names in JSON")
		}
		lt.wasimp = jt.Data[2]
	}

	return nil
}

type MissingSection error

func (lt *LabTab) LoadFromMem(mem *rawMem) error {
	sec, err := mem.sectionContainingHeader(lt.section)
	if err != nil && lt.altSection != "" {
		// Sometimes an old format spelled this incorrectly
		sec, err = mem.sectionContainingHeader(lt.altSection)
	}
	if err != nil {
		return MissingSection(errors.New("Could not get LT section " + lt.section + ": " + err.Error()))
	}

	xcol, err := sec.columnContainsName(lt.xname, lt.tableNum)
	if err != nil && lt.altXname != "" {
		// For some reason this column sometimes has the wrong name in older files
		xcol, err = sec.columnContainsName(lt.altXname, lt.tableNum)
	}
	if err != nil {
		return errors.New("Could not get LT " + lt.section + " xcol: " + err.Error())
	}

	lt.ycol, err = sec.columnContainsName(lt.yname, lt.tableNum)
	if err != nil && lt.altYname != "" {
		// Some old formats use this mis-labeled column that must be converted
		lt.ycol, err = sec.columnContainsName(lt.altYname, lt.tableNum)
		if err != nil {
			return errors.New("Could not get LT " + lt.section + " alt ycol: " + err.Error())
		}

		if lt.altImportFunc != nil {
			lt.altImportFunc(lt)
		}
	} else if err != nil {
		return errors.New("Could not get LT " + lt.section + " ycol: " + err.Error())
	}

	if lt.preImputeAction != nil {
		xcol = lt.preImputeAction(sec, xcol)
	}

	lt.wasimp = lt.ycol.ImputeWithValue(xcol, lt.xcol, lt.precision, lt.logScale)

	if len(lt.xcol) != len(lt.ycol) {
		return fmt.Errorf("Mismatching LT "+lt.section+" lengths %d and %d (%v and %v)", len(lt.xcol), len(lt.ycol), lt.xcol, lt.ycol)
	}

	if lt.postImputeAction != nil {
		lt.postImputeAction()
	}

	return nil
}

func (lt *LabTab) LabelledTable(unused string) LabelledTable {
	return lt
}
