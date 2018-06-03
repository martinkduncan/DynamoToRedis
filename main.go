package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func main() {

	// create native file
	input, err := os.Open("./FRB_CP.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer input.Close()

	reader := csv.NewReader(input)
	reader.FieldsPerRecord = -1 // infinite?????

	data, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	writer, err := os.Create("./FRB_CP.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer writer.Close()

	for _, i := range data {
		if _, err = strconv.Atoi(i[0][0:1]); err == nil {
			y, _ := strconv.Atoi(i[0][0:4])
			m, _ := strconv.Atoi(i[0][6:7])
			d, _ := strconv.Atoi(i[0][9:10])
			f := fmt.Sprintf("%04d%02d%02d", y, m, d)
			t := fmt.Sprintf("*3\r\n$3\r\nSET\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n", len(f), f, len(i[6]), i[6])
			writer.WriteString(t)
		}
	}

	// read native into Redis
	flusher := exec.Command("redis-cli", "flushdb")
	if err := flusher.Run(); err != nil {
		fmt.Println(err)
	}

	loader := exec.Command("/bin/sh", "-c", "cat ./FRB_CP.txt | redis-cli --pipe")
	if err := loader.Run(); err != nil {
		fmt.Println(err)
	}
}
