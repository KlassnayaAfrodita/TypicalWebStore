package handlers

import (
	"GOLANG/net/projext-net/storage"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func (api *Api) Auth(w http.ResponseWriter, r *http.Request) { //! принимаем POST json
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error": "server error"}`, 500)
	}
	defer r.Body.Close()

	var user storage.User

	newerr := json.Unmarshal(body, &user)
	if newerr != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	trueUser, err := storage.User.GetUserByName(user.Name)
	if err != nil {
		http.Error(w, `{"error":"not found"}`, 404)
	}
	if trueUser.Password != user.Password {
		http.Error(w, `bad pass`, 400)
		return
	}

	SID, err := api.session.SetSession(trueUser.ID)
	if err != nil {
		http.Error(w, `{"error":"db error"}`, 500)
	}
	cookie := http.Cookie{
		Name:    "session_id",
		Value:   SID,
		Expires: time.Now().Add(10 * time.Hour),
	}

	http.SetCookie(w, &cookie)
	w.Write([]byte(SID))
}
