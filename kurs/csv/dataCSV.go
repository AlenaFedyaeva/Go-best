package csv

import (
	"fmt"
)

type DataCSV struct {
	Headers []string
	Records [][]string
}

func (data *DataCSV) Print() {
	fmt.Printf("Headers : %v \n", data.Headers)
	fmt.Printf("Records : %v \n", data.Records)
}

func newDataCSV(fname string) (*DataCSV, error) {

	return ReadCsvData(fname)
}



type DataCSVColumns struct{
	Column1 string
	Column2 string
	Column3 string
}

func (data *DataCSVColumns) Print(){
	len:=len(data.Column1)
	for i := 0; i < len; i++ {
		fmt.Println("i ",data.Column1, data.Column2,data.Column3)
	}
}