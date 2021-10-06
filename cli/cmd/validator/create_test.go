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

func blkByNumberRes() map[string]interface{} {
	return map[string]interface{}{
		"result": map[string]interface{}{
			"difficulty":       "0x4ea3f27bc",
			"extraData":        "0x476574682f4c5649562f76312e302e302f6c696e75782f676f312e342e32",
			"gasLimit":         "0x1388",
			"gasUsed":          "0x0",
			"hash":             "0xdc0818cf78f21a8e70579cb46a43643f78291264dda342ae31049421c82d21ae",
			"logsBloom":        "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			"miner":            "0xbb7b8287f3f0a933474a79eae42cbca977791171",
			"mixHash":          "0x4fffe9ae21f1c9e15207b1f472d5bbdd68c9595d461666602f2be20daf5e7843",
			"nonce":            "0x689056015818adbe",
			"number":           "0x1b4",
			"parentHash":       "0xe99e022112df268087ea7eafaf4790497fd21dbeeb6bd7a1721df161a6657a54",
			"receiptsRoot":     "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
			"sha3Uncles":       "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
			"size":             "0x220",
			"stateRoot":        "0xddc8b0234c2e0cad087c8b389aa7ef01f7d79b2570bccb77ce48648aa61c904d",
			"timestamp":        "0x55ba467c",
			"totalDifficulty":  "0x78ed983323d",
			"transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		},
	}
}

func TestValidatorCreate(t *testing.T) {
	walletPK := "15c20889f519082fccd95b385bb304bb29bf531a58afe2a67c89ebf802a23d1b"
	walletAddr := "7015514B3da332d95EE1B94d32ADce4cAa0bAa28"

	t.Run("successfully create one validator for one seed (prater)", func(t *testing.T) {
		var getBalanceCalled int
		var getTransactionCountCalled int
		var gasPriceCalled int
		var getCodeCalled int
		var estimateGasCalled int
		var sendTransactionCalled int
		var blockByNumberCalled int
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			body := make(map[string]interface{})
			err := json.NewDecoder(r.Body).Decode(&body)
			require.NoError(t, err)
			defer require.NoError(t, r.Body.Close())

			switch body["method"] {
			case "eth_getBalance":
				getBalanceCalled++
				require.Equal(t, strings.ToLower("0x"+walletAddr), body["params"].([]interface{})[0].(string))
				balance := hexutil.Big(*new(big.Int).Mul(big.NewInt(32*1e9), big.NewInt(1e9)))
				resp, err := json.Marshal(map[string]interface{}{
					"result": balance.String(),
				})
				require.NoError(t, err)
				_, err = w.Write(resp)
				require.NoError(t, err)
			case "eth_getTransactionCount":
				getTransactionCountCalled++
				require.NotEmpty(t, body["params"].([]interface{})[0].(string))
				resp, err := json.Marshal(map[string]interface{}{
					"result": hexutil.Uint64(5206).String(),
				})
				require.NoError(t, err)
				_, err = w.Write(resp)
				require.NoError(t, err)
			case "eth_gasPrice":
				gasPriceCalled++
				balance := hexutil.Big(*big.NewInt(1e9))
				resp, err := json.Marshal(map[string]interface{}{
					"result": balance.String(),
				})
				require.NoError(t, err)
				_, err = w.Write(resp)
				require.NoError(t, err)
			case "eth_getCode":
				getCodeCalled++
				require.NotEmpty(t, body["params"].([]interface{})[0].(string))
				code := hexutil.Bytes("5206")
				resp, err := json.Marshal(map[string]interface{}{
					"result": code.String(),
				})
				require.NoError(t, err)
				_, err = w.Write(resp)
				require.NoError(t, err)
			case "eth_estimateGas":
				estimateGasCalled++
				resp, err := json.Marshal(map[string]interface{}{
					"result": hexutil.Uint64(1).String(),
				})
				require.NoError(t, err)
				_, err = w.Write(resp)
				require.NoError(t, err)
			case "eth_sendRawTransaction":
				sendTransactionCalled++
				resp, err := json.Marshal(map[string]interface{}{
					"result": "success",
				})
				require.NoError(t, err)
				_, err = w.Write(resp)
				require.NoError(t, err)
			case "eth_getBlockByNumber":
				blockByNumberCalled++
				resp, err := json.Marshal(blkByNumberRes())
				require.NoError(t, err)
				_, err = w.Write(resp)
				require.NoError(t, err)
			case "eth_maxPriorityFeePerGas":
				gasPriceCalled++
				resp, err := json.Marshal(map[string]interface{}{
					"result": hexutil.Uint64(100).String(),
				})
				require.NoError(t, err)
				_, err = w.Write(resp)
				require.NoError(t, err)
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
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		require.NoError(t, err)
		require.Equal(t, 1, getBalanceCalled)
		require.Equal(t, 1, getTransactionCountCalled)
		require.Equal(t, 1, gasPriceCalled)
		// require.Equal(t, 1, getCodeCalled)
		// require.Equal(t, 1, estimateGasCalled)
		require.Equal(t, 1, sendTransactionCalled)
		require.NotEmpty(t, resultOut.String())
	})

	t.Run("successfully create one validator for one seed (mainnet)", func(t *testing.T) {
		var getBalanceCalled int
		var getTransactionCountCalled int
		var gasPriceCalled int
		var getCodeCalled int
		var estimateGasCalled int
		var sendTransactionCalled int
		var blockByNumberCalled int
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
				_, err = w.Write(resp)
				require.NoError(t, err)
			case "eth_getTransactionCount":
				getTransactionCountCalled++
				require.NotEmpty(t, body["params"].([]interface{})[0].(string))
				resp, err := json.Marshal(map[string]interface{}{
					"result": hexutil.Uint64(5206).String(),
				})
				require.NoError(t, err)
				_, err = w.Write(resp)
				require.NoError(t, err)
			case "eth_gasPrice":
				gasPriceCalled++
				balance := hexutil.Big(*big.NewInt(1e9))
				resp, err := json.Marshal(map[string]interface{}{
					"result": balance.String(),
				})
				require.NoError(t, err)
				_, err = w.Write(resp)
				require.NoError(t, err)
			case "eth_getCode":
				getCodeCalled++
				require.NotEmpty(t, body["params"].([]interface{})[0].(string))
				code := hexutil.Bytes("5206")
				resp, err := json.Marshal(map[string]interface{}{
					"result": code.String(),
				})
				require.NoError(t, err)
				_, err = w.Write(resp)
				require.NoError(t, err)
			case "eth_estimateGas":
				estimateGasCalled++
				resp, err := json.Marshal(map[string]interface{}{
					"result": hexutil.Uint64(1).String(),
				})
				require.NoError(t, err)
				_, err = w.Write(resp)
				require.NoError(t, err)
			case "eth_sendRawTransaction":
				sendTransactionCalled++
				resp, err := json.Marshal(map[string]interface{}{
					"result": "success",
				})
				require.NoError(t, err)
				_, err = w.Write(resp)
				require.NoError(t, err)
			case "eth_getBlockByNumber":
				blockByNumberCalled++
				resp, err := json.Marshal(blkByNumberRes())
				require.NoError(t, err)
				_, err = w.Write(resp)
				require.NoError(t, err)
			case "eth_maxPriorityFeePerGas":
				gasPriceCalled++
				resp, err := json.Marshal(map[string]interface{}{
					"result": hexutil.Uint64(100).String(),
				})
				require.NoError(t, err)
				_, err = w.Write(resp)
				require.NoError(t, err)
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
			"--network=mainnet",
		})
		err := cmd.RootCmd.Execute()
		require.NoError(t, err)
		require.Equal(t, 1, getBalanceCalled)
		require.Equal(t, 1, getTransactionCountCalled)
		require.Equal(t, 1, gasPriceCalled)
		// require.Equal(t, 1, getCodeCalled)
		// require.Equal(t, 1, estimateGasCalled)
		require.Equal(t, 1, sendTransactionCalled)
		require.Equal(t, 1, blockByNumberCalled)
		require.NotEmpty(t, resultOut.String())
	})

	t.Run("failed with insufficient funds", func(t *testing.T) {
		var getBalanceCalled int
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			body := make(map[string]interface{})
			err := json.NewDecoder(r.Body).Decode(&body)
			require.NoError(t, err)
			defer require.NoError(t, r.Body.Close())

			switch body["method"] {
			case "eth_getBalance":
				getBalanceCalled++
				require.Equal(t, strings.ToLower("0x"+walletAddr), body["params"].([]interface{})[0].(string))
				balance := hexutil.Big(*new(big.Int).Mul(big.NewInt(1*1e9), big.NewInt(1e9)))
				resp, err := json.Marshal(map[string]interface{}{
					"result": balance.String(),
				})
				require.NoError(t, err)
				_, err = w.Write(resp)
				require.NoError(t, err)
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
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		require.Error(t, err)
		require.EqualError(t, err, "insufficient funds for transfer")
		require.Equal(t, 1, getBalanceCalled)
	})

	t.Run("failed with invalid wallet address", func(t *testing.T) {
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
			"--wallet-addr", "invalidwalletaddr",
			"--validators-per-seed", "1",
			"--seeds-count", "1",
			"--web3-addr", "http://test.test",
			"--network=prater",
		})
		err := cmd.RootCmd.Execute()
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to HEX decode the given wallet address")
	})
}
