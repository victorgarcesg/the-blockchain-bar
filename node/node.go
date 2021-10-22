package node

import (
	"encoding/json"
	"fmt"
	"net/http"
	"the-blockchain-bar/database"
)

const httpPort = 8080

type ErrRes struct {
	Error string `json:"error"`
}

type BalancesRes struct {
	Hash     database.Hash             `json:"hash"`
	Balances map[database.Account]uint `json:"balances"`
}

func Run(dataDir string) error {
	state, err := database.NewStateFromDisk(dataDir)
	if err != nil {
		return err
	}
	defer state.Close()

	http.HandleFunc("/balances/list", func(rw http.ResponseWriter, r *http.Request) {
		listBalancesHandler(rw, r, state)
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
}

func listBalancesHandler(rw http.ResponseWriter, r *http.Request, state *database.State) {
	writeRes(rw, &BalancesRes{state.LatestSnapshot(), state.Balances})
}

func writeRes(rw http.ResponseWriter, content interface{}) {
	contentJson, err := json.Marshal(content)
	if err != nil {
		writeErrRes(rw, err)
		return
	}

	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(contentJson)
}

func writeErrRes(rw http.ResponseWriter, err error) {
	jsonErrRes, _ := json.Marshal(ErrRes{err.Error()})
	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(http.StatusInternalServerError)
	rw.Write(jsonErrRes)
}
