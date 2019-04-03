// prTokenModel
// data structure and db/message interfaces

package main

import (
	"fmt"
	"database/sql"
	"encoding/json"
	"errors"
	"os"
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

// updateToken in database
// Store the UID created as a database key in the prtoken struct
func (t *PrToken) updateToken(db *sql.DB) error {
        jb, err := json.Marshal(t)
        if err != nil {
		fmt.Println("marshall failed")
                panic(err)
        }

	statement := fmt.Sprintf("UPDATE kevlarweb.prtoken SET prtoken='%s' WHERE uid ='%s'",jb, t.Metadata.UID)
	//fmt.Println(statement)
	_, er1 := db.Query(statement)

	if er1 != nil {
		//TODO(jscharber): convert to log message
		fmt.Println("Update failed")
		return er1
	}

	return nil
}
// createToken in database
// Store the UID created as a database key in the prtoken struct
func (t *PrToken) createToken(db *sql.DB) error {
        jb, err := json.Marshal(t)
        if err != nil {
                panic(err)
        }

	statement := fmt.Sprintf("INSERT INTO kevlarweb.prtoken(prtoken) VALUES('%s') RETURNING uid", jb)
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

// getTokens: return a list of tokens
// 
func (t *PrToken) getTokens(db *sql.DB, start, count int) ([]PrToken,error) {
	statement := fmt.Sprintf("SELECT uid, prtoken FROM kevlarweb.prtoken LIMIT %d OFFSET %d", count, start)
	//fmt.Println(statement)
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tokens := []PrToken{}

	for rows.Next() {
		var t PrToken
		var jb []byte
		var uid string
		err := rows.Scan(&uid, &jb)
		err = json.Unmarshal(jb, &t)
		t.Metadata.UID = uid
		//fmt.Println(uid, string(jb))
        	if err != nil {
                	fmt.Println(err)
                	os.Exit(1)
        	}
		tokens = append(tokens, t)
	}

	return tokens, nil
}

// getToken: return a token based on UID
// 
func (t *PrToken) getToken(db *sql.DB, uid string) (error) {
	statement := fmt.Sprintf("SELECT uid, prtoken FROM kevlarweb.prtoken WHERE UID = '%s'", uid)
	rows, err := db.Query(statement)
	//fmt.Println(statement)

	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var jb []byte
		var uid string
		err := rows.Scan(&uid, &jb)
		err = json.Unmarshal(jb, t)
		t.Metadata.UID = uid
		//fmt.Println(uid, string(jb))
        	if err != nil {
                	fmt.Println(err)
                	os.Exit(1)
        	}
	}

	return nil
}

// deleteToken: return a token based on UID
// 
func (t *PrToken) deleteToken(db *sql.DB, uid string) (error) {
	statement := fmt.Sprintf("DELETE FROM kevlarweb.prtoken WHERE UID = '%s'", uid)
	result, err := db.Exec(statement)
	//fmt.Println(statement)
	c, e := result.RowsAffected()

	if e == nil && c == 0 {
		em := "UID: "+uid+ " does not exist"
		return errors.New(em)
	}

	return err
}
