// Package crawler is the module that will perform the crawling
package crawler

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"series-analytics/internal/model"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

type Crawler struct {
	fetcher *Fetcher
	pw *playwright.Playwright
	pwBrowser *playwright.BrowserContext
}

func New() (*Crawler, error) {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}

	browser, err := pw.Chromium.LaunchPersistentContext(chromiumUserData,
		playwright.BrowserTypeLaunchPersistentContextOptions{
			// Headless: new(false),
			ExecutablePath: new("/usr/bin/chromium"),
			Timeout: new(0.0),
		})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	return &Crawler{
		fetcher: NewFetcher(),
		pw: pw,
		pwBrowser: &browser,
	}, nil
}

func (c *Crawler) Close() {
	if err := (*c.pwBrowser).Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err := c.pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}

const top100url = "https://series.naver.com/novel/top100List.series?categoryCode=201&rankingTypeCode="

func (c *Crawler) QueryTop100(ctx context.Context, rankingType model.RankingType) (ranking [100]model.NovelRank, err error) {
	for pageNo := 1; pageNo <= 5; pageNo++ {
		url := top100url + rankingType.ToQuery() + "&page=" + strconv.Itoa(pageNo)
		prdPageFetchResult, nerr := c.fetcher.Fetch(ctx, url)
		if nerr != nil {
			return ranking, nerr
		}

		prds, nerr := parseTop100PageInfo(prdPageFetchResult.Body)
		if nerr != nil {
			return ranking, nerr
		}
		for i, prdNo := range prds {
			rank := (pageNo-1)*20 + i + 1
			idx := rank - 1
			ranking[idx] = model.NovelRank{
				PrdNo:   prdNo,
				Rank:    rank,
				NovelID: -1,
			}
		}
	}

	return
}

const chromiumUserData = "/home/al/.config/chromium"
func (c *Crawler) QueryNovelWithBrowser(ctx context.Context, prdNo string) (novel model.Novel, err error) {
	page, err := (*c.pwBrowser).NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	url := "https://series.naver.com/novel/detail.series?productNo=" + prdNo
	if _, err = page.Goto(url); err != nil {
		log.Fatalf("could not goto: %v", err)
	}
	title, err := page.Locator("div.end_head > h2").First().TextContent()
	if err != nil {
		log.Fatalf("could not parse title: %v", err)
	}
	infoList, err := page.Locator("li.info_lst > ul > li").All()
	if err != nil {
		log.Fatalf("could not parse infoList: %v", err)
	}
	category, err := infoList[1].First().TextContent()
	if err != nil {
		log.Fatalf("could not parse category: %v", err)
	}
	author, err := infoList[2].Locator("a").First().TextContent()
	if err != nil {
		log.Fatalf("could not parse author: %v", err)
	}
	publisher, err := infoList[3].Locator("a").First().TextContent()
	if err != nil {
		log.Fatalf("could not parse publisher: %v", err)
	}
	commentCount, err := page.Locator("span#commentCount").First().TextContent()
	if err != nil {
		log.Fatalf("could not parse commentCount: %v", err)
	}
	downloadCount, err := page.Locator("a.btn_download").First().TextContent()
	if err != nil {
		log.Fatalf("could not parse downloadCount: %v", err)
	}
	ratingStr, err := page.Locator("div.score_area > em").First().TextContent()
	if err != nil {
		log.Fatalf("could not parse ratingStr: %v", err)
	}
	metaDescription, err := page.Locator("meta[name='description']").GetAttribute("content")
	if err != nil {
		log.Fatalf("could not parse metaDescription: %v", err)
	}

	var tags []string
	if metaDescription != "" {
		tags = parseTags(metaDescription)
	}

	// println("> parsing float from : " + ratingStr)
	rating, err := strconv.ParseFloat(ratingStr, 64)
	if err != nil {
		return
	}

	novel.Detail = model.NovelDetail{
		PrdNo:     prdNo,
		Title:     title,
		Author:    author,
		Publisher: publisher,
		Category:  category,
	}
	novel.Stat = model.NovelStat{
		CommentCount:  commentCount,
		DownloadCount: downloadCount,
		Rating:        rating,
	}
	novel.Tags = tags
	return
}

func (c *Crawler) QueryNovel(ctx context.Context, prdNo string) (novel model.Novel, err error) {
	url := "https://series.naver.com/novel/detail.series?productNo=" + prdNo
	prdPageFetchResult, err := c.fetcher.Fetch(ctx, url)
	if err != nil {
		return
	}
	if strings.Contains(prdPageFetchResult.FinalURL, "nidlogin") {
		return c.QueryNovelWithBrowser(ctx, prdNo)
	}

	detail, stat, tags, err := parsePrdPageInfo(prdPageFetchResult.Body)
	if err != nil {
		return
	}
	detail.PrdNo = prdNo
	novel.Detail = detail
	novel.Stat = stat
	novel.Tags = tags

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

func parseTop100PageInfo(htmlBody []byte) (top100PrdNo []string, err error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlBody))
	if err != nil {
		return
	}

	errList := []error{}
	doc.Find("ul.comic_top_lst > li > a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		u, err := url.Parse(href)
		if err != nil {
			errList = append(errList, err)
		} else {
			prdNo := u.Query().Get("productNo")
			top100PrdNo = append(top100PrdNo, prdNo)
		}
	})
	if len(errList) > 0 {
		return nil, errors.New("error parsing top100 page")
	}

	return
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
