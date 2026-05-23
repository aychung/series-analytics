# Series Analytics

Crawler and reporting project for tracking popular Naver Series web novel tags, focused only on novels tagged `로판` or `로맨스`.

The project is currently in planning/scaffolding. See:

- [Crawler plan](docs/series-crawler-plan.md)
- [Example daily report](docs/example-daily-report.md)

## Goal

Build a daily and weekly analytics pipeline that:

- Reads Naver Series ranking pages.
- Finds ranked novels tagged `로판` or `로맨스`.
- Stores novel details and tags in a local database.
- Caches ignored novels so irrelevant details are not refetched every run.
- Generates daily and weekly tag popularity reports.

## Crawl Policy

Do not crawl:

```text
https://series.naver.com/novel/categoryProductList.series
```

Naver Series `robots.txt` disallows `/novel/categoryProductList.series`.

Use Top 100 ranking pages instead:

```text
https://series.naver.com/novel/top100List.series?rankingTypeCode=DAILY&categoryCode=ALL&page=1..5
https://series.naver.com/novel/top100List.series?rankingTypeCode=WEEKLY&categoryCode=ALL&page=1..5
```

Novel detail pages are fetched only when needed to classify/cache a novel:

```text
https://series.naver.com/novel/detail.series?productNo={productNo}
```

## Inclusion Scope

Only novels with at least one of these tags are included in analytics:

```text
로판
로맨스
```

All other ranked novels are stored in an ignore list with their observed tags and a recheck date. Ignored novels are excluded from tag stats and reports, except for report audit sections.

## Planned Stack

- Node.js
- TypeScript
- `pnpm`
- SQLite
- Drizzle ORM or Prisma
- `cheerio` for HTML parsing
- `undici` or `got` for HTTP
- `tsx` for CLI execution

## Planned Commands

```bash
pnpm db:migrate
pnpm crawl:daily
pnpm crawl:weekly
pnpm report:daily
pnpm report:weekly
```

These commands are not implemented yet.

## Planned Output

Reports will be written under:

```text
reports/daily/
reports/weekly/
```

The database will live under:

```text
data/series.sqlite
```

## Current Files

```text
docs/series-crawler-plan.md
docs/example-daily-report.md
README.md
AGENTS.md
package.json
```

