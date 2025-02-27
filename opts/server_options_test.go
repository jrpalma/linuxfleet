package opts

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	testCases := []struct {
		name  string
		input ServerOptions
	}{
		{
			"BasicTest",
			ServerOptions{DatabaseCluster: []string{"db1", "db2"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.input.Marshal()
			assert.NoError(t, err)
		})
	}
}

func TestUnmarshal(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected ServerOptions
	}{
		{
			"BasicTest",
			[]byte("database_cluster: [db1, db2]\n"),
			ServerOptions{DatabaseCluster: []string{"db1", "db2"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result ServerOptions
			err := result.Unmarshal(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestWriteAndReadOptions(t *testing.T) {
	testCases := []struct {
		name     string
		input    ServerOptions
		expected ServerOptions
	}{
		{
			"BasicTest",
			ServerOptions{DatabaseCluster: []string{"db1", "db2"}},
			ServerOptions{DatabaseCluster: []string{"db1", "db2"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "test_options")
			assert.NoError(t, err)
			defer os.Remove(tmpFile.Name())

			err = tc.input.WriteOptions(tmpFile.Name())
			assert.NoError(t, err)

			var result ServerOptions
			err = result.ReadOptions(tmpFile.Name())
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}
