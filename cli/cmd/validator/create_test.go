package validator_test

import (
	"bytes"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/eth2-key-manager/cli/cmd"
	"github.com/bloxapp/eth2-key-manager/cli/cmd/validator"
	"github.com/bloxapp/eth2-key-manager/cli/util/printer"
)

func TestValidatorCreate(t *testing.T) {
	walletPK := "15c20889f519082fccd95b385bb304bb29bf531a58afe2a67c89ebf802a23d1b"
	walletAddr := "7015514B3da332d95EE1B94d32ADce4cAa0bAa28"

	t.Run("successfully create one validator for one seed", func(t *testing.T) {
		var getBalanceCalled int
		var getTransactionCountCalled int
		var gasPriceCalled int
		var getCodeCalled int
		var estimateGasCalled int
		var sendTransactionCalled int
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			body := make(map[string]interface{})
			err := json.NewDecoder(r.Body).Decode(&body)
			require.NoError(t, err)
			defer r.Body.Close()

			switch body["method"] {
			case "eth_getBalance":
				getBalanceCalled++
				require.Equal(t, strings.ToLower("0x"+walletAddr), body["params"].([]interface{})[0].(string))
				balance := hexutil.Big(*new(big.Int).Mul(big.NewInt(32*1e9), big.NewInt(1e9)))
				resp, err := json.Marshal(map[string]interface{}{
					"result": balance.String(),
				})
				require.NoError(t, err)
				w.Write(resp)
				break
			case "eth_getTransactionCount":
				getTransactionCountCalled++
				require.NotEmpty(t, body["params"].([]interface{})[0].(string))
				resp, err := json.Marshal(map[string]interface{}{
					"result": hexutil.Uint64(5206).String(),
				})
				require.NoError(t, err)
				w.Write(resp)
				break
			case "eth_gasPrice":
				gasPriceCalled++
				balance := hexutil.Big(*big.NewInt(1e9))
				resp, err := json.Marshal(map[string]interface{}{
					"result": balance.String(),
				})
				require.NoError(t, err)
				w.Write(resp)
				break
			case "eth_getCode":
				getCodeCalled++
				require.NotEmpty(t, body["params"].([]interface{})[0].(string))
				code := hexutil.Bytes("5206")
				resp, err := json.Marshal(map[string]interface{}{
					"result": code.String(),
				})
				require.NoError(t, err)
				w.Write(resp)
				break
			case "eth_estimateGas":
				estimateGasCalled++
				resp, err := json.Marshal(map[string]interface{}{
					"result": hexutil.Uint64(1).String(),
				})
				require.NoError(t, err)
				w.Write(resp)
				break
			case "eth_sendRawTransaction":
				sendTransactionCalled++
				resp, err := json.Marshal(map[string]interface{}{
					"result": "success",
				})
				require.NoError(t, err)
				w.Write(resp)
				break
			}
		}))
		defer srv.Close()

		var resultOut bytes.Buffer
		var output bytes.Buffer
		cmd.ResultPrinter = printer.New(&output)
		validator.ResultFactory = func(name string) (io.Writer, func(), error) {
			return &resultOut, func() {}, nil
		}
		cmd.RootCmd.SetArgs([]string{
			"validator",
			"create",
			"--wallet-private-key", walletPK,
			"--wallet-addr", walletAddr,
			"--validators-per-seed", "1",
			"--seeds-count", "1",
			"--web3-addr", srv.URL,
		})
		err := cmd.RootCmd.Execute()
		require.NoError(t, err)
		require.Equal(t, 1, getBalanceCalled)
		require.Equal(t, 1, getTransactionCountCalled)
		require.Equal(t, 1, gasPriceCalled)
		require.Equal(t, 1, getCodeCalled)
		require.Equal(t, 1, estimateGasCalled)
		require.Equal(t, 1, sendTransactionCalled)
	})
}
