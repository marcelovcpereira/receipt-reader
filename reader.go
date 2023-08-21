package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type ReceiptItem struct {
	Name  string  `json:"name"`
	Value float32 `json:"value"`
}

type Receipt struct {
	ReceiptType string        `json:"type"`
	Date        time.Time     `json:"date"`
	Currency    string        `json:"currency"`
	Items       []ReceiptItem `json:"items"`
}

func ReadItem(line string) (ReceiptItem, error) {
	if strings.Contains(line, "EUR/kg") || line == "" {
		return ReceiptItem{}, errors.New("could not create item from line: " + line)
	} else {
		var parts = strings.Split(line, " ")
		var _ = parts[len(parts)-1]
		parts = parts[:len(parts)-1]
		var itemValue, err = strconv.ParseFloat(strings.Replace(parts[len(parts)-1], ",", ".", -1), 32)

		if err != nil {
			return ReceiptItem{}, err
		}
		parts = parts[:len(parts)-1]
		var itemName = strings.Join(parts, " ")
		return ReceiptItem{
			Name:  itemName,
			Value: float32(itemValue),
		}, nil
	}
}

func ReadLidlReceipt(file string) {
	output, err := exec.Command("tesseract", file, "stdout", "-l", "spa").Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	var raw = string(output)

	var start = strings.Index(raw, "EUR")
	var end = strings.Index(raw, "Total")
	var items = raw[start+3 : end]
	var itemsArray = strings.Split(items, "\n")
	var productsArray = []ReceiptItem{}
	for _, item := range itemsArray {
		var i, err = ReadItem(item)
		if err != nil {
			fmt.Printf("Skipping item: %v: (%v)\n", item, err)
			continue
		}
		productsArray = append(productsArray, i)
	}

	var receipt = Receipt{
		Currency:    "EUR",
		Items:       productsArray,
		Date:        time.Now(),
		ReceiptType: "lidl",
	}

	var res, e = json.Marshal(receipt)
	if e != nil {
		fmt.Printf("ERROR: %v\n", e)
	} else {
		fmt.Printf("FINAL RECEIPT: %v\n", string(res))
	}

}
