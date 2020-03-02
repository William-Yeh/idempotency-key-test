package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/dustin/go-humanize"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	//"log"
	"math/rand"
	"strconv"
	"time"
)

const USAGE = `Idempotency Key Performance Test.

Usage:
  iktest <dsn> [options] [-n <times>]

Options:
  -n <times>     Repeat n times [default: 100000] (100k).
  -p, --pause    Pause between query sessions for cleaning database cache.
  -v, --verbose  Verbose mode.
`

var (
	dsn           string
	pauseForCache bool
	numOfEntries  int
	allRecords    []Record
)

type Record struct {
	UuidV1    string    `db:"uuid_v1"`
	UuidV4    string    `db:"uuid_v4"`
	Snowflake int64     `db:"snowflake_id"`
	Ctime     time.Time `db:"ctime"`
}

func (record *Record) ToCSV() string {
	return fmt.Sprintf("\"%s\",\"%s\",%d,%d\n",
		record.UuidV1,
		record.UuidV4,
		record.Snowflake,
		record.Ctime.Unix())
}

func main() {
	processCmdline()
	generator, err := NewIKGenerator()
	if err != nil {
		log.Fatal(err)
		return
	}

	genItemsInMemory(generator, numOfEntries)
	//showAllRecords()
	sample := genSample(numOfEntries, numOfEntries/10)
	//fmt.Printf("%v\n", sample)

	t := NewDbTester(dsn)
	t.InitDb()
	t.ClearRecords()

	// INSERT INTO
	t.InsertUuidV1(allRecords)
	t.ClearRecords()
	t.InsertUuidV4(allRecords)
	t.ClearRecords()
	t.InsertSnowflake(allRecords)
	t.ClearRecords()

	t.InsertAll(allRecords)

	// SELECT
	t.CleanDbCache()
	t.SelectUuidV1(allRecords, sample, false)
	t.SelectUuidV1(allRecords, sample, true)

	t.CleanDbCache()
	t.SelectUuidV4(allRecords, sample, false)
	t.SelectUuidV4(allRecords, sample, true)

	t.CleanDbCache()
	t.SelectSnowflake(allRecords, sample, false)
	t.SelectSnowflake(allRecords, sample, true)
}

// processCmdline parses and validates cmdline args
func processCmdline() map[string]interface{} {
	args, _ := docopt.ParseDoc(USAGE)
	//fmt.Println(args)

	dsn = args["<dsn>"].(string)
	numOfEntries, _ = strconv.Atoi(args["-n"].(string))
	pauseForCache = args["--pause"].(bool)

	verbose := args["--verbose"].(bool)
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	return args
}

func genItemsInMemory(generator *IKGenerator, num int) {
	log.Printf("Generating %s items...\n", humanize.Comma(int64(num)))
	allRecords = make([]Record, 0, num)
	for i := 0; i < num; i++ {
		record := &Record{
			UuidV1:    generator.genUuidV1(),
			UuidV4:    generator.genUuidV4(),
			Snowflake: generator.genSnowflakeID(),
			Ctime:     time.Now(),
		}
		//fmt.Printf(record.ToCSV())
		allRecords = append(allRecords, *record)
	}
}

func showAllRecords() {
	for _, record := range allRecords {
		csvLine := record.ToCSV()
		fmt.Print(csvLine)

	}
}

// genSample generates a random sample of size num.
func genSample(max, num int) []int {
	result := make([]int, 0, num)

	rand.Seed(time.Now().Unix())
	for i := 0; i < num; i++ {
		result = append(result, rand.Intn(max))
	}

	return result
}
