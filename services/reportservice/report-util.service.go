package reportservice

import (
	"distribution-system-be/utils/util"
	"encoding/csv"
	"fmt"
	"os"
)

func ExportToCSV(datas interface{}, dateStart, dateEnd, reportName string) string {
	filename := fmt.Sprintf("%v_%v_%v.csv", reportName, dateStart, dateEnd)
	csvFile, _ := os.Create(filename)

	defer func() {
		csvFile.Close()
	}()

	csvwriter := csv.NewWriter(csvFile)
	result := util.ToSliceData(datas)

	for _, data := range result {
		csvwriter.Write(data)
	}
	csvwriter.Flush()
	csvFile.Close()

	return filename
}
