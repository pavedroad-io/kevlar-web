// prTokenModel
// data structure and db/message interfaces

package main

import (
	_ "database/sql"
	_ "fmt"
	"time"
)

type PrToken struct {
	ApiVersion   string    `json:"core.pavedroad.io/v1alpha1"`
	Kind         string    `json:"PrToken"`
    Metadata struct {
	    Name      string    `json:"defauult-name"`
	    Namespace string    `json:"defauult-namespace"`
	    Uid       string    `json:"defauult-name"`
	    Site      string    `json:"github"`
	    EndPoint  string    `json:"api.githubm.com"`
	    Token     string    `json:""`
        Scope     []string
    }
	Created       time.Time
	Updated       time.Time
	Active        bool       `json:"True"`
}
