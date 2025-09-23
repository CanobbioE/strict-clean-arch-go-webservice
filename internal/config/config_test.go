package config_test

import (
	"errors"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/config"
)

func TestLoadConfig(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name       string
		args       args
		want       *config.ServiceCfg
		wantErr    bool
		compareErr func(error) bool
	}{
		{
			name: "file doesn't exists",
			args: args{
				path: "/this/does/not/exist",
			},
			want:    nil,
			wantErr: true,
			compareErr: func(err error) bool {
				return errors.Is(err, os.ErrNotExist)
			},
		},
		{
			name: "file is not parsable",
			args: args{
				path: "./testing/unparsable.yaml",
			},
			want:    nil,
			wantErr: true,
			compareErr: func(err error) bool {
				return strings.Contains(err.Error(), "cannot unmarshal")
			},
		},
		{
			name: "success",
			args: args{
				path: "./testing/sample-config.yaml",
			},
			want:    &config.ServiceCfg{ServerAddress: ":8080"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := config.Load(tt.args.path)
			if (tt.wantErr != (err != nil)) || (tt.wantErr && !tt.compareErr(err)) {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() got = %v, want %v", got, tt.want)
			}
		})
	}
}
