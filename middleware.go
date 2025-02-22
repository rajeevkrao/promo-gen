package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getAppDataPath() string {
	return filepath.Join(os.Getenv("APPDATA"), "promo-gen.txt")
}

func getChannelLink() string {
	var input string
	fmt.Println("Enter your channel link:")
	fmt.Scanln(&input)
	return input
}

func deleteAppData() error {
	appDataPath := getAppDataPath()
	err := os.Remove(appDataPath)
	if err != nil {
		fmt.Println("Error deleting appdata:", err)
		return err
	}
	return nil
}

func loadAppData() string {
	appDataPath := getAppDataPath()
	if _, err := os.Stat(appDataPath); os.IsNotExist(err) {
		link := getChannelLink()
		err := os.WriteFile(appDataPath, []byte(link), 0644)
		if err != nil {
			fmt.Println("Error writing appdata:", err)
			return ""
		}
	}
	data, err := os.ReadFile(appDataPath)
	if err != nil {
		fmt.Println("Error reading appdata:", err)
		return ""
	}
	return string(data)
}

func getDetailsFromId(id string) (ChannelDetails, error) {
	url := fmt.Sprintf("https://yt-det.vercel.app/api?channelId=%s", id)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching JSON:", err)
		return ChannelDetails{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Failed to fetch JSON:", resp.Status)
		return ChannelDetails{}, fmt.Errorf("failed to fetch JSON: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ChannelDetails{}, err
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return ChannelDetails{}, err
	}

	details := ChannelDetails{
		SubscriberCount:       response.Statistics.SubscriberCount,
		CustomUrl:             response.Snippet.CustomUrl,
		Title:                 response.Snippet.Title,
		HiddenSubscriberCount: response.Statistics.HiddenSubscriberCount,
	}

	return details, nil
}

func scrapeIdFromLink(link string) (string, error) {
	resp, err := http.Get(link)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch page: %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	channelId := doc.Find("meta[itemprop='identifier']").AttrOr("content", "Unknown")

	return channelId, nil
}

func formatNumber(num float64) string {
	var formatted string

	switch {
	case num >= 1_000_000_000:
		formatted = fmt.Sprintf("%.2fB", num/1_000_000_000)
	case num >= 1_000_000:
		formatted = fmt.Sprintf("%.2fM", num/1_000_000)
	case num >= 1_000:
		formatted = fmt.Sprintf("%.2fK", num/1_000)
	default:
		formatted = fmt.Sprintf("%.2f", num)
	}

	// Remove trailing zeros and unnecessary decimal points
	formatted = strings.TrimRight(formatted, "0")
	formatted = strings.TrimRight(formatted, ".")

	return formatted
}

func generatePromoText(details ChannelDetails) string {
	title := details.Title
	if details.CustomUrl != "" {
		title = details.CustomUrl
	}

	subscriberCount, err := strconv.ParseFloat(details.SubscriberCount, 64)
	if err != nil {
		fmt.Println("Error parsing subscriber count:", err)
		return ""
	}
	roundedSubscriberCount := formatNumber(subscriberCount)
	return fmt.Sprintf("Subscribe to %s Youtube Channel live with %s subs", title, roundedSubscriberCount)
}