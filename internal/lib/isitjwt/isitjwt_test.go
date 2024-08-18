package isitjwt_test

import (
	"sync"
	"testing"
	"time"

	"github.com/Onnywrite/ssonny/internal/lib/isitjwt"
	"github.com/Onnywrite/ssonny/internal/lib/tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type SignSuite struct {
	suite.Suite
	uid     uuid.UUID
	secret  string
	subject string
	exp     time.Duration
}

func (s *SignSuite) SetupTest() {
	s.uid = uuid.New()
	s.secret = isitjwt.TODOSecret
	s.subject = "test"
	s.exp = time.Hour
}

func (s *SignSuite) TestHappyPath() {
	_, err := isitjwt.Sign(s.secret, s.uid, s.subject, s.exp)
	s.NoError(err)
}

func (s *SignSuite) TestShortSecret() {
	_, err := isitjwt.Sign("secret", s.uid, s.subject, s.exp)
	s.ErrorIs(err, isitjwt.ErrSecretTooShort)
}

type VerifySuite struct {
	suite.Suite
	secret     string
	subject    string
	validToken string
}

func (s *VerifySuite) SetupTest() {
	s.secret = isitjwt.TODOSecret
	s.subject = "test"
	s.validToken, _ = isitjwt.Sign(s.secret, uuid.New(), s.subject, time.Hour)
}

func (s *VerifySuite) TestWrongPeriodCount() {
	tests := []struct {
		token string
		err   error
	}{
		{"token", isitjwt.ErrInvalidToken},
		{"token.1", isitjwt.ErrInvalidToken},
		{"token.1.2", isitjwt.ErrInvalidToken},
		{"token.1.2.3", isitjwt.ErrInvalidToken},
		{"token.1.2.3.4", isitjwt.ErrInvalidToken},
		{"token.1.2.3.4.5", isitjwt.ErrInvalidToken},
	}
	for _, tc := range tests {
		_, err := isitjwt.Verify(s.secret, s.subject, tc.token)
		s.ErrorIs(err, tc.err)
	}
}

func (s *VerifySuite) TestShortSecret() {
	_, err := isitjwt.Verify("secret", s.subject, s.validToken)
	s.ErrorIs(err, isitjwt.ErrSecretTooShort)
}

func (s *VerifySuite) TestDecodingSignatureError() {
	invalidToken := s.validToken + "&"
	_, err := isitjwt.Verify(s.secret, s.subject, invalidToken)
	s.ErrorIs(err, isitjwt.ErrInvalidToken)
}

type E2ESuite struct {
	suite.Suite
	uid     uuid.UUID
	secret  string
	subject string
	exp     time.Duration
}

func (s *E2ESuite) SetupTest() {
	s.uid = uuid.New()
	s.secret = isitjwt.TODOSecret
	s.subject = "test"
	s.exp = time.Hour
}

func (s *E2ESuite) TestHappyPath() {
	token, err := isitjwt.Sign(s.secret, s.uid, s.subject, s.exp)
	s.NoError(err)
	s.NotEmpty(token)

	uid, err := isitjwt.Verify(s.secret, s.subject, token)
	s.NoError(err)
	s.Equal(s.uid, uid)
}

func (s *E2ESuite) TestExpired() {
	token, err := isitjwt.Sign(s.secret, s.uid, s.subject, -time.Hour)
	s.NoError(err)
	s.NotEmpty(token)

	_, err = isitjwt.Verify(s.secret, s.subject, token)
	s.ErrorIs(err, isitjwt.ErrTokenExpired)
}

func (s *E2ESuite) TestWrongSubject() {
	token, err := isitjwt.Sign(s.secret, s.uid, "email", s.exp)
	s.NoError(err)
	s.NotEmpty(token)

	_, err = isitjwt.Verify(s.secret, "test", token)
	s.ErrorIs(err, isitjwt.ErrInvalidSubject)
}

func TestAllIsitjwt(t *testing.T) {
	wg := sync.WaitGroup{}
	tests.RunSuitsParallel(&wg, t, new(SignSuite), new(VerifySuite), new(E2ESuite))
	wg.Wait()
}
