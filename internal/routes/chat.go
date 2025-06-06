package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jdpolicano/go-message-board/internal/controller"
	"github.com/jdpolicano/go-message-board/internal/util"
)

type BoardListResonse struct {
	sessions []string `json:"boards"`
}

type BoardCreateRequest struct {
	user string `json:"user"`
}

func NewListChatHandler(c *controller.Controller) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sessions := c.ListSessionIds()
		payload, e := json.Marshal(BoardListResonse{sessions})
		if e != nil {
			util.Error(w, fmt.Sprintf("encoding payload: %s", e), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(payload)
	}
}
