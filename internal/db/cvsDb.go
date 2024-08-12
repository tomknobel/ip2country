package db

import (
	csv_utils "encoding/csv"
	"github.com/tomknobel/ip2country/internal/models"
	"go.uber.org/zap"
	"os"
)

var logger = zap.NewExample().Sugar()

type csvDb struct {
	csvFilePath string
	file        *os.File
}

func NewCsvDb(dbPath string) *csvDb {
	return &csvDb{
		csvFilePath: dbPath,
	}
}

func (csv *csvDb) Connect() error {
	file, err := os.OpenFile(csv.csvFilePath, os.O_RDWR, 0666)
	csv.file = file

	return err
}
func (csv *csvDb) Close() error {
	return csv.file.Close()
}
func (csv *csvDb) Find(ip string) (*models.Ip, error) {
	csv.Connect()
	defer csv.Close()
	reader := csv_utils.NewReader(csv.file)
	if _, err := reader.Read(); err != nil { // reading the headers
		logger.Errorf("Error reading header:%v", err)
		return nil, err
	}

	// Iterate through the records
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				logger.Info("End of file reached")
				break
			}
			logger.Errorf("Error reading record:%v", err)
			return &models.Ip{}, err
		}
		logger.Info(record)
		// Assume ID is the first column
		if record[0] == ip {
			return &models.Ip{
				Adders:  ip,
				Country: record[1],
				City:    record[2],
			}, nil
		}
	}
	return &models.Ip{}, nil
}
