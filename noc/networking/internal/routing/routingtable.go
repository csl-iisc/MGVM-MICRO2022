package routing

import "gitlab.com/akita/akita"

// Table is a routing table that can find the next-hop port according to the
// final destination.
type Table interface {
	FindPort(dst akita.Port) akita.Port
	DefineRoute(finalDst, outputPort akita.Port)
	DefineDefaultRoute(outputPort akita.Port)
}

// NewTable creates a new Table.
func NewTable() Table {
	t := &table{}
	t.t = make(map[akita.Port]akita.Port)
	return t
}

type table struct {
	t           map[akita.Port]akita.Port
	defaultPort akita.Port
}

func (t table) FindPort(dst akita.Port) akita.Port {
	out, found := t.t[dst]
	if found {
		return out
	}
	return t.defaultPort
}

func (t *table) DefineRoute(finalDst, outputPort akita.Port) {
	t.t[finalDst] = outputPort
}

func (t *table) DefineDefaultRoute(outputPort akita.Port) {
	t.defaultPort = outputPort
}
