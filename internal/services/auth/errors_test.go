package auth

import (
	"os"
	"testing"

	"github.com/Onnywrite/ssonny/internal/storage/repo"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestUserFailed(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		expErr error
	}{
		{
			name:   "empty_result",
			err:    repo.ErrEmptyResult,
			expErr: ErrUserNotFound,
		},
		{
			name:   "unique",
			err:    repo.ErrUnique,
			expErr: ErrUserAlreadyExists,
		},
		{
			name:   "checked",
			err:    repo.ErrChecked,
			expErr: ErrInvalidData,
		},
		{
			name:   "foreign_key",
			err:    repo.ErrFK,
			expErr: ErrDependencyNotFound,
		},
		{
			name:   "null",
			err:    repo.ErrNull,
			expErr: ErrInvalidData,
		},
		{
			name:   "internal",
			err:    gofakeit.Error(),
			expErr: ErrInternal,
		},
	}

	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)

	t.Parallel()

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			err := userFailed(&logger, tc.err)
			assert.ErrorIs(tt, err, tc.expErr)
		})
	}

	t.Run("panic", func(tt *testing.T) {
		assert.PanicsWithValue(tt,
			"nil error passed, check log for details",
			func() {
				_ = userFailed(&logger, nil)
			},
		)
	})
}
