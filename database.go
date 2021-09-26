package main

import (
	"time"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
	stmtInsertHouerData *sql.Stmt
	stmtInsertFilNodes *sql.Stmt
)

func initDB() {
	var err error
	db, err = sql.Open("mysql", config.user + ":" + config.password + "@/" + config.name)
	if err != nil {
		log.Fatal(err)
	}

	stmtInsertHouerData, err = db.Prepare("insert into hour_data values(null, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	stmtInsertFilNodes, err = db.Prepare("insert into fil_node values(null, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

}

func insertHourData(lowcaseB, countDrawnsFil float64) (int64, error) {
	createTime := time.Now()
	result, err := stmtInsertHouerData.Exec(createTime, lowcaseB, countDrawnsFil)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func insertFilNodes(id int64, filNodes map[string]cacheFilNode_T) error {
	for k, filNode := range filNodes {
		_, err := stmtInsertFilNodes.Exec(k, filNode.Address, filNode.Balance, filNode.WorkerBalance, filNode.QualityAdjPower, filNode.AvailableBalance, filNode.Pledge, filNode.VestingFunds, filNode.SingleT, id)
		if err != nil {
			return err
		}
	}

	return nil
}
