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
	regources "gitlab.com/tokend/regources/generated"
)

type BlobsQ struct {
	Log       *logan.Entry
	connector *horizon.Connector
	masterKp  *keypair.Full
}

func NewBlobsQ(log *logan.Entry, connector *horizon.Connector, kp *keypair.Full) data.HorizonBlobsQ {
	return &BlobsQ{
		Log:       log,
		connector: connector,
		masterKp:  kp,
	}
}

func (b *BlobsQ) New() data.HorizonBlobsQ {
	return NewBlobsQ(b.Log, b.connector, b.masterKp)
}

type BlobRequest struct {
	data.Blob
}

func (b BlobRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Blob)
}

func (b *BlobsQ) WriteBlob(blob *data.Blob) (xdr.Uint64, error) {
	builder, err := b.connector.TXBuilder()

	if err != nil {
		return 0, errors.Wrap(err, "failed to get builder")
	}

	if err != nil {
		return 0, errors.Wrap(err, "failed to marshal blob")
	}

	envelope, err := builder.Transaction(*b.masterKp).Op(xdrbuild.CreateData{
		Type:  1,
		Value: BlobRequest{*blob},
	}).Sign(*b.masterKp).Marshal()

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

func (b *BlobsQ) GetBlob(id string) (*data.Blob, error) {
	blobBytes, err := b.connector.Client().Get("/v3/data/" + id)

	if err != nil {
		return nil, errors.Wrap(err, "failed to get blob")
	}

	if blobBytes == nil {
		return nil, nil
	}

	var v3Data struct {
		Data regources.Data `json:"data"`
	}

	err = json.Unmarshal(blobBytes, &v3Data)

	var blob = data.Blob{}

	if err := json.Unmarshal(v3Data.Data.Attributes.Value, &blob); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal blob")
	}

	return &blob, nil
}
