// Package crawler is the module that will perform the crawling
package crawler

import (
	"bytes"
	"context"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Crawler struct {
	fetcher *Fetcher
}

type NovelStat struct {
	CommentCount string
	DownloadCount string
	Rating float64
}

func New() (*Crawler, error) {
	return &Crawler{
		fetcher: NewFetcher(),
	}, nil
}

func (c *Crawler)QueryNovelStat(ctx context.Context, prdNo string) (*NovelStat, error) {
	url := "https://series.naver.com/novel/detail.series?productNo=" + prdNo
	prdPageFetchResult, err := c.fetcher.Fetch(ctx, url)
	if err != nil {
		return nil, err
	}
	return parsePrdPageStat(prdPageFetchResult.Body)
}

func parsePrdPageStat(htmlBody []byte) (*NovelStat, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	commentCount := cleanText(doc.Find("span#commentCount").First().Text())
	downloadCount := cleanText(doc.Find("a.btn_download").First().Text())
	ratingStr := cleanText(doc.Find("div.score_area > em").First().Text())

	rating, err := strconv.ParseFloat(ratingStr, 64)
	if err != nil {
		return nil, err
	}

	return &NovelStat{
		Rating: rating,
		DownloadCount: downloadCount,
		CommentCount: commentCount,
	}, nil
}

func cleanText(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}
