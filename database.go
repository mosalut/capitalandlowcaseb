package main

import (
	"strconv"
	"time"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
	stmtInsertHourData *sql.Stmt
	stmtInsert5MinsData *sql.Stmt
	stmtInsertFilNodes *sql.Stmt
	stmtQueryHourData *sql.Stmt
	stmtQueryCurve *sql.Stmt
	stmtQuery5MinsData *sql.Stmt
	stmtQueryFilNodes *sql.Stmt
)

func initDB() {
	var err error
	db, err = sql.Open("mysql", config.user + ":" + config.password + "@/" + config.name)
	if err != nil {
		log.Fatal(err)
	}

	stmtInsertHourData, err = db.Prepare("insert into hour_data values(null, from_unixtime(unix_timestamp(current_timestamp) div ? * ?), ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	stmtInsert5MinsData, err = db.Prepare("insert into 5_mins_data values(null, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	stmtInsertFilNodes, err = db.Prepare("insert into fil_node values(null, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	stmtQueryHourData, err = db.Prepare("select id, lowcase_b, drawn_fil from hour_data order by create_time desc")
	if err != nil {
		log.Fatal(err)
	}

	stmtQueryCurve, err = db.Prepare("select * from hour_data order by create_time desc limit 24")
	if err != nil {
		log.Fatal(err)
	}

	stmtQuery5MinsData, err = db.Prepare("select create_time, cfil_to_fil from 5_mins_data order by create_time desc limit 288")
	if err != nil {
		log.Fatal(err)
	}

	stmtQueryFilNodes, err = db.Prepare("select node_name, address, balance, worker_balance, quality_adj_power, available_balance, pledge, vestingFunds, singletT from fil_node, hour_data where hour_data_id = ?")
	if err != nil {
		log.Fatal(err)
	}
}

func insertHourData(lowcaseB, countDrawnsFil float64) (int64, error) {
	period := strconv.FormatInt(config.period, 10)
	result, err := stmtInsertHourData.Exec(period, period, lowcaseB, countDrawnsFil)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func insertFilNodes(id int64, filNodes map[string]filNode_T) error {
	for k, filNode := range filNodes {
		_, err := stmtInsertFilNodes.Exec(k, filNode.Address, filNode.Balance, filNode.WorkerBalance, filNode.QualityAdjPower, filNode.AvailableBalance, filNode.Pledge, filNode.VestingFunds, filNode.SingleT, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func initCacheData() {
	var id int
	err := stmtQueryHourData.QueryRow().Scan(&id, &cache.LowcaseB, &cache.DrawnFil)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := stmtQueryFilNodes.Query(id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		var address string
		var balance float64
		var workerBalance float64
		var qualityAdjPower float64
		var availableBalance float64
		var pledge float64
		var vestingFunds float64
		var singletT float64
		err := rows.Scan(&key, &address, &balance, &workerBalance, &qualityAdjPower, &availableBalance, &pledge, &vestingFunds, &singletT)
		if err != nil {
			log.Fatal(err)
		}
		cache.FilNodes[key] = filNode_T{address, balance, workerBalance, qualityAdjPower, availableBalance, pledge, vestingFunds, singletT}
	}

	rows5, err := stmtQuery5MinsData.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows5.Close()
	for rows5.Next() {
	}
}

func getCurveData() ([]data24_T, []data24_T, []data24_T, []data24_T, error) {
	rows, err := stmtQueryCurve.Query()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer rows.Close()

	lowcaseBs := make([]data24_T, 0, 24)
	var timeStr string
	for rows.Next() {
		lowcaseB := data24_T{}
		var id int
		err := rows.Scan(&id, &timeStr, &lowcaseB.Value, &cache.DrawnFil)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		lowcaseB.CreateTime, err = time.Parse("2006-01-02 15:04:05", timeStr)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		lowcaseBs = append([]data24_T{lowcaseB}, lowcaseBs...)
	}

	worthDeposits, err := getData24()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	filDrawns, err := getData24()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	cfilDrawns := getCfilDrawnsData()

	rows5, err := stmtQuery5MinsData.Query()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer rows5.Close()
	for rows5.Next() {
	}

	return lowcaseBs, worthDeposits, filDrawns, cfilDrawns, nil
}
