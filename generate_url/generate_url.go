package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"

	"github.com/alufers/owm-patrons/common"
	"github.com/kevinburke/nacl"
	"github.com/kevinburke/nacl/secretbox"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	key, err := nacl.Load(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		common.WriteJSON(w, 500, "Failed to create private key")
		return
	}
	data, err := json.Marshal(&common.J{})

	encrypted := secretbox.EasySeal(data, key)

	base64Encoded := base64.StdEncoding.EncodeToString(encrypted)

	common.WriteJSON(w, 200, &common.J{"url": os.Getenv("NOW_URL") + "/new-patron/" + base64Encoded})
}
