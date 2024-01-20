package api

import (
	"net/http"

	"github.com/pkg/errors"

	"rollbringer/pkg/repositories/database"
)

var (
	errUnauthorized = errors.New("unauthorized")
)

type apiError struct {
	statusCode int

	msg string
}

func (e *apiError) Error() string {
	return e.msg
}

func (api *API) err(w http.ResponseWriter, e error, statusCode int) {
	if statusCode >= 500 {
		api.Logger.Error().Stack().Err(e).Msg("Internal server error")
		e = errors.New("internal server error")
	}

	http.Error(w, e.Error(), statusCode)
}

func (api *API) dbErr(w http.ResponseWriter, err error) {

	cause := errors.Cause(err)
	res := &apiError{msg: cause.Error()}

	switch cause {
	case database.ErrUnauthorized:
		res.statusCode = http.StatusUnauthorized

	case database.ErrUserNotFound:
		res.statusCode = http.StatusNotFound

	case database.ErrGameNotFound:
		res.statusCode = http.StatusNotFound
	case database.ErrMaxGames:
		res.statusCode = http.StatusForbidden

	default:
		res.statusCode = http.StatusInternalServerError
		res.msg = err.Error()
	}

	api.err(w, res, res.statusCode)
}
