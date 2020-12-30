package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	urls = []string{
		// // this one is article, so it's parse-able
		//"https://www.nytimes.com/2019/02/20/climate/climate-national-security-threat.html",
		// // while this one is not an article, so readability will fail to parse.
		// "https://www.nytimes.com/",
		//"https://www.reuters.com/article/us-usa-drones-faa/u-s-to-allow-small-drones-to-fly-over-people-at-night-idUSKBN2921R8",
		//"https://www.dw.com/en/french-fashion-designer-pierre-cardin-dies-at-98/a-56083346",
		//"https://popey.com/blog/2020/12/contributing-without-code/",
		"http://www.chicagocoinclub.org/projects/PiN/juh.html",
	}
)

func main() {
	ctx := context.Background()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err.Error())
	}
	for _, url := range urls {
		dir, err := runToPdf(ctx, url, dir+"/fonts")
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println(dir + "/main.pdf")
	}
}
