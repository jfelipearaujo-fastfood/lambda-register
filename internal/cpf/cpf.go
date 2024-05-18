package cpf

import (
	"regexp"
	"strings"
)

var (
	cpfFirstDigitTable  = []int{10, 9, 8, 7, 6, 5, 4, 3, 2}
	cpfSecondDigitTable = []int{11, 10, 9, 8, 7, 6, 5, 4, 3, 2}
)

const (
	CPFFormatPattern string = `([\d]{3})([\d]{3})([\d]{3})([\d]{2})`
)

type CPF string

func NewCPF(s string) CPF {
	return CPF(Clean(s))
}

func (c *CPF) IsValid() bool {
	return ValidateCPF(string(*c))
}

func (c *CPF) String() string {
	str := string(*c)

	expr, err := regexp.Compile(CPFFormatPattern)
	if err != nil {
		return str
	}

	if !c.IsValid() {
		return str
	}

	return expr.ReplaceAllString(str, "$1.$2.$3-$4")
}

func (c *CPF) Mask() string {
	cpf := string(*c)

	return strings.ReplaceAll(cpf, cpf[3:(len(cpf)-2)], strings.Repeat("*", len(cpf)-5))
}
