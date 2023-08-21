package main

import "fmt"

func main() {
	purchases, err := CollectFromCarrefour("marcelovcpereira@gmail.com", "celo1987")
	if err != nil {
		fmt.Printf("Error collecting from carrefour: %v\n", err)
	}
	fmt.Printf("Found %d purchases", len(purchases))
	// ReadLidlReceipt("/home/marce/images/lidl1.jpg")

}
