package limitProvider

import "github.com/pkg/errors"

func GetToken(model int, value int) (int, error) {
	switch model {
	case 1: //
		return 3, nil
	case 2:
		return 3, nil
	default:
		return -1, errors.New("model value error")
	}
}
