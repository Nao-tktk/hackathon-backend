package usecase

import (
	"fmt"

	"db/model"
)

type UserUsecase struct {
	Repo UserRepository
}

func NewUserUsecase(repo UserRepository) *UserUsecase {
	return &UserUsecase{Repo: repo}
}

func (u *UserUsecase) SearchUser(name string) ([]model.User, error) {
	if name == "" {
		return nil, fmt.Errorf("name is empty")
	}
	return u.Repo.FindByName(name)
}

type RegisterUserReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (u *UserUsecase) RegisterUser(req RegisterUserReq) (int, error) {
	if err := u.validateRegisterRequest(req); err != nil {
		return 0, err
	}

	user := &model.User{
		Name:     req.Name,
		Password: req.Password,
	}

	id, err := u.Repo.Insert(user)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (u *UserUsecase) validateRegisterRequest(req RegisterUserReq) error {
	if req.Name == "" {
		return fmt.Errorf("invalid name: empty")
	}
	if len(req.Name) > 50 {
		return fmt.Errorf("invalid name: too long")
	}
	if len(req.Password) < 4 {
		return fmt.Errorf("invalid password: password too short; password must have at least 4 characters")
	}
	return nil
}

type LoginReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Login: 名前で検索し、パスワードが一致するか確認する
func (u *UserUsecase) Login(req LoginReq) (int, error) {
	// 1. 名前でユーザーを探す
	users, err := u.Repo.FindByName(req.Name)
	if err != nil {
		return 0, err
	}
	if len(users) == 0 {
		return 0, fmt.Errorf("user not found")
	}

	// 2. パスワード照合 (ハッカソン用: 平文比較)
	// FindByNameはリストを返すので、先頭のユーザーを使います
	targetUser := users[0]
	if targetUser.Password != req.Password {
		return 0, fmt.Errorf("invalid password")
	}

	return targetUser.ID, nil
}

// ▼▼▼ 追加: ソーシャルログイン用リクエスト型 ▼▼▼
type SocialLoginReq struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// ▼▼▼ 追加: ソーシャルログインのロジック ▼▼▼
func (u *UserUsecase) SocialLogin(req SocialLoginReq) (int, string, error) {
	// 1. そのメールアドレス(Nameカラムに保存)のユーザーがいるか確認
	users, err := u.Repo.FindByName(req.Email)
	if err != nil {
		return 0, "", err
	}

	// 2. 既にいるなら、そのユーザー情報を返す
	if len(users) > 0 {
		return users[0].ID, users[0].Name, nil
	}

	// 3. いないなら、新規作成する
	// ハッカソン仕様: Nameカラムにメアドを入れ、パスワードはダミーを入れる
	newUser := &model.User{
		Name:     req.Email,
		Password: "google_login", // ソーシャルログイン用のダミーパスワード
	}

	id, err := u.Repo.Insert(newUser)
	if err != nil {
		return 0, "", err
	}

	// 表示名はリクエストのName(Googleの名前)を返したいが、
	// DB上はEmailで管理しているので、ここでは便宜上EmailをNameとして扱うか、
	// フロントの表示用に req.Name を返す設計にします。
	return id, req.Name, nil
}
