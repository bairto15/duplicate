package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)


func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user_id1_str := ps.ByName("user_id1")
	user_id2_str := ps.ByName("user_id2")

	user_id1, err := strconv.ParseInt(user_id1_str, 10, 32)
	if err != nil {
		h.logger.Error(err)
		fmt.Fprint(w, err.Error())
		return
	}

	user_id2, err := strconv.ParseInt(user_id2_str, 10, 32)
	if err != nil {
		h.logger.Error(err)
		fmt.Fprint(w, err.Error())
		return
	}

	dupes := h.service.IsDuplicate(int32(user_id1), int32(user_id2))

	result := result{dupes}

	response, err := json.Marshal(result)
	if err != nil {
		h.logger.Error(err)
		fmt.Fprint(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}