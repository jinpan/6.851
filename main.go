package main

import (
    "encoding/csv"
    "fmt"
    "os"
    "strconv"
)

func main() {
    sht := makeSimpleHashTable(uint(1024))
    fmt.Println("Finished initializing")

    file, _ := os.Open("tests/simple.csv")
    csvReader := csv.NewReader(file)
    for {
        data, err := csvReader.Read()
        if err != nil {
            break
        }
        key, _ := strconv.ParseUint(data[0], 10, 0)
        val := data[1]
        sht.insert(uint(key), val)
    }

}

