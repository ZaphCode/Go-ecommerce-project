package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_jwtAuthServiceImpl_CreateToken(t *testing.T) {
	type args struct {
		claims        Claims
		exp           time.Duration
		isRefreshType bool
	}
	tests := []struct {
		name    string
		s       *jwtAuthServiceImpl
		args    args
		wantErr bool
	}{
		{
			name: "Correct funtionality ",
			s:    &jwtAuthServiceImpl{},
			args: args{
				claims: Claims{
					ID:   uuid.New(),
					Role: "user",
				},
				exp:           time.Millisecond * 100,
				isRefreshType: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.CreateToken(tt.args.claims, tt.args.exp, tt.args.isRefreshType)
			if (err != nil) != tt.wantErr {
				t.Errorf("jwtAuthServiceImpl.CreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log("Success!", got)
		})
	}
}

func Test_jwtAuthServiceImpl_CreateTokens(t *testing.T) {
	type args struct {
		claims      Claims
		access_exp  time.Duration
		refresh_exp time.Duration
	}
	tests := []struct {
		name    string
		s       *jwtAuthServiceImpl
		args    args
		wantErr bool
	}{
		{
			name: "Success case (create both tokens)",
			s:    &jwtAuthServiceImpl{},
			args: args{
				claims: Claims{
					ID:   uuid.New(),
					Role: "user",
				},
				access_exp:  time.Hour * 24 * 30,
				refresh_exp: time.Hour * 25 * 30,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			at, rt, err := tt.s.CreateTokens(tt.args.claims, tt.args.access_exp, tt.args.refresh_exp)
			if (err != nil) != tt.wantErr {
				t.Errorf("jwtAuthServiceImpl.CreateTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log("Success!")
			t.Log("Access!", at)
			t.Log("Refresh!", rt)
		})
	}
}

func Test_jwtAuthServiceImpl_DecodeToken(t *testing.T) {
	type args struct {
		jwtoken     string
		refreshType bool
	}
	tests := []struct {
		name    string
		s       *jwtAuthServiceImpl
		args    args
		want    *Claims
		wantErr bool
	}{
		{
			name: "Decode normal access token",
			s:    &jwtAuthServiceImpl{},
			args: args{
				jwtoken:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6ImNjYjkwMTVmLTFkMWEtNDI2YS05OTY5LTM1NmZhZWQwMDZkNCIsIlJvbGUiOiJ1c2VyIiwiaXNzIjoiY2NiOTAxNWYtMWQxYS00MjZhLTk5NjktMzU2ZmFlZDAwNmQ0IiwiZXhwIjoxNzA0NTIwNTg0fQ.Wep1wFETQ2qX5Bny4QzZLSLGadbT4BH4vD24vheA3FE",
				refreshType: false,
			},
			wantErr: false,
		},
		{
			name: "Decode normal refresh token",
			s:    &jwtAuthServiceImpl{},
			args: args{
				jwtoken:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6ImNjYjkwMTVmLTFkMWEtNDI2YS05OTY5LTM1NmZhZWQwMDZkNCIsIlJvbGUiOiJ1c2VyIiwiaXNzIjoiY2NiOTAxNWYtMWQxYS00MjZhLTk5NjktMzU2ZmFlZDAwNmQ0IiwiZXhwIjoxNzA1ODE2NTg0fQ.1rxjJ8yt6pWELx1VtgpDESrwMtZealf-3XOPM3N2Q3M",
				refreshType: true,
			},
			wantErr: false,
		},
		{
			name: "Decode expired access token",
			s:    &jwtAuthServiceImpl{},
			args: args{
				jwtoken:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IjVkNWFiMjUzLWIzZTctNDc5ZC1hNzY2LTg0ZjNjNDk1MDdjZCIsIlJvbGUiOiJ1c2VyIiwiaXNzIjoiNWQ1YWIyNTMtYjNlNy00NzlkLWE3NjYtODRmM2M0OTUwN2NkIiwiZXhwIjoxNjczNDE3MTU1fQ.t48Q9o26qpbVeCfCOEMHJwfi3Tot3TIDKO3qRWhx-z0",
				refreshType: false,
			},
			wantErr: true,
		},
		{
			name: "Decode invalid access token",
			s:    &jwtAuthServiceImpl{},
			args: args{
				jwtoken:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6ImY0MzY3MDFjLWE4ZGYtNDA0MC1hOWRlLTcxM2I5MWIwYzI5NSIsIlJvbGUiOiJ1c2VyIiwiaXNzIjoiZjQzNjcwMWMtYThkZi00MDQwLWE5ZGUtNzEzYjkxYjBjMjk1IiwiZXhwIjoxNjczNDE2ODQ5fQ.dfuMU0sxE0MJVgik2a2cyAhbyjpdsfP2p7X2gBwAJ8Y",
				refreshType: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.DecodeToken(tt.args.jwtoken, tt.args.refreshType)
			if (err != nil) != tt.wantErr {
				t.Errorf("jwtAuthServiceImpl.DecodeToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("jwtAuthServiceImpl.DecodeToken() = %v, want %v", got, tt.want)
			// }
			t.Logf("Success! %+v", got)
		})
	}
}
