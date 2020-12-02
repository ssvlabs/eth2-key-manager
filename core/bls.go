package core

import (
	"sync"

	"github.com/herumi/bls-eth-go-binary/bls"
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
		wg.Done()
	})
	wg.Wait()
	return err
}
