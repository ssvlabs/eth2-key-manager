package printer_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/eth2-key-manager/cli/util/printer"
)

func TestText(t *testing.T) {
	text := "some text"

	var buf bytes.Buffer
	pr := printer.New(&buf)
	pr.Text(text)
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
	pr := printer.New(&buf)
	err = pr.JSON(obj)
	require.NoError(t, err)
	require.Equal(t, string(expectedData)+"\n", buf.String())
}

func TestError(t *testing.T) {
	testError := errors.New("some error text")

	var buf bytes.Buffer
	pr := printer.New(&buf)
	pr.Error(testError)
	require.Equal(t, "Error: "+testError.Error()+"\n", buf.String())
}
