package user

import (
	"testing"

	"github.com/ZaphCode/clean-arch/src/utils"
)

func Test_memoryUserRepositoryImpl_Find(t *testing.T) {
	repo := NewMemoryUserRepository()

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Simple find",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUsers, err := repo.Find()
			if (err != nil) != tt.wantErr {
				t.Errorf("memoryUserRepositoryImpl.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			utils.PrettyPrint(gotUsers)
		})
	}
}
