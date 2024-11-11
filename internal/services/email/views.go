package email

type VerificationEmail struct {
	Recipient    string
	UserNickname string
	Token        string
}

type Notification struct {
	Recipient    string
	UserNickname string
	Message      string
}
