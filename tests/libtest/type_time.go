package libtest

import (
	"database/sql"

	"testing"

	"time"
)

// DoTestTime tests the handling of the Time.
func DoTestTime(t *testing.T) {
	TestForEachDB("TestTime", t, testTime)
	//
}

func testTime(t *testing.T, db *sql.DB, tableName string) {
	pass := make([]interface{}, len(samplesTime))
	mySamples := make([]time.Time, len(samplesTime))

	for i, sample := range samplesTime {

		mySample := sample

		pass[i] = mySample
		mySamples[i] = mySample
	}

	rows, err := SetupTableInsert(db, tableName, "time", pass...)
	if err != nil {
		t.Errorf("Error preparing table: %v", err)
		return
	}
	defer rows.Close()

	i := 0
	var recv time.Time
	for rows.Next() {
		err = rows.Scan(&recv)
		if err != nil {
			t.Errorf("Scan failed on %dth scan: %v", i, err)
			continue
		}

		if recv != mySamples[i] {

			t.Errorf("Received value does not match passed parameter")
			t.Errorf("Expected: %v", mySamples[i])
			t.Errorf("Received: %v", recv)
		}

		i++
	}

	if err := rows.Err(); err != nil {
		t.Errorf("Error preparing rows: %v", err)
	}
}
