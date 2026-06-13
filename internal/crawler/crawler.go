// Package crawler is the module that will perform the crawling
package crawler

import (
	"bytes"
	"context"
	"regexp"
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

func (c *Crawler) QueryTop100(ctx context.Context) ([]model.NovelDetail, error) {
	url := "https://series.naver.com/novel/top100List.series?rankingTypeCode=DAILY&categoryCode=201&page=1"
	prdPageFetchResult, err := c.fetcher.Fetch(ctx, url)
	if err != nil {
		return nil, err
	}

	parseTop100PageInfo(prdPageFetchResult.Body)

	return nil, nil
}

func (c *Crawler) QueryNovelInfo(ctx context.Context, prdNo string) (detail model.NovelDetail, stat model.NovelStat, tags []string, err error) {
	url := "https://series.naver.com/novel/detail.series?productNo=" + prdNo
	prdPageFetchResult, err := c.fetcher.Fetch(ctx, url)
	if err != nil {
		return
	}

	detail, stat, tags, err = parsePrdPageInfo(prdPageFetchResult.Body)
	if err != nil {
		return
	}
	detail.PrdNo = prdNo

	return
}

func parsePrdPageInfo(htmlBody []byte) (detail model.NovelDetail, stat model.NovelStat, tags []string, err error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlBody))
	if err != nil {
		return
	}

	title := cleanText(doc.Find("div.end_head > h2").First().Text())
	infoList := doc.Find("li.info_lst > ul > li")
	category := cleanText(infoList.Eq(1).First().Text())
	author := cleanText(infoList.Eq(2).Find("a").First().Text())
	publisher := cleanText(infoList.Eq(3).Find("a").First().Text())
	commentCount := cleanText(doc.Find("span#commentCount").First().Text())
	downloadCount := cleanText(doc.Find("a.btn_download").First().Text())
	ratingStr := cleanText(doc.Find("div.score_area > em").First().Text())
	metaDescription, exists := doc.Find("meta[name='description']").Attr("content")

	if exists {
		tags = parseTags(metaDescription)
	}

	rating, err := strconv.ParseFloat(ratingStr, 64)
	if err != nil {
		return
	}

	return model.NovelDetail{
			PrdNo:     "",
			Title:     title,
			Author:    author,
			Publisher: publisher,
			Category:  category,
		}, model.NovelStat{
			CommentCount:  commentCount,
			DownloadCount: downloadCount,
			Rating:        rating,
		}, tags, nil
}

func parseTop100PageInfo(htmlBody []byte) error {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlBody))
	if err != nil {
		return err
	}

	doc.Find("")
	return nil
}

func cleanText(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

func parseTags(s string) []string {
	re := regexp.MustCompile(`#([^\s,]+)`)

	matches := re.FindAllStringSubmatch(s, -1)

	var tagList []string
	for _, match := range matches {
		if len(match) > 1 {
			tagList = append(tagList, match[1])
		}
	}
	return tagList
}
