package main

import (
    "encoding/csv"
    "os"
    "fmt"
    "strconv"
)

var (
    u = 32000000
)

func main() {

    sht := makeSimpleHashTable(16, true)

    file, _ := os.Open("tests/insert.csv")
    csvReader := csv.NewReader(file)
    for i := 0; ; i++ {
        data, err := csvReader.Read()
        if err != nil {
            break
        }
        key, _ := strconv.ParseInt(data[1], 0, 0)
        val := data[2]
        if i % 10000 == 0 {
            fmt.Println(i)
        }
        sht.insert(int(key), val)
    }

}

