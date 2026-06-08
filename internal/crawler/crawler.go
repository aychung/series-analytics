// Package crawler is the module that will perform the crawling
package crawler

import (
	"bytes"
	"context"
	"strconv"
	"strings"

	"series-analytics/internal/model"

	"github.com/PuerkitoBio/goquery"
)

type Crawler struct {
	fetcher *Fetcher
}

func New() (*Crawler, error) {
	return &Crawler{
		fetcher: NewFetcher(),
	}, nil
}

func (c *Crawler) QueryNovelInfo(ctx context.Context, prdNo string) (*model.NovelDetail, *model.NovelStat, error) {
	url := "https://series.naver.com/novel/detail.series?productNo=" + prdNo
	prdPageFetchResult, err := c.fetcher.Fetch(ctx, url)
	if err != nil {
		return nil, nil, err
	}

	detail, stat, err := parsePrdPageInfo(prdPageFetchResult.Body)
	if err != nil {
		return nil, nil, err
	}
	detail.PrdNo = prdNo

	return detail, stat, nil
}

func parsePrdPageInfo(htmlBody []byte) (*model.NovelDetail, *model.NovelStat, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlBody))
	if err != nil {
		return nil, nil, err
	}

	title := cleanText(doc.Find("div.end_head > h2").First().Text())
	infoList := doc.Find("li.info_lst > ul > li")
	category := cleanText(infoList.Eq(1).First().Text())
	author := cleanText(infoList.Eq(2).Find("a").First().Text())
	publisher := cleanText(infoList.Eq(3).Find("a").First().Text())
	commentCount := cleanText(doc.Find("span#commentCount").First().Text())
	downloadCount := cleanText(doc.Find("a.btn_download").First().Text())
	ratingStr := cleanText(doc.Find("div.score_area > em").First().Text())

	rating, err := strconv.ParseFloat(ratingStr, 64)
	if err != nil {
		return nil, nil, err
	}

	return &model.NovelDetail{
			PrdNo:     "",
			Title:     title,
			Author:    author,
			Publisher: publisher,
			Category:  category,
		}, &model.NovelStat{
			CommentCount:  commentCount,
			DownloadCount: downloadCount,
			Rating:        rating,
		}, nil
}

func cleanText(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}
