module github.com/bloxapp/eth2-key-manager

go 1.15

require (
	github.com/ethereum/go-ethereum v1.9.24
	github.com/google/uuid v1.1.2
	github.com/herumi/bls-eth-go-binary v0.0.0-20201104034342-d782bdf735de
	github.com/pkg/errors v0.9.1
	github.com/prysmaticlabs/ethereumapis v0.0.0-20201003171600-a72e5f77d233
	github.com/prysmaticlabs/go-ssz v0.0.0-20200612203617-6d5c9aa213ae
	github.com/prysmaticlabs/prysm v1.0.0-beta.2
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.6.1
	github.com/tyler-smith/go-bip39 v1.0.2
	github.com/wealdtech/go-eth2-util v1.6.2
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897
	golang.org/x/text v0.3.4
)

replace github.com/ethereum/go-ethereum => github.com/prysmaticlabs/bazel-go-ethereum v0.0.0-20201113091623-013fd65b3791
