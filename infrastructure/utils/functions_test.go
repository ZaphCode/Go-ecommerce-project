package utils

import "testing"

func TestRandomString(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "omg",
			args: args{
				n: 20,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RandomString(tt.args.n)
			t.Logf("\n >>>>>> %s", got)
			if len(got) != tt.args.n {
				t.Fatal("bad string")
			}
		})
	}
}
