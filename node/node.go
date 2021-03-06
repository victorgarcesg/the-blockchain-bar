package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"the-blockchain-bar/database"
)

const httpPort = 8080

type ErrRes struct {
	Error string `json:"error"`
}

type TxAddReq struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value uint   `json:"value"`
	Data  string `json:"data"`
}

type TxAddRes struct {
	Hash database.Hash `json:"block_hash"`
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

	http.HandleFunc("/tx/add", func(rw http.ResponseWriter, r *http.Request) {
		txAddHandler(rw, r, state)
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
}

func listBalancesHandler(rw http.ResponseWriter, r *http.Request, state *database.State) {
	writeRes(rw, &BalancesRes{state.LatestSnapshot(), state.Balances})
}

func txAddHandler(rw http.ResponseWriter, r *http.Request, state *database.State) {
	req := TxAddReq{}
	err := readReq(r, &req)
	if err != nil {
		writeErrRes(rw, err)
		return
	}

	tx := database.NewTx(database.NewAccount(req.From), database.NewAccount(req.To), req.Value, req.Data)

	err = state.AddTx(tx)
	if err != nil {
		writeErrRes(rw, err)
		return
	}

	hash, err := state.Persist()
	if err != nil {
		writeErrRes(rw, err)
		return
	}

	writeRes(rw, TxAddRes{hash})
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

func readReq(r *http.Request, reqBody interface{}) error {
	reqBodyJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	err = json.Unmarshal(reqBodyJson, reqBody)
	if err != nil {
		return fmt.Errorf("unable to unmarshal request body %s", err.Error())
	}

	return nil
}
