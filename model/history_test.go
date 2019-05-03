package model

import (
	"testing"
	"time"

	. "github.com/ahmetb/go-linq"
	"github.com/jmoiron/sqlx"
	"github.com/yamadashi/EscaTSGen/db"
)

func DefaultDB() *sqlx.DB {
	configs, err := db.NewConfigsFromFile("../dbconfig.yml")
	if err != nil {
		panic(err)
	}
	dbx, err := configs.Open("test")
	if err != nil {
		panic(err)
	}

	return dbx
}

//timeのポインタを取るかどうかは考える余地がある
var now = time.Now()
var tommorow = now.Add(24 * time.Hour)
var yesterday = now.Add(-24 * time.Hour)
var historiesMock = []History{
	History{
		HistoryInfo{
			0,
			&now,
		},
		[]Shift{
			Shift{"ui", "12:00", "17:00"},
			Shift{"rin", "17:00", "21:00"},
		},
	},
	History{
		HistoryInfo{
			0,
			&tommorow,
		},
		[]Shift{
			Shift{"tsubasa", "12:00", "16:00"},
			Shift{"yuzu", "19:30", "23:00"},
		},
	},
	History{
		HistoryInfo{
			0,
			&yesterday,
		},
		[]Shift{
			Shift{"remu", "18:00", "23:00"},
		},
	},
}

func TestPost(t *testing.T) {
	dbx := DefaultDB()
	defer dbx.Close()

	//PostHistory Test
	for _, history := range historiesMock {
		tx := dbx.MustBegin()
		defer func() {
			if err := recover(); err != nil {
				tx.Rollback()
			}
		}()
		if err := PostHistory(tx, history); err != nil {
			t.Fatalf("PostHistory error : %s", err)
		}
		tx.Commit()
	}
}

func TestGet(t *testing.T) {
	dbx := DefaultDB()
	defer dbx.Close()

	//AllHistories Test
	allHistories, err := AllHistories(dbx)
	if err != nil {
		t.Fatalf("AllHistories error : %s", err)
	} else if len(allHistories) != len(historiesMock) {
		t.Fatalf("len(histories) want %d got %d", len(historiesMock), len(allHistories))
	}

	//HisotriesWithNum Test
	var num int = 2
	histories, err := HistoriesWithNum(dbx, num)
	if err != nil {
		t.Fatalf("HistoriesWithNum error : %s", err)
	} else if len(histories) != num {
		t.Fatalf("len(histories) want %d got %d", num, len(histories))
	}
}

func TestDelete(t *testing.T) {
	dbx := DefaultDB()
	defer dbx.Close()

	allHistories, err := AllHistories(dbx)
	if err != nil {
		t.Fatalf("AllHistories error : %s", err)
	}

	//DeleteHistory
	tx2 := dbx.MustBegin()
	if err := DeleteHistory(tx2, allHistories[0].Info.ID); err != nil {
		t.Fatalf("DeleteHistory error : %s", err)
	}
	tx2.Commit()

	histories, err := AllHistories(dbx)
	if err != nil {
		t.Fatalf("AllHistories error : %s", err)
	} else if len(histories) != len(historiesMock)-1 {
		t.Fatalf("len(histories) want %d got %d", len(historiesMock)-1, len(histories))
	}

	//DeleteHistoriesIn
	IDs := make([]int64, 0, len(histories))
	From(histories).Select(func(h interface{}) interface{} {
		return h.(History).Info.ID
	}).ToSlice(&IDs)
	tx3 := dbx.MustBegin()
	if err = DeleteHistoriesIn(tx3, IDs); err != nil {
		t.Fatalf("DeleteHistory error : %s", err)
	}
	tx3.Commit()
	histories, err = AllHistories(dbx)
	if err != nil {
		t.Fatalf("AllHistories error : %s", err)
	} else if len(histories) != 0 {
		t.Fatalf("len(hsitories) want %d got %d", 0, len(histories))
	}
}
