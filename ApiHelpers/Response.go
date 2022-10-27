package ApiHelpers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type token struct {
	Access	string
	Refresh	string
}

func RespondJSON(w *gin.Context, status int, payload interface{}) {
	_, ok := payload.(map[string]string)
	if ok {
		var err error
		p, err := json.Marshal(payload)
		if err != nil {
			status = 400
			p = []byte("ERROR")
		} 
		w.Data(status, "application/json", p)
		return
	}
	
	w.JSON(status, payload)
}
