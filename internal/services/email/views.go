package email

type VerificationEmail struct {
	Recipient    string
	UserNickname string
	Link         string
}

type PasswordResetEmail struct {
	Recipient    string
	UserNickname string
	Link         string
}
