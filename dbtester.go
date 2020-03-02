package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"os/exec"
	//"log"
	"runtime"
	"time"
)

type DbTester struct {
	db *sqlx.DB
}

const schema = `
--
-- DSN: postgres://iktest@localhost/iktest?sslmode=disable
--

CREATE TABLE IF NOT EXISTS ik_table (
    uuid_v1      uuid,
    uuid_v4      uuid,
    snowflake_id bigint,
    ctime TIMESTAMP WITH TIME ZONE default now()
);
CREATE UNIQUE INDEX IF NOT EXISTS uuid_v1 ON ik_table (uuid_v1);
CREATE UNIQUE INDEX IF NOT EXISTS uuid_v4 ON ik_table (uuid_v4);
CREATE UNIQUE INDEX IF NOT EXISTS snowflake ON ik_table (snowflake_id);
`

func NewDbTester(dsn string) *DbTester {
	client := &DbTester{}
	client.db = sqlx.MustConnect("postgres", dsn)
	return client
}

func (d *DbTester) InitDb() {
	d.db.MustExec(schema)
}

// CleanDbCache cleans up cache/buffer owned by db and os.
func (d *DbTester) CleanDbCache() {
	if pauseForCache {
		fmt.Printf("NOTE: Please clear PostgreSQL cache, and press ENTER when done...")
		var s string
		fmt.Scanln(&s)
		return
	}

	var cmdList []string
	switch runtime.GOOS {
	case "linux":
		cmdList = []string{
			"service postgresql stop",
			"sync",
			"echo 3 > /proc/sys/vm/drop_caches",
			"service postgresql start",
		}

	case "darwin":
		cmdList = []string{
			"pg_ctl -D /usr/local/var/postgres stop",
			"sync",
			"sudo purge",
			"pg_ctl -D /usr/local/var/postgres start",
		}
	}

	log.Debug("Clearing PostgreSQL cache and restarting...")
	for _, cmd := range cmdList {
		log.Debug(cmd)
		c := exec.Command("/bin/sh", "-c", cmd)
		if err := c.Run(); err != nil {
			log.Fatal(err)
		}
	}

	log.Debug("Reconnecting to db...")
	for i := 0; i < 100; i++ {
		if err := d.db.Ping(); err == nil {
			break
		}
	}
}

func (d *DbTester) ClearRecords() {
	d.db.MustExec(`DELETE FROM ik_table`)
}

func (d *DbTester) InsertUuidV1(records []Record) {
	defer timeTrack(time.Now(), "InsertUuidV1")
	for _, record := range records {
		d.db.MustExec(`INSERT INTO ik_table (uuid_v1) VALUES ($1)`, record.UuidV1)
	}
}

func (d *DbTester) InsertUuidV4(records []Record) {
	defer timeTrack(time.Now(), "InsertUuidV4")
	for _, record := range records {
		d.db.MustExec(`INSERT INTO ik_table (uuid_v4) VALUES ($1)`, record.UuidV4)
	}
}

func (d *DbTester) InsertSnowflake(records []Record) {
	defer timeTrack(time.Now(), "InsertSnowflake")
	for _, record := range records {
		d.db.MustExec(`INSERT INTO ik_table (snowflake_id) VALUES ($1)`, record.Snowflake)
	}
}

func (d *DbTester) InsertAll(records []Record) {
	defer timeTrack(time.Now(), "InsertAll")
	for _, record := range records {
		d.db.MustExec(`INSERT INTO ik_table (uuid_v1, uuid_v4, snowflake_id) VALUES ($1, $2, $3)`, record.UuidV1, record.UuidV4, record.Snowflake)
	}
}

func (d *DbTester) SelectUuidV1(records []Record, sample []int, withCache bool) {
	name := "SelectUuidV1"
	if withCache {
		name += "/cache"
	} else {
		name += "/clean"
	}
	defer timeTrack(time.Now(), name)

	for _, idx := range sample {
		row := d.db.QueryRowx(`SELECT * FROM ik_table WHERE uuid_v1 = $1`, records[idx].UuidV1)
		var record Record
		if err := row.StructScan(&record); err != nil {
			log.Fatal("SelectUuidV1 fail")
		}
		//fmt.Println(record.ToCSV())
	}
}

func (d *DbTester) SelectUuidV4(records []Record, sample []int, withCache bool) {
	name := "SelectUuidV4"
	if withCache {
		name += "/cache"
	} else {
		name += "/clean"
	}
	defer timeTrack(time.Now(), name)

	for _, idx := range sample {
		row := d.db.QueryRowx(`SELECT * FROM ik_table WHERE uuid_v4 = $1`, records[idx].UuidV4)
		var record Record
		if err := row.StructScan(&record); err != nil {
			log.Fatal("SelectUuidV4 fail")
		}
		//fmt.Println(record.ToCSV())
	}
}

func (d *DbTester) SelectSnowflake(records []Record, sample []int, withCache bool) {
	name := "SelectSnowflake"
	if withCache {
		name += "/cache"
	} else {
		name += "/clean"
	}
	defer timeTrack(time.Now(), name)

	for _, idx := range sample {
		row := d.db.QueryRowx(`SELECT * FROM ik_table WHERE snowflake_id = $1`, records[idx].Snowflake)
		var record Record
		if err := row.StructScan(&record); err != nil {
			log.Fatal("SelectSnowflake fail")
		}
		//fmt.Println(record.ToCSV())
	}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s: %s", name, elapsed)
}
