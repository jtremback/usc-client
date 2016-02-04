package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// func TestNewEscrowProvider(t *testing.T) {
// 	// handler := new(EchoHandler)
// 	expectedBody := "Hello"

// 	rw := httptest.NewRecorder()

// 	req, err := http.NewRequest(
// 		"POST",
// 		fmt.Sprintf("http://example.com/"),
// 		strings.NewReader(`{"name": "hokey","address": "crap", "pubkey": "ZnVjaw==", "privkey": "ZnVjaw=="}`),
// 	)

// 	if err != nil {
// 		t.Errorf("Failed to create request.")
// 	}

// 	addEscrowProvider(rw, req)

// 	switch rw.Body.String() {
// 	case expectedBody:
// 		// body is equal so no need to do anything
// 	default:
// 		t.Errorf("Body (%s) did not match expectation (%s).",
// 			rw.Body.String(),
// 			expectedBody)
// 	}
// }

func TestShouldGetPosts(t *testing.T) {
	db, err := sql.Open("sqlite3", "/tmp/foo.db")
	if err != nil {
		t.Fatal(err)
	}

	// create app with mocked db, request and response to test
	app := &api{db}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS channels (
		ChannelId 							 TEXT PRIMARY KEY,
		Phase										 INT,

		OpeningTx         			 BLOB,
		OpeningTxEnvelope 			 BLOB,

		LastUpdateTx         		 BLOB,
		LastUpdateTxEnvelope 		 BLOB,

		LastFullUpdateTx         BLOB,
		LastFullUpdateTxEnvelope BLOB,

		EscrowProvider 					 TEXT,
		Accounts 								 TEXT,

		Me           						 INT,
		Fulfillments 						 TEXT
	);`)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Exec(`DROP TABLE channels;`)

	_, err = db.Exec(`INSERT INTO channels values(?,?,?,?,?,?,?,?,?,?,?,?)`,
		"ernel",
		4,
		[]byte{110, 100, 22},
		[]byte{110, 100, 22},
		[]byte{110, 100, 22},
		[]byte{110, 100, 22},
		[]byte{110, 100, 22},
		[]byte{110, 100, 22},
		"erp",
		"fop",
		1,
		"english")
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("http://example.com/"),
		strings.NewReader(`{"name": "hokey","address": "crap", "pubkey": "ZnVjaw==", "privkey": "ZnVjaw=="}`),
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("bingo")
	w := httptest.NewRecorder()

	app.getChannels(w, req)
	fmt.Println(w)
	if w.Code != 200 {
		t.Fatalf("expected status code to be 200, but got: %d", w.Code)
	}
}
