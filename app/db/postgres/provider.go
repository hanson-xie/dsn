package postgres

import (
	"encoding/json"
	"fmt"
	"github.com/Bedrock-Technology/Dsn/app/db"
	"github.com/Bedrock-Technology/Dsn/log"
)

type Provider struct {
}

func (r *Provider) ExecCmd(dsn, sqlCmd string) ([]map[string]interface{}, error) {
	db, err := db.GetDBConnection(dsn)
	if err != nil {
		return nil, err
	}
	log.Debugf("ExecCmd", "sqlCmd:", sqlCmd)
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	rows, err := sqlDB.Query(sqlCmd)
	if err != nil {
		log.Infof("ExecCmd", "Failed to execute SQL:", err)
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	columns, _ := rows.Columns()

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)

		row := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			row[col] = v
		}

		results = append(results, row)
	}

	log.Debugf("ExecCmd", "query result is: %s", func() string {
		jsonData, err := json.Marshal(results)
		if err != nil {
			return fmt.Sprintf("Failed to convert results to JSON: %v", err)
		}
		return string(jsonData)
	}())
	return results, nil
}
