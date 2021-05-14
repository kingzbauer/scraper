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
	handle(graph.Add("imageParser", components.NewImageParser()))
	handle(graph.Add("printer", components.NewPrinter()))
	handle(graph.Add("errorHandler", components.NewErrorHandler()))
	handle(graph.Add("debugger", components.NewDebugger()))
	handle(graph.Add("imageRetriever", components.NewImageRetriever()))
	handle(graph.Add("localFileSaver", components.NewLocalFileSaver("./downloads")))

	handle(graph.Connect("retriever", "Out", "extractor", "In"))
	handle(graph.Connect("extractor", "Out", "imageParser", "In"))
	handle(graph.Connect("imageParser", "Out", "imageRetriever", "In"))
	handle(graph.Connect("imageRetriever", "Out", "localFileSaver", "In"))
	handle(graph.Connect("localFileSaver", "Debug", "debugger", "In[localFileSaver]"))
	handle(graph.Connect("imageRetriever", "Debug", "debugger", "In[imageRetriever]"))
	// Error handler connections
	handle(graph.Connect("retriever", "Err", "errorHandler", "In[retriever]"))
	handle(graph.Connect("imageRetriever", "Err", "errorHandler", "In[imageRetriever]"))
	handle(graph.Connect("localFileSaver", "Err", "errorHandler", "In[localFileSaver]"))

	graph.MapInPort("In", "retriever", "In")
	urlChan := make(chan string)
	handle(graph.SetInPort("In", urlChan))

	wait := goflow.Run(graph)

	for {
		var url string
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
