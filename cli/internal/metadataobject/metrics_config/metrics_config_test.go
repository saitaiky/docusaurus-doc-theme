package metricsconfig

import (
	"fmt"
	"io/ioutil"
	"testing"

	goyaml "github.com/goccy/go-yaml"
	"github.com/hasura/graphql-engine/cli/v2/internal/metadatautil"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestMetricsConfigBuild(t *testing.T) {
	type fields struct {
		MetadataDir string
		logger      *logrus.Logger
	}
	tests := []struct {
		name       string
		fields     fields
		wantGolden string
		wantErr    bool
		assertErr  require.ErrorAssertionFunc
	}{
		{
			"can build metrics config",
			fields{
				MetadataDir: "testdata/build_test/t1/metadata",
				logger:      logrus.New(),
			},
			"testdata/build_test/t1/want.golden.json",
			false,
			require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := &MetricsConfigObject{
				MetadataDir: tt.fields.MetadataDir,
				logger:      tt.fields.logger,
			}
			got, err := tc.Build()
			tt.assertErr(t, err)
			if !tt.wantErr {
				assert.NoError(t, err)
				gotbs, err := yaml.Marshal(got)
				assert.NoError(t, err)
				jsonbs, err := goyaml.YAMLToJSON(gotbs)
				assert.NoError(t, err)

				// uncomment following lines to update golden file
				// assert.NoError(t, ioutil.WriteFile(tt.wantGolden, jsonbs, os.ModePerm))

				wantbs, err := ioutil.ReadFile(tt.wantGolden)
				assert.NoError(t, err)
				assert.Equal(t, string(wantbs), string(jsonbs))
			}
		})
	}
}

func TestMetricsConfigExport(t *testing.T) {
	type fields struct {
		MetadataDir string
		logger      *logrus.Logger
	}
	type args struct {
		metadata map[string]yaml.Node
	}
	tests := []struct {
		id        string
		name      string
		fields    fields
		args      args
		want      map[string][]byte
		wantErr   bool
		assertErr require.ErrorAssertionFunc
	}{
		{
			"t1",
			"can export metrics config",
			fields{
				MetadataDir: "testdata/metadata",
				logger:      logrus.New(),
			},
			args{
				metadata: func() map[string]yaml.Node {
					bs, err := ioutil.ReadFile("testdata/export_test/t1/metadata.json")
					assert.NoError(t, err)
					yamlbs, err := metadatautil.JSONToYAML(bs)
					assert.NoError(t, err)
					var v map[string]yaml.Node
					assert.NoError(t, yaml.Unmarshal(yamlbs, &v))
					return v
				}(),
			},
			func() map[string][]byte {
				m := map[string][]byte{
					"testdata/metadata/metrics_config.yaml": func() []byte {
						bs, err := ioutil.ReadFile("testdata/export_test/t1/want.metrics_config.yaml")
						assert.NoError(t, err)
						return bs
					}(),
				}
				return m
			}(),
			false,
			require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := &MetricsConfigObject{
				MetadataDir: tt.fields.MetadataDir,
				logger:      tt.fields.logger,
			}
			got, err := tc.Export(tt.args.metadata)
			tt.assertErr(t, err)
			for k, v := range got {
				fmt.Println(k, string(v))
			}
			if !tt.wantErr {
				for k, v := range got {
					assert.Contains(t, tt.want, k)
					// uncomment to update golden files
					//assert.NoError(t, ioutil.WriteFile(fmt.Sprintf("testdata/export_test/%v/want.%v", tt.id, filepath.Base(k)), v, os.ModePerm))

					assert.Equalf(t, string(tt.want[k]), string(v), "%v", k)
				}
			}
		})
	}
}
