package tokens

import "github.com/Onnywrite/ssonny/internal/domain/models"

type Pair struct {
	Access  AccessString
	Refresh RefreshString
}

func NewPair(usr *models.User, rotation uint64) (Pair, error) {
	access := Access{
		UserId: usr.Id.String(),
		Email:  usr.Email,
	}
	refresh := Refresh{
		UserId:   usr.Id.String(),
		Rotation: rotation,
	}

	accessStr, err := access.Sign()
	if err != nil {
		return Pair{}, err
	}

	refreshStr, err := refresh.Sign()
	if err != nil {
		return Pair{}, err
	}

	return Pair{
		Access:  accessStr,
		Refresh: refreshStr,
	}, nil
}
