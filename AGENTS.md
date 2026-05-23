# AGENTS.md

Guidance for agents working in this repository.

## Project Purpose

This project tracks Naver Series web novel tag popularity for novels tagged `로판` or `로맨스`.

The core workflow is:

1. Crawl robots-permitted Top 100 ranking pages.
2. Fetch detail pages only when needed.
3. Include novels tagged `로판` or `로맨스`.
4. Store all other inspected novels in an ignore list.
5. Generate daily and weekly reports from included novels only.

## Required Reading

Before implementation work, read:

- [docs/series-crawler-plan.md](docs/series-crawler-plan.md)
- [docs/example-daily-report.md](docs/example-daily-report.md)

## Crawl Rules

Do not crawl `/novel/categoryProductList.series`. It is disallowed by Naver Series `robots.txt`.

Allowed planned ranking source:

```text
/novel/top100List.series?rankingTypeCode=DAILY&categoryCode=ALL&page=1..5
/novel/top100List.series?rankingTypeCode=WEEKLY&categoryCode=ALL&page=1..5
```

Detail pages may be fetched for classification and cached metadata:

```text
/novel/detail.series?productNo={productNo}
```

Use conservative HTTP behavior:

- Set a clear user agent.
- Use request timeouts.
- Use retry with backoff.
- Rate-limit requests.
- Re-check `robots.txt` periodically.

## Data Scope

Analytics must include only novels whose parsed tags contain at least one exact required tag:

```text
로판
로맨스
```

Novels without either required tag must be written to `ignored_novels` and excluded from:

- `novels`
- `novel_tags`
- `ranking_entries`
- tag stats
- main report analytics

Ignored novels may appear only in audit/report-health sections.

## Database Expectations

The planned schema includes:

- `novels`
- `tags`
- `novel_tags`
- `ignored_novels`
- `ranking_snapshots`
- `ranking_entries`
- `ignored_ranking_entries`
- `tag_daily_stats`
- `crawl_runs`
- `reports`

Prefer migrations over ad hoc schema creation once implementation begins.

## Implementation Style

- Use TypeScript.
- Keep parsing code isolated from storage code.
- Keep HTTP logic isolated from parsing logic.
- Add tests around parsers using saved HTML fixtures.
- Do not hard-code dates in implementation code.
- Keep generated data out of git unless explicitly requested.

Suggested structure:

```text
src/
  cli.ts
  config.ts
  db/
  crawler/
  analytics/
  reports/
```

## Report Rules

Daily and weekly reports should be based on included `로판`/`로맨스` novels only.

Reports should include:

- crawl summary
- scope filter summary
- top included tags by count
- top included tags by weighted score
- risers and fallers
- new included tags
- top included novels
- ignored novel summary
- data quality notes

Use [docs/example-daily-report.md](docs/example-daily-report.md) as the shape to preserve unless the user asks to change it.

## Safety Notes

- Do not add code that crawls disallowed category pages.
- Do not scrape viewer/content pages.
- Do not bypass access controls or login-only areas.
- Do not refetch every detail page every run; respect the cache and ignore TTL design.

