package main

import (
	mycsv "Go-best/homework1/csv"
)


func main() {
	// csvReader()
	fname:="/home/alena/go/src/Go-best/homework1/records.csv"
	data1,_:=mycsv.ReadCsvData(fname)
	data1.Print()
	data2, _:= mycsv.ReadCSVColumns(fname)
	data2.Print()
}
