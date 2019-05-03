//resultを返すべきか
//[]Shiftをshifts.idでorderすべきか

package model

import (
	"database/sql"
	"strings"
	"time"

	. "github.com/ahmetb/go-linq"
	"github.com/jmoiron/sqlx"
)

//linqでOrderByするためCompareToメソッドを実装
type ComparableTime struct {
	time.Time
}

func (me ComparableTime) CompareTo(other Comparable) int {
	switch {
	case me.After(other.(ComparableTime).Time):
		return 1
	case me.Before(other.(ComparableTime).Time):
		return -1
	default:
		return 0
	}
}

//Get時に一時的に格納する構造体
type joinedData struct {
	ID      int64
	Name    string
	Begin   string
	End     string
	Created *time.Time
}

type Shift struct {
	Name  string `json:"name"`
	Begin string `json:"begin"`
	End   string `json:"end"`
}

//Post時はID指定の必要なし
type HistoryInfo struct {
	ID      int64      `json:"id"`
	Created *time.Time `json:"created"`
}

type History struct {
	Info   HistoryInfo `json:"info"`
	Shifts []Shift     `json:"shifts"`
}

const aveNumOfCast = 8

func moldJoinedDataToHistory(joined []joinedData) []History {
	historyInfoOrderedByID := make([]HistoryInfo, 0, len(joined)/aveNumOfCast)
	From(joined).DistinctBy(func(j interface{}) interface{} {
		return j.(joinedData).ID
	}).OrderBy(func(j interface{}) interface{} {
		return j.(joinedData).ID
	}).Select(func(j interface{}) interface{} {
		return HistoryInfo{
			j.(joinedData).ID,
			j.(joinedData).Created,
		}
	}).ToSlice(&historyInfoOrderedByID)

	//キャストの順番も考慮する必要があるかも(shifts.idの順にする?)
	shiftOrderedByID := make([]interface{}, 0, len(historyInfoOrderedByID))
	From(joined).GroupBy(func(j interface{}) interface{} {
		return j.(joinedData).ID
	}, func(j interface{}) interface{} {
		return Shift{
			j.(joinedData).Name,
			j.(joinedData).Begin,
			j.(joinedData).End,
		}
	}).OrderBy(func(group interface{}) interface{} {
		return group.(Group).Key
	}).Select(func(group interface{}) interface{} {
		return group.(Group).Group
	}).ToSlice(&shiftOrderedByID)

	histories := make([]History, 0, len(historyInfoOrderedByID))
	From(historyInfoOrderedByID).Zip(From(shiftOrderedByID), func(c, s interface{}) interface{} {
		//shiftOrderedByIDは[][]interface{}のスライスなので[][]Shiftにキャストできない
		//[]interface{}の[]Shiftへのキャストならできる
		shifts := make([]Shift, 0, 10)
		From(s).Select(func(si interface{}) interface{} {
			return si.(Shift)
		}).ToSlice(&shifts)

		return History{
			c.(HistoryInfo),
			shifts,
		}
	}).OrderByDescending(func(h interface{}) interface{} {
		return ComparableTime{*h.(History).Info.Created}
	}).ToSlice(&histories)

	return histories
}

func AllHistories(dbx *sqlx.DB) ([]History, error) {
	var joined []joinedData
	err := dbx.Select(&joined, `
	SELECT histories.id, name, begin, end, created
	FROM shifts INNER JOIN histories
	ON histories.id = shifts.history_id`)

	if err != nil {
		return nil, err
	}

	return moldJoinedDataToHistory(joined), nil
}

func HistoriesWithNum(dbx *sqlx.DB, num int) ([]History, error) {
	var joined []joinedData
	err := dbx.Select(&joined, `
	SELECT limited.id, name, begin, end, created
	FROM shifts INNER JOIN
	(SELECT * FROM histories ORDER BY created DESC LIMIT ?) AS limited
	ON limited.id = shifts.history_id`, num)

	if err != nil {
		return nil, err
	}

	return moldJoinedDataToHistory(joined), nil
}

func (history *History) Post(tx *sqlx.Tx) (res sql.Result, err error) {
	stmtHistory, err := tx.Prepare("INSERT INTO histories (created) VALUES (?)")
	if err != nil {
		return
	}
	defer stmtHistory.Close()
	res, err = stmtHistory.Exec(history.Info.Created)
	if err != nil {
		return
	}
	historyID, err := res.LastInsertId()
	if err != nil {
		return
	}

	queryStr := "INSERT INTO shifts (history_id,name,begin,end) VALUES "
	expand := func(shift Shift) []interface{} {
		return []interface{}{shift.Name, shift.Begin, shift.End}
	}
	expandedList := make([]interface{}, 0, 4*len(history.Shifts))
	for _, shift := range history.Shifts {
		queryStr += "(?,?,?,?),"
		//一つのappendにまとめられない スライスの展開の仕様を確認すべきかも
		expandedList = append(expandedList, historyID)
		expandedList = append(expandedList, expand(shift)...)
	}
	stmtShift, err := tx.Prepare(strings.TrimRight(queryStr, ","))
	if err != nil {
		return
	}
	defer stmtShift.Close()
	_, err = stmtShift.Exec(expandedList...)
	if err != nil {
		return
	}
	return
}

func (history *History) Delete(tx *sqlx.Tx) (res sql.Result, err error) {
	stmt, err := tx.Prepare("DELETE FROM histories WHERE histories.id = ?")
	if err != nil {
		return
	}
	defer stmt.Close()
	res, err = stmt.Exec(history.Info.ID)
	return
}

//使わないかも
func DeleteHistoriesIn(tx *sqlx.Tx, IDs []int64) (err error) {
	queryStr := "DELETE FROM histories WHERE id IN ("
	//queryStrの可変部分を組み立て
	//プレースホルダへのバインドのため[]int64を[]interface{}に変換
	bindings := make([]interface{}, 0, len(IDs))
	for _, id := range IDs {
		queryStr += "?,"
		bindings = append(bindings, id)
	}
	queryStr = strings.TrimRight(queryStr, ",")
	queryStr += ")"

	stmt, err := tx.Prepare(queryStr)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(bindings...)
	return
}
