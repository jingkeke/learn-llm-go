package main

import (
	"context"
	"fmt"
	"github.com/smhanov/laconic/search"
)

func main() {
	p := search.NewDuckDuckGo()
	res, err := p.Search(context.Background(), "2026-04-23 A股 涨停分析")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Got %d results\n", len(res))
	for _, r := range res {
		fmt.Println("-", r.Title)
	}
}
