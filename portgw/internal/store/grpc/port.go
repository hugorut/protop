package grpc

import (
	"github.com/sirupsen/logrus"

	"github.com/hugorut/protop/portgw/internal"
)

// PortStore provides implements the port.PortStore interface.
// It used a gprc endpoint to persist ports to an external service.
type PortStore struct {
	Logger logrus.FieldLogger
}

// Store wraps a grpc endpoint on the PortDomain service adhearing
// to the PortStore base interface.
//
// @TODO this is currently not implemented because of time constraints, simply logs to stdout.
func (p PortStore) Store(port internal.Port) error {
	p.Logger.Infof("storing port %s", port.Name)
	return nil
}
