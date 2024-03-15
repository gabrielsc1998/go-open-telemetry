package domain

import "errors"

type Cep struct {
	Value string
}

func NewCep(value string) (*Cep, error) {
	cep := &Cep{Value: value}
	isValid := cep.validate()
	if !isValid {
		return nil, errors.New("invalid cep")
	}
	return cep, nil
}

func (c *Cep) validate() bool {
	return c.Value != "" && len(c.Value) == 8
}
