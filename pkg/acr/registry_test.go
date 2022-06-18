package acr

import (
	"reflect"
	"testing"
)

func Test_parseServerURL(t *testing.T) {
	type args struct {
		rawURL string
	}
	tests := []struct {
		name    string
		args    args
		want    *Registry
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseServerURL(tt.args.rawURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseServerURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseServerURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
