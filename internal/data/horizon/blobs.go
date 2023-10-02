package horizon

import (
	"context"
	"encoding/json"
	"github.com/lukachi/blob-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/keypair"
)

type BlobsQ interface {
	WriteBlob()
	Getblob(xdr.Uint64)
}

type Blobs struct {
	BlobsQ

	Log       *logan.Entry
	connector *horizon.Connector
	builder   *xdrbuild.Builder
	source    keypair.Address
	signer    keypair.Full
}

func NewBlobsQ(log *logan.Entry, connector *horizon.Connector, builder *xdrbuild.Builder, source keypair.Address, signer keypair.Full) *Blobs {
	return &Blobs{
		Log:       log,
		connector: connector,
		builder:   builder,
		source:    source,
		signer:    signer,
	}
}

func (b *Blobs) WriteBlob(blob *data.Blob) (xdr.Uint64, error) {
	envelope, err := b.builder.Transaction(b.source).Op(xdrbuild.CreateData{
		Type:  1,
		Value: json.RawMessage(blob.Value),
	}).Sign(b.signer).Marshal()

	if err != nil {
		return 0, errors.Wrap(err, "failed to Marshal tx")
	}

	result := b.connector.Submitter().Submit(context.TODO(), envelope)

	if result.Err != nil {
		if len(result.OpCodes) == 1 {
			switch result.OpCodes[0] {
			case "op_not_ready":
				return 0, nil
			}
		}

		fields := logan.F{
			"result error": result.Err.Error(),
		}

		if result.Err == horizon.ErrSubmitRejected {
			fields["tx code"] = result.TXCode
			fields["op codes"] = result.OpCodes
			fields["result xdr"] = result.ResultXDR
		}

		return 0, errors.Wrap(result.Err, "failed to submit tx", fields)
	}

	txResult := xdr.TransactionResult{}

	err = txResult.Scan(result.ResultXDR)

	if err != nil {
		return 0, errors.Wrap(err, "failed to get transaction result")
	}

	opResults, isOk := txResult.Result.GetResults()

	if !isOk {
		return 0, errors.New("failed to get transaction result")
	}

	createDataResult, isOk := opResults[0].Tr.CreateDataResult.GetSuccess()

	if !isOk {
		return 0, errors.New("failed to create data")
	}

	return createDataResult.DataId, err
}

func (b *Blobs) GetBlob(id xdr.Uint64) (*data.Blob, error) {
	blobBytes, err := b.connector.Client().Get("/blobs/" + string(id))

	if err != nil {
		return nil, errors.Wrap(err, "failed to get blob")
	}

	if blobBytes == nil {
		return nil, nil
	}

	var blob struct {
		Data data.Blob `json:"data"`
	}

	err = json.Unmarshal(blobBytes, &blob)

	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal blob")
	}

	return &blob.Data, nil
}
