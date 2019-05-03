package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/yamadashi/EscaTSGen/model"
)

type History struct {
	DB     *sqlx.DB
	numget int
}

func (h *History) GetAll(w http.ResponseWriter, r *http.Request) error {
	histories, err := model.AllHistories(h.DB)
	if err != nil {
		return err
	}
	return JSON(w, 200, histories)
}

func (h *History) Get(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	numStr, ok := params["num"]
	var num int
	if !ok {
		num = 0
	}
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return err
	}
	histories, err := model.HistoriesWithNum(h.DB, num)
	if err != nil {
		return err
	}
	return JSON(w, 200, histories)
}

func (h *History) PostOne(w http.ResponseWriter, r *http.Request) error {
	var history model.History
	if err := json.NewDecoder(r.Body).Decode(&history); err != nil {
		return err
	}
	if err := TXHandler(h.DB, func(tx *sqlx.Tx) error {
		res, err := history.Post(tx)
		if err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		history.Info.ID, err = res.LastInsertId()
		return err
	}); err != nil {
		return err
	}
	return JSON(w, http.StatusCreated, history)
}

func (h *History) Delete(w http.ResponseWriter, r *http.Request) error {
	var history model.History
	if err := json.NewDecoder(r.Body).Decode(&history); err != nil {
		return err
	}
	if err := TXHandler(h.DB, func(tx *sqlx.Tx) error {
		_, err := history.Delete(tx)
		if err != nil {
			return err
		}
		return tx.Commit()
	}); err != nil {
		return err
	}
	return JSON(w, http.StatusOK, history)
}
