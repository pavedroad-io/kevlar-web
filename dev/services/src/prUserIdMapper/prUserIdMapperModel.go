// prUserIdMapperModel
// data structure and db/message interfaces

package main

import (
	"database/sql"
	_ "encoding/json"
	_ "errors"
	"fmt"
	 "log"
   "errors"
  _  "os"
)

// prUserIdMapper data structure for token storage
type prUserIdMapper struct {
	APIVersion string `json:"apiVersion"`
  ObjVersion string `json:"objVersion"`
	Kind       string `json:"kind"`
  Credential string `json:"login"`
  UserUUID string `json:"userUUID"`
  LoginCount int `json:"loginCount"`
	Created string `json:"created,ignoreempty"`
	Updated string `json:"updated,ignoreempty"`
	Active  string `json:"active"`
}

// updateUserIdMapper in database
func (t *prUserIdMapper) updateUserIdMapper(db *sql.DB) error {
  update := `
	UPDATE pavedroad.pruseridmapper
  SET apiVersion = '%s', 
      objVersion = '%s',
      kind = '%s', 
      userUUID = '%s', 
      loginCount = '%d', 
      updated = '%s', 
      active = '%s'
  WHERE
      credential = '%s';
  `
	statement := fmt.Sprintf(update, t.APIVersion, t.ObjVersion, t.Kind, 
  t.UserUUID, t.LoginCount, t.Updated, t.Active, t.Credential)
  // fmt.Println(statement)

	_, er1 := db.Query(statement)

	if er1 != nil {
		log.Println("Update failed")
		return er1
	}

	return nil
}

// createUserIdMapper in database
func (t *prUserIdMapper) createUserIdMapper(db *sql.DB) error {

	statement := fmt.Sprintf("INSERT INTO pavedroad.pruseridmapper(apiVersion, objVersion, kind, credential, userUUID, loginCount, created, updated, active) VALUES('%s', '%s', '%s', '%s', '%s', '%d', '%s', '%s', '%s') RETURNING credential;", 
    t.APIVersion, t.ObjVersion, t.Kind, t.Credential, t.UserUUID, t.LoginCount,
    t.Created, t.Updated, t.Active)
  //fmt.Println(statement)

	rows, er1 := db.Query(statement)

	if er1 != nil {
    log.Printf("Insert failed for: %s", t.Credential)
    log.Printf("SQL Error: %s", er1)
		return er1
	}

	defer rows.Close()

	return nil
}

// getUserIdMappers: return a list of tokens
//
func (t *prUserIdMapper) getUserIdMappers(db *sql.DB, start, count int) ([]prUserIdMapper, error) {
	statement := fmt.Sprintf("SELECT apiVersion, objVersion, kind, credential, userUUID, loginCount, created, updated, active FROM pavedroad.pruseridmapper LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	mappings := []prUserIdMapper{}

	for rows.Next() {
		var t prUserIdMapper
		err := rows.Scan(&t.APIVersion, &t.ObjVersion, &t.Kind, &t.Credential, 
    &t.UserUUID, &t.LoginCount, &t.Created, &t.Updated, &t.Active)

	  if err != nil {
      log.Printf("SQL rows.Scan failed: %s", err)
  		return mappings, err
  	}

		mappings = append(mappings, t)
	}

	return mappings, nil
}

// getUserIdMapper: return a token based on UID
//
func (t *prUserIdMapper) getUserIdMapper(db *sql.DB, key string) error {
	statement := fmt.Sprintf("SELECT apiVersion, objVersion, kind, credential, userUUID, loginCount, created, updated, active FROM pavedroad.pruseridmapper WHERE credential = '%s'", key)
	row := db.QueryRow(statement)

  // Fill in mapper
  switch err := row.Scan(&t.APIVersion, &t.ObjVersion, &t.Kind, &t.Credential, 
  &t.UserUUID, &t.LoginCount, &t.Created, &t.Updated, &t.Active); err {

  case sql.ErrNoRows:
    m := fmt.Sprintf("credential %s does not exist", key)
	  return errors.New(m)
  case nil:
    break;
  default:
       //Some error to catch
       panic(err)
     }
	return nil
}

// deleteUserIdMapper: return a token based on UID
//
func (t *prUserIdMapper) deleteUserIdMapper(db *sql.DB, cred string) error {
	statement := fmt.Sprintf("DELETE FROM pavedroad.pruseridmapper WHERE credential = '%s'", cred)
	result, err := db.Exec(statement)
	c, e := result.RowsAffected()

	if e == nil && c == 0 {
    em := fmt.Sprintf("credential %s does not exist", cred)
    log.Println(em)
    log.Println(e)
		return errors.New(em)
	}

	return err
}
