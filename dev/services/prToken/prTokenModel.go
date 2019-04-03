// prTokenModel
// data structure and db/message interfaces

package main

import (
	"fmt"
	"database/sql"
	"encoding/json"
)

type PrToken struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name      string   `json:"name"`
		Namespace string   `json:"namespace"`
		UID       string   `json:"uid"`
		Site      string   `json:"site"`
		EndPoint  string   `json:"endPoint"`
		Token     string   `json:"token"`
		Scope     []string `json:"scope"`
	} `json:"metadata"`
	Created string `json:"created"`
	Updated string `json:"updated"`
	Active  bool   `json:"active"`
}

// createToken in database
// Store the UID created as a database key in the prtoken struct
func (t *PrToken) createToken(db *sql.DB) error {
        jb, err := json.Marshal(t)
        if err != nil {
                panic(err)
        }

	statement := fmt.Sprintf("INSERT INTO prtoken(prtoken) VALUES('%s') RETURNING uid", jb)
	//fmt.Println(statement)
	rows, er1 := db.Query(statement)

	if er1 != nil {
		//TODO(jscharber): convert to log message
		fmt.Println("Insert failed")
		return er1
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&t.Metadata.UID)
		if err != nil {
			return err
		}
		fmt.Println(t.Metadata.UID);	
	}

	return nil
}
