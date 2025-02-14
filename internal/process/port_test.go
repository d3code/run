package process

import "testing"

func TestKillPortProcess(t *testing.T) {
	type args struct {
		port int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestKillPortProcess",
			args: args{port: 8080},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			KillPortProcess(tt.args.port)
		})
	}
}
