package horizon

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/horizon-connector"
)

type BlobsQ struct {
	Log     *logan.Entry
	Horizon *horizon.Connector
}

func NewBlobsQ(log *logan.Entry, horizon *horizon.Connector) *BlobsQ {
	return &BlobsQ{
		Log:     log,
		Horizon: horizon,
	}
}
