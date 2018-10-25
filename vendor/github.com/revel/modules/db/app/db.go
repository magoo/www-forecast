// Copyright (c) 2012-2016 The Revel Framework Authors, All rights reserved.
// Revel Framework source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

// Package db module configures a database connection for the application.
//
// Developers use this module by importing and calling db.Init().
// A "Transactional" controller type is provided as a way to import interceptors
// that manage the transaction
//
// In particular, a transaction is begun before each request and committed on
// success.  If a panic occurred during the request, the transaction is rolled
// back.  (The application may also roll the transaction back itself.)
package db

import (
	"database/sql"

	"github.com/revel/revel"
)

// Database connection variables
var (
	Db     *sql.DB
	Driver string
	Spec   string
)

// Init method used to initialize DB module on `OnAppStart`
func Init() {
	// Read configuration.
	var found bool
	if Driver, found = revel.Config.String("db.driver"); !found {
		revel.RevelLog.Fatal("db.driver not configured")
	}
	if Spec, found = revel.Config.String("db.spec"); !found {
		revel.RevelLog.Fatal("db.spec not configured")
	}

	// Open a connection.
	var err error
	Db, err = sql.Open(Driver, Spec)
	if err != nil {
		revel.RevelLog.Fatal("Open database connection error", "error", err, "driver", Driver, "spec", Spec)
	}
}

// Transactional definition for database transaction
type Transactional struct {
	*revel.Controller
	Txn *sql.Tx
}

// Begin a transaction
func (c *Transactional) Begin() revel.Result {
	txn, err := Db.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	return nil
}

// Rollback if it's still going (must have panicked).
func (c *Transactional) Rollback() revel.Result {
	if c.Txn != nil {
		if err := c.Txn.Rollback(); err != nil {
			if err != sql.ErrTxDone {
				panic(err)
			}
		}
		c.Txn = nil
	}
	return nil
}

// Commit the transaction.
func (c *Transactional) Commit() revel.Result {
	if c.Txn != nil {
		if err := c.Txn.Commit(); err != nil {
			if err != sql.ErrTxDone {
				panic(err)
			}
		}
		c.Txn = nil
	}
	return nil
}

func init() {
	revel.InterceptMethod((*Transactional).Begin, revel.BEFORE)
	revel.InterceptMethod((*Transactional).Commit, revel.AFTER)
	revel.InterceptMethod((*Transactional).Rollback, revel.FINALLY)
}
