module github.com/bloxapp/eth2-key-manager

go 1.13

require (
	github.com/ethereum/go-ethereum v1.9.20
	github.com/google/uuid v1.1.1
	github.com/pkg/errors v0.9.1
	github.com/prysmaticlabs/ethereumapis v0.0.0-20200827165051-58ccb36e36b9
	github.com/prysmaticlabs/go-ssz v0.0.0-20200612203617-6d5c9aa213ae
	github.com/prysmaticlabs/prysm v1.0.0-alpha.24
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.6.1
	github.com/tyler-smith/go-bip39 v1.0.2
	github.com/wealdtech/eth2-signer-api v1.3.0
	github.com/wealdtech/go-eth2-types/v2 v2.5.0
	github.com/wealdtech/go-eth2-util v1.5.0
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de
	golang.org/x/text v0.3.3
)

replace github.com/ethereum/go-ethereum => github.com/prysmaticlabs/bazel-go-ethereum v0.0.0-20200530091827-df74fa9e9621

replace github.com/herumi/bls-eth-go-binary => github.com/herumi/bls-eth-go-binary v0.0.0-20200605082007-3a76b4c6c599
