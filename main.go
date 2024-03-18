package main

import (
	"fmt"
	"github.com/soerenchrist/logsync/graph"
)

func main() {
	g, err := graph.ReadGraph("C:\\Users\\schrist\\OneDrive\\Logseq\\Personal")
	logErr(err)

	err = graph.SaveGraph(g, "save.json")
	logErr(err)
}

func logErr(err error) {
	fmt.Printf("error: %v", err)
}
