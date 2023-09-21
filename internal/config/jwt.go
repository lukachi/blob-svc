package config

import (
	"sync"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
)

type JWTer interface {
	Secret() string
}

func NewJWTSecret(getter kv.Getter) JWTer {
	return &JWT{
		getter: getter,
	}
}

type JWT struct {
	getter kv.Getter
	once   sync.Once
	value  string
	err    error
}

func (j *JWT) Secret() string {
	j.once.Do(func() {
		var config struct {
			Secret string `fig:"secret,required"`
		}
		err := figure.
			Out(&config).
			From(kv.MustGetStringMap(j.getter, "jwt")).
			Please()

		if err != nil {
			j.err = errors.Wrap(err, "failed to figure out listener")
			return
		}

		j.value, j.err = config.Secret, nil
	})

	if j.err != nil {
		panic(j.err)
	}

	return j.value
}
