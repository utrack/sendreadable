package main

import "log"

var (
	urls = []string{
		// // this one is article, so it's parse-able
		//"https://www.nytimes.com/2019/02/20/climate/climate-national-security-threat.html",
		// // while this one is not an article, so readability will fail to parse.
		// "https://www.nytimes.com/",
		"https://www.reuters.com/article/us-usa-drones-faa/u-s-to-allow-small-drones-to-fly-over-people-at-night-idUSKBN2921R8",
	}
)

func main() {
	for _, url := range urls {
		err := runToPdf(url)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
