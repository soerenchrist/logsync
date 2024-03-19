package main

import (
	"fmt"
	"github.com/soerenchrist/logsync/client/internal/config"
	"github.com/soerenchrist/logsync/client/internal/sync"
)

func main() {

	conf, err := config.Read()
	if err != nil {
		fmt.Printf("Failed to read config: %v", err)
		panic(-1)
	}

	fmt.Printf("%v\n", conf)

	sync.Start(conf)

	/*
		g, err := graph.ReadGraph("")
		logErr(err)

		f, err := os.Open("save.json")
		defer f.Close()
		logErr(err)

		err = graph.SaveGraph(g, f)
		logErr(err)

	*/
}

func logErr(err error) {
	fmt.Printf("error: %v", err)
}
