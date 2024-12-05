package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Topic struct {
	Name string
	Tag  string
}

func main() {
	// Read the maximum number of paragraphs per topic
	maxParagraphs := flag.Int("maxParagraphs", 100, "The maximum number of paragraphs per topic")

	// Read the output file name
	outputFile := flag.String("outputFile", "data.csv", "The output file name")

	// Parse the flags
	flag.Parse()

	// Define base URL
	baseURL := "https://vnexpress.net/"

	// Define topics
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
	file, err := os.Create(*outputFile)
	if err != nil {
		fmt.Println("Error creating file:", err)
	}

	// Create a CSV writer
	writer := csv.NewWriter(file)
	writer.Write([]string{"Paragraph", "Topic"})

	// Create a link collector
	linkCollector := colly.NewCollector()

	// Create a paragraph collector
	paragraphCollector := colly.NewCollector()

	var topicName string
	var countParagraphs int
	var foundLink bool

	// Set up the paragraph collector
	paragraphCollector.OnHTML("article.fck_detail", func(e *colly.HTMLElement) {
		if countParagraphs >= *maxParagraphs {
			return
		}
		paragraphs := e.ChildTexts("p:not(.Image)")
		for _, paragraph := range paragraphs {
			if countParagraphs >= *maxParagraphs {
				break
			}
			if len(strings.Fields(paragraph)) < 10 {
				continue
			}
			countParagraphs++
			writer.Write([]string{paragraph, topicName})
		}
	})

	// Set up the link collector
	linkCollector.OnHTML(".title-news", func(e *colly.HTMLElement) {
		if countParagraphs >= *maxParagraphs {
			return
		}
		foundLink = true
		link := e.ChildAttr("a[title]", "href")
		fmt.Println("Visiting", link)
		paragraphCollector.Visit(link)
	})

	for _, topic := range topics {
		// Visit multiple pages of a topic
		topicName = topic.Name
		countParagraphs = 0
		foundLink = false
		for page := 1; countParagraphs < *maxParagraphs; page++ {
			url := fmt.Sprintf("%s%s-p%d", baseURL, topic.Tag, page)
			linkCollector.Visit(url)
			if !foundLink {
				break
			}
		}
	}

	defer func() {
		writer.Flush()
		file.Close()
	}()
}
