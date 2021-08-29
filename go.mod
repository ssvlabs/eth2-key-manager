module github.com/bloxapp/eth2-key-manager

go 1.15

require (
	github.com/bazelbuild/buildtools v0.0.0-20200528175155-f4e8394f069d // indirect
	github.com/confluentinc/confluent-kafka-go v1.4.2 // indirect
	github.com/ethereum/go-ethereum v1.9.25
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/google/uuid v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.14.6 // indirect
	github.com/herumi/bls-eth-go-binary v0.0.0-20210130185500-57372fb27371
	github.com/kr/text v0.2.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prestonvanloon/go-recaptcha v0.0.0-20190217191114-0834cef6e8bd // indirect
	github.com/protolambda/zssz v0.1.5 // indirect
	github.com/prysmaticlabs/eth2-types v0.0.0-20210303084904-c9735a06829d
	github.com/prysmaticlabs/go-ssz v0.0.0-20200612203617-6d5c9aa213ae
	github.com/prysmaticlabs/prysm v1.4.2-0.20210827024218-7757b49f067e
	github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.7.0
	github.com/tyler-smith/go-bip39 v1.1.0
	github.com/wealdtech/go-eth2-util v1.6.3
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	golang.org/x/text v0.3.6
	gopkg.in/confluentinc/confluent-kafka-go.v1 v1.4.2 // indirect
)

replace github.com/ethereum/go-ethereum => github.com/prysmaticlabs/bazel-go-ethereum v0.0.0-20210707101027-e8523651bf6f
