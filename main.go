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

    bht := makeBackgroundHashTable(1024)

    file, _ := os.Open("tests/simple.csv")
    csvReader := csv.NewReader(file)
    for i := 0; ; i++ {
        data, err := csvReader.Read()
        if err != nil {
            break
        }
        key, _ := strconv.ParseInt(data[0], 10, 0)
        val := data[1]
        if i % 10000 == 0 {
            fmt.Println(i)
        }
        bht.insert(int(key), val)
    }

}

