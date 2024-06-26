package domain

import "fmt"

type House struct {
	ID       uint64
	Address  string
	PhotoURL string
}

func (h *House) GetFileName() string {
	return fmt.Sprintf("%d-%s", h.ID, h.Address)
}
