package testutils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func WithJSONFile(t testing.TB, in interface{}, f func(path string)) {
	dir, err := ioutil.TempDir("", "fitbitsleeptest")
	require.NoError(t, err, "Unable to create temp directory)")
	file, err := ioutil.TempFile(dir, "jsonfile")
	require.NoError(t, err, "Couldn't create tempfile")
	e := json.NewEncoder(file)
	require.NoError(t, e.Encode(in), "unable to encode JSON from struct")
	defer os.Remove(file.Name())
	defer os.RemoveAll(dir)

	f(file.Name())
}
