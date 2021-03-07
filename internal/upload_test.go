package internal

import "testing"

func TestUpload(t *testing.T) {
	type args struct {
		fp  string
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test upload",
			args: args{
				fp:  "../output/fs",
				url: "http://localhost:8002/upload",
			},
			want:    "OK",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Upload(tt.args.fp, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Upload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Upload() got = %v, want %v", got, tt.want)
			}
		})
	}
}
