package testutils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func WithFile(t testing.TB, dir string, fn func(*os.File)) {
	f, err := ioutil.TempFile(dir, "withfilefitbit")
	require.NoError(t, err)
	defer os.Remove(f.Name())
	fn(f)
}

func WithJSONFile(t testing.TB, in interface{}, f func(path string)) {
	WithTempDir(t, func(dir string) {
		WithFile(t, dir, func(file *os.File) {
			e := json.NewEncoder(file)
			require.NoError(t, e.Encode(in), "unable to encode JSON from struct")
			f(file.Name())
		})
	})
}

func WithTempDir(t testing.TB, fn func(path string)) {
	dir, err := ioutil.TempDir("", "fitbitsleeptest")
	require.NoError(t, err, "Unable to create temp directory)")
	defer os.RemoveAll(dir)

	fn(dir)
}
