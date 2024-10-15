package tokens

import (
	"sync"
	"testing"
	"time"

	"github.com/Onnywrite/ssonny/internal/lib/tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type SignSuite struct {
	suite.Suite
	uid uuid.UUID
	gen Generator
}

func (s *SignSuite) SetupTest() {
	s.uid = uuid.New()
	s.gen = NewWithConfig(Config{
		Issuer:     "",
		Secret:     []byte("SecrEt"),
		AccessExp:  time.Hour,
		RefreshExp: time.Hour,
		IdExp:      time.Hour,
		EmailExp:   time.Hour,
	})
}

func (s *SignSuite) TestHappyPath() {
	_, err := s.gen.SignEmail(s.uid)
	s.NoError(err)
}

type VerifySuite struct {
	suite.Suite
	gen        Generator
	validToken string
}

func (s *VerifySuite) SetupTest() {
	s.gen = NewWithConfig(Config{
		Issuer:     "",
		Secret:     []byte("SecrEt"),
		AccessExp:  time.Hour,
		RefreshExp: time.Hour,
		IdExp:      time.Hour,
		EmailExp:   time.Hour,
	})
	s.validToken, _ = s.gen.SignEmail(uuid.New())
}

func (s *VerifySuite) TestWrongPeriodCount() {
	tests := []struct {
		token string
		err   error
	}{
		{"token", ErrInvalidToken},
		{"token.1", ErrInvalidToken},
		{"token.1.2", ErrInvalidToken},
		{"token.1.2.3", ErrInvalidToken},
		{"token.1.2.3.4", ErrInvalidToken},
		{"token.1.2.3.4.5", ErrInvalidToken},
	}
	for _, tc := range tests {
		s.Run(tc.token, func() {
			_, err := s.gen.ParseEmail(tc.token)
			s.ErrorIs(err, tc.err)
		})
	}
}

func (s *VerifySuite) TestDecodingSignatureError() {
	invalidToken := s.validToken + "&"
	_, err := s.gen.ParseEmail(invalidToken)
	s.ErrorIs(err, ErrInvalidToken)
}

type E2ESuite struct {
	suite.Suite
	gen Generator
	uid uuid.UUID
}

func (s *E2ESuite) SetupTest() {
	s.uid = uuid.New()
	s.gen = NewWithConfig(Config{
		Issuer:     "",
		Secret:     []byte("SecrEt"),
		AccessExp:  time.Hour,
		RefreshExp: time.Hour,
		IdExp:      time.Hour,
		EmailExp:   time.Hour,
	})
}

func (s *E2ESuite) TestHappyPath() {
	token, err := s.gen.SignEmail(s.uid)
	s.NoError(err)
	s.NotEmpty(token)

	uid, err := s.gen.ParseEmail(token)
	s.NoError(err)
	s.Equal(s.uid, uid)
}

func (s *E2ESuite) TestExpired() {
	expiredGen := NewWithConfig(Config{
		Issuer:     "",
		Secret:     []byte("SecrEt"),
		AccessExp:  time.Hour,
		RefreshExp: time.Hour,
		IdExp:      time.Hour,
		EmailExp:   -time.Hour,
	})
	token, err := expiredGen.SignEmail(s.uid)
	s.NoError(err)
	s.NotEmpty(token)

	_, err = s.gen.ParseEmail(token)
	s.ErrorIs(err, ErrExpired)
}

func (s *E2ESuite) TestWrongSubject() {
	token, err := s.gen.signMyToken(s.uid.String(), "notemail", time.Hour, 1)
	s.NoError(err)
	s.NotEmpty(token)

	_, err = s.gen.ParseEmail(token)
	s.ErrorIs(err, ErrInvalidSubject)
}

func (s *E2ESuite) TestWrongSignature() {
	wrongSecretGen := NewWithConfig(Config{
		Issuer:     "",
		Secret:     []byte("wrongSecret"),
		AccessExp:  time.Hour,
		RefreshExp: time.Hour,
		IdExp:      time.Hour,
		EmailExp:   time.Hour,
	})
	token, err := wrongSecretGen.SignEmail(s.uid)
	s.NoError(err)
	s.NotEmpty(token)

	_, err = s.gen.ParseEmail(token)
	s.ErrorIs(err, ErrInvalidToken)
}

func TestAllIsitjwt(t *testing.T) {
	wg := sync.WaitGroup{}
	tests.RunSuitsParallel(t, &wg, new(SignSuite), new(VerifySuite), new(E2ESuite))
	wg.Wait()
}
