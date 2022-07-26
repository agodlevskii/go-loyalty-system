package order

import (
	"go-loyalty-system/user"
	"io"
	"net/http"
	"time"
)

var errToStat = map[string]int{
	ErrSameUser:  http.StatusOK,
	ErrOtherUser: http.StatusConflict,
}

func GetOrders(db Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func UpdateOrders(db Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		usr, ok := r.Context().Value(user.Key).(string)
		if !ok || usr == `` {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		id, err := io.ReadAll(r.Body)
		if err != nil || id == nil || len(id) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !validateOrderNumber(string(id)) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		_, err = db.Add(Order{
			Number:     string(id),
			Status:     StatusNew,
			Accrual:    0,
			UploadedAt: time.Now().Round(time.Microsecond),
			User:       usr,
		})

		if err != nil {
			if code, ok := errToStat[err.Error()]; ok {
				w.WriteHeader(code)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
