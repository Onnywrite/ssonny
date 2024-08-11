package email

type VerificationEmail struct {
	Recipient    string
	UserNickname string
	Token        string
}
