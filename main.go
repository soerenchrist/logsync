package main

import (
	"fmt"
	"github.com/soerenchrist/logsync/graph"
	"os"
)

func main() {
	g, err := graph.ReadGraph("C:\\Users\\schrist\\OneDrive\\Logseq\\Personal")
	logErr(err)

	f, err := os.Open("save.json")
	defer f.Close()
	logErr(err)

	err = graph.SaveGraph(g, f)
	logErr(err)
}

func logErr(err error) {
	fmt.Printf("error: %v", err)
}
