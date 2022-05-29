package geolocation

import (
	"encoding/csv"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/atulya-pandey/golang-geolocation/datastore"
)

type stats struct {
	timeElapsed int
	Accepted    int
	Discarded   int
}

type CsvImporter struct {
	datastore datastore.GeoLocationStore
	stats     *stats
}

func New(store datastore.GeoLocationStore) Importer {
	return CsvImporter{
		datastore: store,
		stats:     &stats{},
	}
}

func (c CsvImporter) Parse(csvFilePath string) (string, error) {
	csvFile, err := os.Open(csvFilePath)

	if err != nil {
		log.Println(err)
	}

	log.Print("Successfully Opened CSV file")

	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()

	parsedCsvLines, discarded := parse(csvLines)

	c.stats.Accepted = len(parsedCsvLines)
	c.stats.Discarded = discarded

	parasedCsvFileName := strings.Replace(csvFilePath, ".csv", "", 1) + "_parsed.csv"

	parsedCsvFile, err := os.Create(parasedCsvFileName)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer parsedCsvFile.Close()

	csvwriter := csv.NewWriter(parsedCsvFile)

	for _, row := range parsedCsvLines {
		_ = csvwriter.Write(row)
	}

	csvwriter.Flush()

	return parasedCsvFileName, nil
}

func (c CsvImporter) Import(geolocationData string) (string, error) {
	rowCount, err := c.datastore.LoadData(geolocationData)

	if err != nil {
		return "Failed", err
	}

	log.Printf("Added %d records in geo_loc table", rowCount)

	return "Success", nil
}

func (c CsvImporter) Statistics() {
	log.Printf("Accepted: %d", c.stats.Accepted)
	log.Printf("Discarded: %d", c.stats.Discarded)
}

func parse(csvRows [][]string) ([][]string, int) {
	// Remove duplicates
	// Check missing values
	// Check format
	dedupMap := make(map[string]bool, 0)
	var discardedRecords int

	parsedCsvRows := make([][]string, 0)

	for _, row := range csvRows {
		ip := row[0]

		if _, exist := dedupMap[ip]; exist {
			discardedRecords++
			continue
		} else if !checkFormat(row) {
			discardedRecords++
		} else {
			dedupMap[ip] = true
			parsedCsvRows = append(parsedCsvRows, row)
		}
	}

	return parsedCsvRows, discardedRecords
}

func checkFormat(row []string) bool {
	re := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)

	if !re.MatchString(row[0]) {
		log.Printf("Invalid IP %s. Discarding...", row[0])
		return false
	} else if len(row) > 7 {
		log.Printf("CSV Row length greater than expected. Discarding...")
		return false
	} else if len(row) < 7 {
		log.Printf("CSV Row length less than expected. Discarding...")
		return false
	} else {
		return true
	}
}
