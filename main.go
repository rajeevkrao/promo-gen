package main

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
)

func main() {
	parameters := ""
	if len(os.Args) > 1 {
		parameters = os.Args[1]
	}

	if parameters == "-r" || parameters == "-R" {
		err := deleteAppData()
		if err != nil {
			fmt.Println("Error deleting appdata:", err)
			return
		}
		fmt.Println("Appdata deleted.")
		return
	}

	link := loadAppData()

	id, err := scrapeIdFromLink(link)
	if err != nil {
		fmt.Println("Error scraping ID from link:", err)
		return
	}

	details, err := getDetailsFromId(id)
	if err != nil {
		fmt.Println("Error fetching channel details:", err)
		return
	}

	promoText := generatePromoText(details)
	fmt.Println(promoText)

	err = clipboard.WriteAll(promoText)
	if err != nil {
		fmt.Println("Error copying to clipboard:", err)
		return
	}
}
