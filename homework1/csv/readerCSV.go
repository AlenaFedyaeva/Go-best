package csv

import (
	"encoding/csv"
	"fmt"
	"os"
)

func ReadCSVColumns(fname string) (data *DataCSVColumns, err error) {
	
	
	lines, err := ReadCsv(fname)
    if err != nil {
        return nil,err
    }

    // Loop through lines & turn into object
    for _, line := range lines {
        data = &DataCSVColumns{
            Column1: line[0],
            Column2: line[1],
			Column3: line[2],
        }
        fmt.Println(data.Column1 + " " + data.Column2 + " " + data.Column3)
    }

	return  data, nil
}

// ReadCsv accepts a file and returns its content as a multi-dimentional type
// with lines and each column. Only parses to string type.
func ReadCsv(filename string) ([][]string, error) {

    // Open CSV file
    f, err := os.Open(filename)
    if err != nil {
        return [][]string{}, err
    }
    defer f.Close()

    // Read File into a Variable
    lines, err := csv.NewReader(f).ReadAll()
    if err != nil {
        return [][]string{}, err
    }

    return lines, nil
}

func ReadCsvData(fname string) (data *DataCSV,err error) {
	// Open CSV file
	f, err := os.Open(fname)
	if err != nil {
		fmt.Println("An error encountered ::", err)
		return nil, err
	}


	defer func ()  {
		cerr:=f.Close()
		if err==nil {
			err=cerr
		}
	}()


	// Setup the reader
	reader := csv.NewReader(f)

	// Read 1 record from file
	header, err := reader.Read()
	if err != nil {
		fmt.Println("An error encountered ::", err)
		return nil, err
	}


	records, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }

	data =&DataCSV{
		Headers: header,
		Records: records,
	}

	data.Print()
	return data, nil 
}
