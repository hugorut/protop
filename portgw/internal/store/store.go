package store

import "github.com/hugorut/protop/portgw/internal"

type PortStore interface {
	Store(port internal.Port) error
}
