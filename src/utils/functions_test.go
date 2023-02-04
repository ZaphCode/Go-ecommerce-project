package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestItemInSlice(t *testing.T) {
	type args[T any] struct {
		a    T
		list []T
	}
	tests := []struct {
		name string
		args args[string]
		want bool
	}{
		{
			name: "in array",
			args: args[string]{
				a:    "test",
				list: []string{"t", "te", "test"},
			},
			want: true,
		},
		{
			name: "not in array",
			args: args[string]{
				a:    "uwu",
				list: []string{"t", "te", "test"},
			},
			want: false,
		},
		{
			name: "ivalid provider",
			args: args[string]{
				a:    "uwu",
				list: GetOAuthProviders(),
			},
			want: false,
		},
		{
			name: "valid role",
			args: args[string]{
				a:    UserRole,
				list: GetUserRoles(),
			},
			want: true,
		},
	}
	for i, tt := range tests {
		tt = tests[i]
		t.Run(tt.name, func(t *testing.T) {
			if got := ItemInSlice(tt.args.a, tt.args.list); got != tt.want {
				t.Errorf("ItemInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStructFields(t *testing.T) {
	type veicle struct {
		speed int
		Name  string
	}

	type moto struct {
		veicle
		Wheels int
		brand  string
	}

	type args struct {
		s interface{}
	}

	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "Get the composited field",
			args:    args{s: new(moto)},
			want:    []string{"veicle", "Wheels", "brand"},
			wantErr: false,
		},
		{
			name:    "Empty struct",
			args:    args{s: struct{}{}},
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "Empty struct pointer",
			args:    args{s: &struct{}{}},
			want:    []string{},
			wantErr: false,
		},
		{
			name: "Normal work",
			args: args{s: struct {
				Name   string
				Age    int
				Single bool
			}{"Carlos", 24, true}},
			want:    []string{"Name", "Age", "Single"},
			wantErr: false,
		},
		{
			name: "Normal work with pointer",
			args: args{s: &struct {
				Amount float32
				Tags   []string
			}{}},
			want:    []string{"Amount", "Tags"},
			wantErr: false,
		},
		{
			name: "Normal work with new() struct",
			args: args{s: new(struct {
				Amount int
				Name   string
			})},
			want:    []string{"Amount", "Name"},
			wantErr: false,
		},
		{
			name:    "Normal work with new() struct",
			args:    args{s: new(int)},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Not struct",
			args:    args{s: "sexo"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Nil value",
			args:    args{s: nil},
			want:    nil,
			wantErr: true,
		},
	}

	for i, tt := range tests {
		tt = tests[i]
		t.Run(tt.name, func(t *testing.T) {
			r := require.New(t)

			got, err := GetStructFields(tt.args.s)

			r.Equal(tt.wantErr, (err != nil), "expect err fail")

			r.Equalf(tt.want, got, "GetStructFields() = %v, want %v", got, tt.want)
		})
	}
}

func TestRandomString(t *testing.T) {
	testCases := []struct {
		desc string
		l    int
	}{
		{
			desc: "20 caracters length",
			l:    20,
		},
		{
			desc: "10 caracters length",
			l:    10,
		},
		{
			desc: "15 caracters length",
			l:    15,
		},
	}
	for i, tC := range testCases {
		tC = testCases[i]
		t.Run(tC.desc, func(t *testing.T) {
			rs := RandomString(tC.l)

			require.Equal(t, tC.l, len(rs), "ivalid length")

			t.Log(rs)
		})
	}
}

func TestGetStructAttr(t *testing.T) {
	type args struct {
		strc      interface{}
		fieldName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Proper work",
			args: args{
				strc:      struct{ Name string }{"John Doe"},
				fieldName: "Name",
			},
			wantErr: false,
		},
		{
			name: "Not field",
			args: args{
				strc:      struct{ Name string }{"John Doe"},
				fieldName: "Email",
			},
			wantErr: true,
		},
		{
			name: "Not field in pointer",
			args: args{
				strc:      &struct{ Name string }{"John Doe"},
				fieldName: "Email",
			},
			wantErr: true,
		},
		{
			name: "existing zero field",
			args: args{
				strc:      struct{ Name string }{},
				fieldName: "Name",
			},
			wantErr: false,
		},
		{
			name: "Nil input",
			args: args{
				strc:      nil,
				fieldName: "Name",
			},
			wantErr: true,
		},
		{
			name: "Not struct",
			args: args{
				strc:      "Sexo",
				fieldName: "Name",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := GetStructAttr(tt.args.strc, tt.args.fieldName)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetStructAttr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				t.Log(got.Interface())
				t.Log(got.Type())
			}
		})
	}
}

func TestIsZeroValue(t *testing.T) {
	testCases := []struct {
		desc   string
		value  interface{}
		expect bool
	}{
		{
			desc:   "no zero string",
			value:  "hello",
			expect: false,
		},
		{
			desc:   "zero string",
			value:  "",
			expect: true,
		},
		{
			desc:   "no zero int",
			value:  10,
			expect: false,
		},
		{
			desc:   "zero int",
			value:  0,
			expect: true,
		},
		{
			desc:   "no zero float",
			value:  float32(154),
			expect: false,
		},
		{
			desc:   "no zero struct",
			value:  struct{ A int }{1},
			expect: false,
		},
		{
			desc:   "zero struct",
			value:  struct{ A string }{},
			expect: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			require.Equal(t, tC.expect, isZeroValue(tC.value), "wrong result")
		})
	}
}
