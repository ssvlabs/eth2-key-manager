package core

import (
	"sync"

	"github.com/herumi/bls-eth-go-binary/bls"
	e2types "github.com/wealdtech/go-eth2-types/v2"
)

// initBLS initializes BLS ONLY ONCE!
var initBLSOnce sync.Once

func InitBLS() error {
	var err error
	var wg sync.WaitGroup
	initBLSOnce.Do(func() {
		wg.Add(1)

		if err = bls.Init(bls.BLS12_381); err != nil {
			return
		}
		err = bls.SetETHmode(bls.EthModeDraft07)

		err = e2types.InitBLS()
		wg.Done()
	})
	wg.Wait()
	return err
}
