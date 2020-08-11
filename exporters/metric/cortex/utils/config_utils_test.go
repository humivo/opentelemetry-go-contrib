// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils_test

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/contrib/exporters/metric/cortex"
	"go.opentelemetry.io/contrib/exporters/metric/cortex/utils"
)

// initYAML creates a YAML file at a given filepath in a in-memory file system.
func initYAML(yamlBytes []byte, path string) (afero.Fs, error) {
	// Create an in-memory file system.
	fs := afero.NewMemMapFs()

	// Retrieve the directory path.
	dirPath := filepath.Dir(path)

	// Create the directory and the file in the directory.
	if err := fs.MkdirAll(dirPath, 0755); err != nil {
		return nil, err
	}
	if err := afero.WriteFile(fs, path, yamlBytes, 0644); err != nil {
		return nil, err
	}

	return fs, nil
}

// TestNewConfig tests whether NewConfig returns a correct Config struct. It checks whether the YAML
// file was read correctly and whether validation of the struct succeeded.
func TestNewConfig(t *testing.T) {
	tests := []struct {
		testName       string
		yamlByteString []byte
		fileName       string
		directoryPath  string
		expectedConfig *cortex.Config
		expectedError  error
	}{
		{
			testName:       "Valid Config file",
			yamlByteString: validYAML,
			fileName:       "config.yml",
			directoryPath:  "/test",
			expectedConfig: &validConfig,
			expectedError:  nil,
		},
		{
			testName:       "No Timeout",
			yamlByteString: noTimeoutYAML,
			fileName:       "config.yml",
			directoryPath:  "/test",
			expectedConfig: &validConfig,
			expectedError:  nil,
		},
		{
			testName:       "No Endpoint URL",
			yamlByteString: noEndpointYAML,
			fileName:       "config.yml",
			directoryPath:  "/test",
			expectedConfig: &validConfig,
			expectedError:  nil,
		},
		{
			testName:       "Two passwords",
			yamlByteString: twoPasswordsYAML,
			fileName:       "config.yml",
			directoryPath:  "/test",
			expectedConfig: nil,
			expectedError:  cortex.ErrTwoPasswords,
		},
		{
			testName:       "Two Bearer Tokens",
			yamlByteString: twoBearerTokensYAML,
			fileName:       "config.yml",
			directoryPath:  "/test",
			expectedConfig: nil,
			expectedError:  cortex.ErrTwoBearerTokens,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Create YAML file.
			fullPath := test.directoryPath + "/" + test.fileName
			fs, err := initYAML(test.yamlByteString, fullPath)
			require.Nil(t, err)

			// Create new Config struct from the specified YAML file with an in-memory filesystem.
			config, err := utils.NewConfig(
				test.fileName,
				utils.WithFilepath(test.directoryPath),
				utils.WithFilesystem(fs),
			)

			// Verify error and struct contents.
			require.Equal(t, err, test.expectedError)
			require.Equal(t, config, test.expectedConfig)
		})
	}
}
