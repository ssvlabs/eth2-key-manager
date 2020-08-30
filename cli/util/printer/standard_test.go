package printer_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bloxapp/eth-key-manager/cli/util/printer"
)

func TestText(t *testing.T) {
	text := "some text"

	var buf bytes.Buffer
	printer := printer.New(&buf)
	printer.Text(text)
	require.Equal(t, text+"\n", buf.String())
}

func TestJSON(t *testing.T) {
	obj := struct {
		data string
	}{
		data: "some data",
	}
	expectedData, err := json.MarshalIndent(obj, "", "  ")
	require.NoError(t, err)

	var buf bytes.Buffer
	printer := printer.New(&buf)
	err = printer.JSON(obj)
	require.NoError(t, err)
	require.Equal(t, string(expectedData)+"\n", buf.String())
}

func TestError(t *testing.T) {
	testError := fmt.Errorf("some error text")

	var buf bytes.Buffer
	printer := printer.New(&buf)
	printer.Error(testError)
	require.Equal(t, "Error: "+testError.Error()+"\n", buf.String())
}
