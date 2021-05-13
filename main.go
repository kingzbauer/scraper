package main

import (
	"fmt"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/trustmaster/goflow"

	"github.com/kingzbauer/scraper/components"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	graph := goflow.NewGraph()

	handle(graph.Add("retriever", components.NewRetriever(soup.Get)))
	handle(graph.Add("extractor", components.NewExtractor()))
	handle(graph.Add("spreader", components.NewSpreader()))
	handle(graph.Add("printer", components.NewPrinter()))

	handle(graph.Connect("retriever", "Out", "extractor", "In"))
	handle(graph.Connect("extractor", "Out", "spreader", "In"))
	handle(graph.Connect("spreader", "Out", "printer", "In"))

	graph.MapInPort("In", "retriever", "In")
	urlChan := make(chan string)
	handle(graph.SetInPort("In", urlChan))

	wait := goflow.Run(graph)

	for {
		var url string
		fmt.Printf("URL: ")
		fmt.Scanf("%s", &url)
		url = strings.TrimSpace(url)
		if url == "q" {
			close(urlChan)
			break
		}
		urlChan <- url
	}

	<-wait
}
