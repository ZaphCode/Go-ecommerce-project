package user

import (
	"fmt"
	"testing"
	"time"

	"github.com/ZaphCode/clean-arch/config"
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/ZaphCode/clean-arch/src/utils"
)

func Test_firestoreUserRepositoryImpl_Find(t *testing.T) {
	config.MustLoadConfig("./../../../config")
	config.MustLoadFirebaseConfig("./../../../config")

	tests := []struct {
		name    string
		r       *firestoreUserRepositoryImpl
		want    []domain.User
		wantErr bool
	}{
		{
			name: "Normal work",
			r: &firestoreUserRepositoryImpl{
				client:   utils.GetFirestoreClient(config.GetFirebaseApp()),
				collName: "users",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			got, err := tt.r.Find()
			if (err != nil) != tt.wantErr {
				t.Fatalf("firestoreUserRepositoryImpl.Find() error = %v, wantErr %v", err, tt.wantErr)
			}
			utils.PrettyPrint(got)
			fmt.Println("Took:", time.Since(start).Milliseconds(), "ms")
		})
	}
}

// func Beanchmark_firestoreUserRepositoryImpl_Find(t *testing.B) {
// 	config.LoadConfig("./../../../config")
// 	config.LoadFirebaseConfig("./../../../config")

// 	repo := &firestoreUserRepositoryImpl{
// 		client:   utils.GetFirestoreClient(config.GetFirebaseApp()),
// 		collName: "users",
// 	}
// 	for i := 0; i < t.N; i++ {
// 		users, err := repo.Find()
// 		if err != nil {
// 			t.Errorf("firestoreUserRepositoryImpl.Find() error = %v, wantErr %v", err, tt.wantErr)
// 			return
// 		}
// 		utils.PrettyPrint(users)
// 	}
// }
