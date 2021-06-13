package models

// PasswordReset TODO: 有効期限を追加する
type PasswordReset struct {
	Id    uint
	Email string
	Token string
}
