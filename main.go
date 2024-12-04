package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/gocolly/colly/v2"
)

type Topic struct {
	Name string
	Tag string
}

func main() {
	// Read the maximum number of titles per topic
	var maxTitles int
	fmt.Print("Enter the maximum number of titles per topic: ")
	_, err := fmt.Scanf("%d", &maxTitles)
	if err != nil {
		fmt.Println("Error reading the maximum number of titles per topic:", err)
	}

	// Define base URL
	baseURL := "https://vnexpress.net/"

	// Define topics
	// topics := []string{
	// 	"thoi-su", "the-gioi", "kinh-doanh",
	// 	"bat-dong-san", "khoa-hoc", "giai-tri",
	// 	"the-thao", "phap-luat", "giao-duc",
	// 	"suc-khoe", "doi-song", "du-lich",
	// 	"so-hoa", "oto-xe-may", "y-kien",
	// }
	topics := []Topic{
		{"Thời sự", "thoi-su"},
		{"Thế giới", "the-gioi"},
		{"Kinh doanh", "kinh-doanh"},
		{"Bất động sản", "bat-dong-san"},
		{"Khoa học", "khoa-hoc"},
		{"Giải trí", "giai-tri"},
		{"Thể thao", "the-thao"},
		{"Pháp luật", "phap-luat"},
		{"Giáo dục", "giao-duc"},
		{"Sức khỏe", "suc-khoe"},
		{"Đời sống", "doi-song"},
		{"Du lịch", "du-lich"},
		{"Số hóa", "so-hoa"},
		{"Ôtô - Xe máy", "oto-xe-may"},
		{"Ý kiến", "y-kien"},
	}

	// Create a CSV file
	file, err := os.Create("data.csv")
	if err != nil {
		fmt.Println("Error creating file:", err)
	}

	// Create a CSV writer
	writer := csv.NewWriter(file)
	writer.Write([]string{"Title", "Topic"})

	// Create a Collector
	c := colly.NewCollector()

	var topicName string
	var countTitles int
	var foundTitle bool

	// Set up the Collector
	c.OnHTML(".title-news", func(e *colly.HTMLElement) {
		if countTitles >= maxTitles {
			return
		}
		title := e.ChildText("a[title]")
		err := writer.Write([]string{title, topicName})
		if err != nil {
			fmt.Println("Error writing record to CSV:", err)
		}
		countTitles++
		foundTitle = true
	})

	for _, topic := range topics {
		// Visit URL
		topicName = topic.Name
		countTitles = 0
		for page := 1; countTitles < maxTitles; page++ {
			var url string
			if page == 1 {
				url = baseURL + topic.Tag
			} else {
				url = fmt.Sprintf("%s%s-p%d", baseURL, topic.Tag, page)
			}
			foundTitle = false
			c.Visit(url)
			if !foundTitle {
				break
			}
		}
	}

	defer func() {
		writer.Flush()
		file.Close()
	}()
}
