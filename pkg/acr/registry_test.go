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
		{
			name: "person public",
			args: args{
				rawURL: "registry.cn-beijing.aliyuncs.com/foo/bar:v1",
			},
			want: &Registry{
				IsEE:         false,
				InstanceId:   "",
				InstanceName: "",
				Region:       "cn-beijing",
				Domain:       "registry.cn-beijing.aliyuncs.com",
			},
			wantErr: false,
		},
		{
			name: "person vpc",
			args: args{
				rawURL: "registry-vpc.cn-beijing.aliyuncs.com/foo/bar:v1",
			},
			want: &Registry{
				IsEE:         false,
				InstanceId:   "",
				InstanceName: "",
				Region:       "cn-beijing",
				Domain:       "registry-vpc.cn-beijing.aliyuncs.com",
			},
			wantErr: false,
		},
		{
			name: "person internal",
			args: args{
				rawURL: "registry-internal.cn-beijing.aliyuncs.com/foo/bar:v1",
			},
			want: &Registry{
				IsEE:         false,
				InstanceId:   "",
				InstanceName: "",
				Region:       "cn-beijing",
				Domain:       "registry-internal.cn-beijing.aliyuncs.com",
			},
			wantErr: false,
		},
		{
			name: "person public (intl)",
			args: args{
				rawURL: "registry-intl.ap-southeast-1.aliyuncs.com/foo/bar:v1",
			},
			want: &Registry{
				IsEE:         false,
				InstanceId:   "",
				InstanceName: "",
				Region:       "ap-southeast-1",
				Domain:       "registry-intl.ap-southeast-1.aliyuncs.com",
			},
			wantErr: false,
		},
		{
			name: "person vpc (intl)",
			args: args{
				rawURL: "registry-intl-vpc.ap-southeast-1.aliyuncs.com/foo/bar:v1",
			},
			want: &Registry{
				IsEE:         false,
				InstanceId:   "",
				InstanceName: "",
				Region:       "ap-southeast-1",
				Domain:       "registry-intl-vpc.ap-southeast-1.aliyuncs.com",
			},
			wantErr: false,
		},
		{
			name: "ee public",
			args: args{
				rawURL: "foobar-registry.cn-beijing.cr.aliyuncs.com/foo/bar:v2",
			},
			want: &Registry{
				IsEE:         true,
				InstanceId:   "",
				InstanceName: "foobar",
				Region:       "cn-beijing",
				Domain:       "foobar-registry.cn-beijing.cr.aliyuncs.com",
			},
			wantErr: false,
		},
		{
			name: "ee vpc",
			args: args{
				rawURL: "foobar-registry-vpc.cn-beijing.cr.aliyuncs.com/foo/bar:v2",
			},
			want: &Registry{
				IsEE:         true,
				InstanceId:   "",
				InstanceName: "foobar",
				Region:       "cn-beijing",
				Domain:       "foobar-registry-vpc.cn-beijing.cr.aliyuncs.com",
			},
			wantErr: false,
		},
		{
			name: "unknown domain",
			args: args{
				rawURL: "alpine:3.15",
			},
			want:    nil,
			wantErr: true,
		},
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
