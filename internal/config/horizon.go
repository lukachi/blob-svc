package config

import (
	"net/url"
	"os"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/tokend/horizon-connector"
	"gitlab.com/tokend/keypair"
)

func (c *config) Horizon() *horizon.Connector {
	if c.horizon != nil {
		return c.horizon
	}

	var config struct {
		URL    url.URL `fig:"url,required"`
		Signer keypair.Full
	}

	err := figure.
		Out(&config).
		From(kv.MustGetStringMap(c.getter, "horizon")).
		Please()
	if err != nil {
		panic(errors.Wrap(err, "failed to figure out horizon"))
	}

	masterSeed := os.Getenv("MASTER_SEED")
	if len(masterSeed) == 0 {
		panic("the MASTER_SEED enviroment variable does not exist")
	}

	kp, err := keypair.ParseSeed(masterSeed)
	if err != nil {
		panic(err)
	}

	config.Signer = kp

	c.horizon = horizon.NewConnector(&config.URL).WithSigner(config.Signer)

	return c.horizon
}
