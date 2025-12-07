package usecase

import (
	"db/model"
	"testing"
)

type MockRepo struct{}

func (m *MockRepo) FindByName(name string) ([]model.User, error) { return nil, nil }
func (m *MockRepo) Insert(user *model.User) error                { return nil }

func TestUserUsecase_validateRegisterRequest(t *testing.T) {
	u := NewUserUsecase(&MockRepo{})
	tests := []struct {
		name    string
		req     RegisterUserReq
		wantErr bool
	}{
		{"正常", RegisterUserReq{Name: "Taro", Age: 25}, false},
		{"名前空", RegisterUserReq{Name: "", Age: 25}, true},
		{"名前長すぎ", RegisterUserReq{Name: "123456789012345678901234567890123456789012345678901", Age: 25}, true},
		{"年齢若すぎ", RegisterUserReq{Name: "Jiro", Age: 19}, true},
		{"年齢高すぎ", RegisterUserReq{Name: "Saburo", Age: 81}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := u.validateRegisterRequest(tt.req); (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
