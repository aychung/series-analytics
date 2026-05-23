# Naver Series Novel Tag Analytics Plan

## Goal

Build a crawler and reporting pipeline that tracks popular Naver Series web novel tags over time.

The crawler should:

- Collect daily and weekly ranking data from robots-permitted ranking pages.
- Cache novel detail metadata in a database so details are not fetched every crawl.
- Extract tags from novel detail pages.
- Keep only novels tagged `로판` or `로맨스` for analytics.
- Track ignored novels in a database table so the crawler can skip known irrelevant novels later.
- Generate daily and weekly reports that show which romance/rofan-adjacent tags are popular, rising, falling, and attached to high-ranked novels.

## Important Crawl Constraint

Do not crawl this category endpoint:

```text
https://series.naver.com/novel/categoryProductList.series?categoryTypeCode=genre&genreCode=201
```

Naver Series `robots.txt` disallows:

```text
/novel/categoryProductList.series
```

Use `top100List.series` instead. It is a better source for popularity analysis because it exposes ranked daily, weekly, monthly, and category-filtered lists.

## Scope Filter

Only novels with at least one of these required tags should be included in the analytics dataset:

```text
로판
로맨스
```

All other ranked novels should be ignored after their detail tags are inspected. Ignored novels should be written to the `ignored_novels` table with the observed tags and ignore reason.

Important behavior:

- Ranking pages are still parsed from Top 100 because they are the popularity source.
- New ranked novels need one detail fetch to determine whether they have `로판` or `로맨스`.
- If a novel has neither required tag, it is added to `ignored_novels`.
- Future crawls can skip detail fetches for ignored novels until the ignore TTL expires.
- Ignored novels should not appear in `novels`, `novel_tags`, tag stats, or report analytics unless we explicitly add an audit section.

Default ignore TTL:

```text
30 days
```

The TTL exists because tags can change or parsing can improve later. After the TTL, the crawler may re-check an ignored novel if it appears in the ranking again.

## Source Pages

Daily ranking pages:

```text
https://series.naver.com/novel/top100List.series?rankingTypeCode=DAILY&categoryCode=ALL&page=1
https://series.naver.com/novel/top100List.series?rankingTypeCode=DAILY&categoryCode=ALL&page=2
https://series.naver.com/novel/top100List.series?rankingTypeCode=DAILY&categoryCode=ALL&page=3
https://series.naver.com/novel/top100List.series?rankingTypeCode=DAILY&categoryCode=ALL&page=4
https://series.naver.com/novel/top100List.series?rankingTypeCode=DAILY&categoryCode=ALL&page=5
```

Weekly ranking pages:

```text
https://series.naver.com/novel/top100List.series?rankingTypeCode=WEEKLY&categoryCode=ALL&page=1..5
```

Novel detail pages:

```text
https://series.naver.com/novel/detail.series?productNo={productNo}
```

Detail pages include tag-like metadata in the page meta description, for example:

```text
#NOVEL, #로판, #궁정로맨스, #로맨틱코미디
```

## Recommended Stack

- Node.js
- TypeScript
- `pnpm`
- SQLite for local storage
- Drizzle ORM or Prisma for schema and migrations
- `cheerio` for HTML parsing
- `undici` or `got` for HTTP requests
- `tsx` for running TypeScript commands
- Cron, GitHub Actions, or a small scheduler process for daily execution

SQLite is enough for the first version because this project stores small daily snapshots and cached metadata. Move to Postgres later if the project adds a web dashboard, multiple users, or broader crawl coverage.

## Project Structure

```text
src/
  cli.ts
  config.ts
  db/
    client.ts
    migrate.ts
    schema.ts
  crawler/
    crawlNovelDetail.ts
    crawlTop100.ts
    http.ts
    parsers.ts
    robots.ts
  analytics/
    dailyReport.ts
    tagStats.ts
    weeklyReport.ts
  reports/
    renderJson.ts
    renderMarkdown.ts

data/
  series.sqlite

reports/
  daily/
  weekly/

docs/
  series-crawler-plan.md
  example-daily-report.md
```

## Database Schema

### `novels`

Stores stable metadata and detail-page cache state for included novels only. A novel is included when its parsed tags contain `로판` or `로맨스`.

```sql
id integer primary key
product_no text not null unique
title text not null
author text
publisher text
category text
status text
description text
cover_url text
detail_url text not null
first_seen_at datetime not null
last_seen_at datetime not null
detail_fetched_at datetime
detail_hash text
detail_fetch_status text -- success, failed, skipped
detail_fetch_error text
created_at datetime not null
updated_at datetime not null
```

### `tags`

Canonical tag names.

```sql
id integer primary key
name text not null unique
created_at datetime not null
```

### `novel_tags`

Many-to-many relationship between novels and tags.

```sql
novel_id integer not null references novels(id)
tag_id integer not null references tags(id)
created_at datetime not null
primary key (novel_id, tag_id)
```

### `ignored_novels`

Stores ranked novels that were inspected and excluded because they do not match the required `로판`/`로맨스` scope.

```sql
id integer primary key
product_no text not null unique
title text not null
detail_url text not null
ranking_type text
category_code text
first_seen_at datetime not null
last_seen_at datetime not null
detail_fetched_at datetime
ignore_reason text not null -- missing_required_tags, detail_fetch_failed, parser_failed
observed_tags_json text -- JSON array of parsed non-matching tags
required_tags_checked_json text not null -- JSON array, default ["로판","로맨스"]
next_recheck_at datetime
created_at datetime not null
updated_at datetime not null
```

Suggested indexes:

```sql
create index idx_ignored_novels_next_recheck_at on ignored_novels(next_recheck_at);
create index idx_ignored_novels_last_seen_at on ignored_novels(last_seen_at);
```

### `ranking_snapshots`

One crawl snapshot for a ranking type/category/date.

```sql
id integer primary key
snapshot_date date not null
ranking_type text not null -- DAILY, WEEKLY, MONTHLY
category_code text not null -- ALL, 201, 207, 202, etc.
source text not null -- naver_series_top100
source_url text not null
crawled_at datetime not null
created_at datetime not null
unique(snapshot_date, ranking_type, category_code)
```

### `ranking_entries`

Included ranked novels inside a snapshot. Excluded ranked novels are stored in `ignored_ranking_entries`.

```sql
id integer primary key
snapshot_id integer not null references ranking_snapshots(id)
novel_id integer not null references novels(id)
product_no text not null
rank integer not null
title_at_crawl text not null
detail_url text not null
created_at datetime not null
unique(snapshot_id, rank)
unique(snapshot_id, product_no)
```

### `ignored_ranking_entries`

Stores the ranking position where an ignored novel appeared. This lets reports show how many ranked items were filtered out without mixing them into analytics.

```sql
id integer primary key
snapshot_id integer not null references ranking_snapshots(id)
product_no text not null
title_at_crawl text not null
rank integer not null
detail_url text not null
ignore_reason text not null
created_at datetime not null
unique(snapshot_id, product_no)
```

### `tag_daily_stats`

Materialized daily tag statistics. This can be computed on demand at first, then stored once the report format stabilizes.

```sql
id integer primary key
snapshot_id integer not null references ranking_snapshots(id)
tag_id integer not null references tags(id)
novel_count integer not null
weighted_score integer not null
average_rank real not null
best_rank integer not null
worst_rank integer not null
rank_1_20_count integer not null
rank_21_50_count integer not null
rank_51_100_count integer not null
created_at datetime not null
unique(snapshot_id, tag_id)
```

### `crawl_runs`

Tracks operational health.

```sql
id integer primary key
run_type text not null -- daily, weekly, backfill, manual
started_at datetime not null
finished_at datetime
status text not null -- running, success, partial, failed
pages_fetched integer not null default 0
details_fetched integer not null default 0
details_skipped_cached integer not null default 0
details_failed integer not null default 0
error_message text
created_at datetime not null
```

### `reports`

Stores generated report output for review and history.

```sql
id integer primary key
report_type text not null -- daily, weekly
report_date date not null
format text not null -- markdown, json, html
content text not null
created_at datetime not null
unique(report_type, report_date, format)
```

## Cache Policy

Ranking pages are fetched on every crawl because ranking position is the daily signal.

Novel detail pages are fetched only when one of these is true:

- The novel is new.
- The novel has no tags.
- `detail_fetched_at` is older than the configured TTL.
- The previous detail fetch failed.
- A manual refresh is requested.
- The novel is in `ignored_novels`, appears in rankings again, and `next_recheck_at` has passed.

Default TTL:

```text
30 days
```

Default request pacing:

```text
1 request every 1-3 seconds
2 retry attempts with exponential backoff
15 second request timeout
```

This keeps daily crawls light. A normal daily run should fetch 5 ranking pages and only a small number of new or stale detail pages.

Ignored novels also use a cache. A known ignored novel should be skipped until `next_recheck_at` unless manually refreshed.

## Ranking Crawl Flow

1. Start a `crawl_runs` row.
2. Fetch and check `robots.txt` periodically.
3. Fetch ranking pages `page=1..5`.
4. Parse each ranking item:
   - rank
   - `productNo`
   - title
   - detail URL
5. Create or replace the `ranking_snapshots` row for the date/type/category.
6. For each ranked item, check whether it already exists in `novels`.
7. If it exists in `novels`, insert a `ranking_entries` row and refresh details only if stale.
8. If it exists in `ignored_novels` and `next_recheck_at` has not passed, insert an `ignored_ranking_entries` row and skip the detail fetch.
9. If it is new or due for recheck, fetch the detail page.
10. Parse metadata and tags.
11. If tags contain `로판` or `로맨스`:
    - Upsert into `novels`.
    - Remove any stale `ignored_novels` row for the same `productNo`.
    - Upsert tags into `tags`.
    - Replace that novel's `novel_tags` rows.
    - Insert a `ranking_entries` row.
12. If tags do not contain `로판` or `로맨스`:
    - Upsert into `ignored_novels`.
    - Store observed tags and reason `missing_required_tags`.
    - Insert an `ignored_ranking_entries` row.
13. Compute tag stats from included `ranking_entries` only.
14. Generate report files and store report rows.
15. Finish the `crawl_runs` row.

## Detail Parser

The detail parser should extract:

- Title
- Product number
- Cover image URL
- Description
- Meta description
- Tags from hashtag tokens
- Category/genre if available
- Author if available
- Publisher if available
- Completion/status if available

Tag extraction rule:

```text
Find hashtag tokens from meta description, normalize whitespace, remove leading "#", discard generic "NOVEL" unless explicitly requested.
```

Example:

```text
"160 화 완결, #NOVEL, #로판, #궁정로맨스, #로맨틱코미디"
```

Extracted tags:

```text
로판
궁정로맨스
로맨틱코미디
```

Inclusion rule:

```text
include = tags contains "로판" or tags contains "로맨스"
```

Examples:

```text
["로판", "궁정로맨스"] => include
["로맨스", "재회물"] => include
["판타지", "무협", "회귀"] => ignore
```

## Popularity Metrics

Each tag should have these base metrics:

- `novel_count`: number of included ranked novels with the tag.
- `weighted_score`: rank-aware popularity score.
- `average_rank`: average ranking position for novels with the tag.
- `best_rank`: highest-ranked novel using the tag.
- `rank_1_20_count`: count of novels in ranks 1-20.
- `rank_21_50_count`: count of novels in ranks 21-50.
- `rank_51_100_count`: count of novels in ranks 51-100.

All metrics are computed from included `로판`/`로맨스` novels only. Ignored novels are excluded from tag popularity calculations.

Default weighted score:

```text
score = 101 - rank
```

Examples:

```text
Rank 1 => 100 points
Rank 20 => 81 points
Rank 100 => 1 point
```

Trend metrics:

- `count_delta_1d`: today's count minus yesterday's count.
- `score_delta_1d`: today's weighted score minus yesterday's score.
- `count_delta_7d`: today's count minus count from seven days ago.
- `score_delta_7d`: today's weighted score minus score from seven days ago.
- `new_today`: tag appeared today but not yesterday.
- `dropped_today`: tag appeared yesterday but not today.

## Daily Report Requirements

The daily report should answer:

- How much of the Top 100 matched the `로판`/`로맨스` scope?
- Which tags are most common among included romance/rofan novels today?
- Which included tags are strongest after considering rank?
- Which included tags rose or fell compared to yesterday?
- Which new included tags appeared today?
- Which high-ranked included novels are driving each major tag?
- How healthy was the crawl?

Report sections:

1. Header and crawl summary
2. Scope filter summary
3. Top tags by included novel count
4. Top tags by weighted score
5. Biggest risers
6. Biggest fallers
7. New tags
8. Top included novels and their tags
9. Ignored novel summary
10. Notable tag-to-novel examples
11. Data quality notes

See [example-daily-report.md](./example-daily-report.md).

## Weekly Report Requirements

The weekly report should answer:

- How many ranked novels matched the `로판`/`로맨스` scope across the week?
- Which included tags dominated the week?
- Which included tags gained momentum week over week?
- Which included tags appeared consistently across the week?
- Which included tags had high rank-weighted performance despite lower count?
- Which included novels repeatedly drove the week's strongest tags?

Weekly report sections:

1. Header and date range
2. Scope filter summary
3. Weekly top tags by total included appearances
4. Weekly top tags by weighted score
5. Week-over-week risers
6. Week-over-week fallers
7. Consistent included tags, measured by days present
8. Breakout included tags, measured by late-week acceleration
9. Top included novels contributing to weekly tag strength
10. Ignored novel summary
11. Data quality notes

Weekly stats should initially aggregate seven daily snapshots. Later, also compare against Naver's `rankingTypeCode=WEEKLY` snapshot as a separate source.

## CLI Commands

Target scripts:

```json
{
  "scripts": {
    "db:migrate": "tsx src/db/migrate.ts",
    "crawl:daily": "tsx src/cli.ts crawl daily",
    "crawl:weekly": "tsx src/cli.ts crawl weekly",
    "report:daily": "tsx src/cli.ts report daily",
    "report:weekly": "tsx src/cli.ts report weekly",
    "backfill": "tsx src/cli.ts backfill"
  }
}
```

Example manual commands:

```bash
pnpm crawl:daily
pnpm report:daily
pnpm report:weekly
```

## Scheduling

Use Korea time if reports should align with Naver's market day.

Example server cron:

```cron
15 6 * * * cd /home/al/Workspace/node/series-analytics && pnpm crawl:daily && pnpm report:daily
30 6 * * 1 cd /home/al/Workspace/node/series-analytics && pnpm report:weekly
```

## Implementation Phases

### Phase 1: Foundation

- Add TypeScript tooling.
- Add SQLite and ORM.
- Create migrations.
- Add typed config.
- Add basic CLI shell.

### Phase 2: Crawler

- Build HTTP client with user-agent, timeouts, retry, and rate limiting.
- Build `robots.txt` checker.
- Build Top 100 parser.
- Build novel detail parser.
- Add detail-page cache policy.

### Phase 3: Storage

- Upsert novels.
- Upsert ignored novels.
- Upsert tags.
- Store ranking snapshots.
- Store ranking entries.
- Store ignored ranking entries.
- Track crawl runs and failures.

### Phase 4: Analytics

- Compute tag stats for a snapshot.
- Compute daily deltas.
- Compute seven-day weekly aggregates.
- Store materialized stats if useful.

### Phase 5: Reports

- Generate markdown reports.
- Generate JSON reports for later dashboard use.
- Store generated reports in the database.
- Write files under `reports/daily` and `reports/weekly`.

### Phase 6: Operations

- Add scheduled runs.
- Add logging.
- Add crawl health checks.
- Add parser regression tests with saved HTML fixtures.
- Add README usage instructions.

## Open Decisions

- Whether to include the generic `NOVEL` tag in reports. Default: exclude.
- Whether `로맨스판타지` or similar variants should count as `로판` if Naver does not provide the exact `로판` tag. Default: exact match only.
- Whether to report only `categoryCode=ALL` or also per-category reports such as romance and rofan.
- Whether to use Naver weekly rankings directly or only seven-day aggregation. Default: aggregate daily snapshots first.
- Whether reports should be markdown only or markdown plus JSON. Default: both.
- Whether to later build a web dashboard. Default: defer until the report format stabilizes.
