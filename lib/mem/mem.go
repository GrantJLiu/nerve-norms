package mem

import (
	"errors"
	"io"
	"strings"
)

type Mem struct {
	Header
	Sections []Section
	ExcitabilityVariables
}

func Import(data io.Reader) (Mem, error) {
	reader := NewReader(data)
	mem := Mem{}
	mem.ExcitabilityVariables.Values = make(map[string]float64)

	err := mem.Header.Parse(reader)
	if err != nil {
		return mem, err
	}

	for err == nil {
		err = mem.importSection(reader)
	}

	if err != io.EOF && err != nil {
		return mem, errors.New("Error encountered before EOF: " + err.Error())
	}

	return mem, nil
}

func (mem *Mem) importSection(reader *Reader) error {
	str, err := reader.skipNewlines()
	if err != nil {
		return err
	}

	if len(str) < 2 || str[0] != ' ' {
		return errors.New("Could not parse invalid section: '" + str + "'")
	}

	if strings.Contains(str, "DERIVED EXCITABILITY VARIABLES") {
		return mem.ExcitabilityVariables.Parse(reader)
	}

	sec := Section{Header: strings.TrimSpace(str)}
	err = sec.parse(reader)
	if err != nil {
		return err
	}
	mem.Sections = append(mem.Sections, sec)

	return nil
}

func (mem Mem) String() string {
	str := "Mem{\n"
	str += "\t" + mem.Header.String() + ",\n"
	for _, sec := range mem.Sections {
		str += "\t" + sec.String() + ",\n"
	}
	str += "\t" + mem.ExcitabilityVariables.String() + ",\n"
	str += "}"
	return str
}
