# Jira-Style Task Backlog

This backlog breaks the Naver Series romance/rofan tag analytics project into implementation-ready tasks.

## Epic: Build Naver Series Romance/Rofan Tag Analytics Pipeline

Build a crawler, database, analytics layer, and report generator that tracks popular Naver Series tags for novels tagged `로판` or `로맨스`.

### SER-1: Initialize TypeScript Project Tooling

- [x] Task complete

Type: Task  
Priority: High  
Labels: `setup`, `typescript`, `tooling`

Description:

Set up the Node.js project for TypeScript development and command-line execution.

Acceptance Criteria:

- TypeScript config exists.
- `tsx` is available for running CLI scripts.
- Basic `src/` directory structure exists.
- Package scripts are added for future crawler/report commands.
- Generated build/cache files are ignored by git.

Dependencies:

- None

### SER-2: Choose and Configure SQLite ORM

- [ ] Task complete

Type: Task  
Priority: High  
Labels: `database`, `sqlite`, `orm`

Description:

Select Drizzle ORM or Prisma and configure SQLite access for local development.

Acceptance Criteria:

- SQLite client can open `data/series.sqlite`.
- ORM schema location is defined.
- Migration command is available.
- Database path is configurable.

Dependencies:

- SER-1

### SER-3: Create Initial Database Schema

- [ ] Task complete

Type: Task  
Priority: High  
Labels: `database`, `schema`, `migration`

Description:

Create database tables needed for included novels, ignored novels, ranking snapshots, ranking entries, tag stats, crawl runs, and reports.

Acceptance Criteria:

- `novels` table exists.
- `tags` table exists.
- `novel_tags` table exists.
- `ignored_novels` table exists.
- `ranking_snapshots` table exists.
- `ranking_entries` table exists.
- `ignored_ranking_entries` table exists.
- `tag_daily_stats` table exists.
- `crawl_runs` table exists.
- `reports` table exists.
- Unique constraints and indexes from the plan are included.

Dependencies:

- SER-2

### SER-4: Implement Configuration Module

- [ ] Task complete

Type: Task  
Priority: Medium  
Labels: `config`

Description:

Create a typed configuration module for crawl URLs, required tags, cache TTLs, request pacing, database path, and output directories.

Acceptance Criteria:

- Required tags default to `로판` and `로맨스`.
- Detail cache TTL is configurable.
- Ignore recheck TTL is configurable.
- Ranking type/category/page count are configurable.
- Report output paths are configurable.

Dependencies:

- SER-1

### SER-5: Implement HTTP Client

- [ ] Task complete

Type: Task  
Priority: High  
Labels: `crawler`, `http`, `reliability`

Description:

Build a reusable HTTP client for fetching Naver Series pages conservatively.

Acceptance Criteria:

- Uses a clear user-agent.
- Supports request timeout.
- Supports retry with backoff.
- Supports rate limiting.
- Returns response body and status.
- Logs failed requests with enough context for debugging.

Dependencies:

- SER-1
- SER-4

### SER-6: Implement Robots Policy Check

- [ ] Task complete

Type: Task  
Priority: High  
Labels: `crawler`, `robots`, `compliance`

Description:

Add a robots policy checker that confirms planned crawl paths are allowed and prevents accidental crawling of disallowed paths.

Acceptance Criteria:

- Fetches or loads Naver Series `robots.txt`.
- Confirms `/novel/top100List.series` is allowed before crawling.
- Blocks `/novel/categoryProductList.series`.
- Fails closed if policy cannot be checked.
- Has tests for allowed and disallowed paths.

Dependencies:

- SER-5

### SER-7: Build Top 100 Ranking Parser

- [ ] Task complete

Type: Task  
Priority: High  
Labels: `crawler`, `parser`, `ranking`

Description:

Parse ranking pages and extract rank, product number, title, and detail URL.

Acceptance Criteria:

- Parses all ranking entries from a saved Top 100 page fixture.
- Extracts `productNo`.
- Extracts rank.
- Extracts title.
- Extracts normalized detail URL.
- Handles pages 1-5.
- Parser has fixture-based tests.

Dependencies:

- SER-1

### SER-8: Build Novel Detail Parser

- [ ] Task complete

Type: Task  
Priority: High  
Labels: `crawler`, `parser`, `tags`

Description:

Parse detail pages and extract novel metadata and tags.

Acceptance Criteria:

- Extracts title.
- Extracts product number.
- Extracts cover URL when available.
- Extracts description/meta description.
- Extracts hashtag tags.
- Removes generic `NOVEL` tag by default.
- Identifies whether tags contain `로판` or `로맨스`.
- Parser has fixture-based tests.

Dependencies:

- SER-1

### SER-9: Implement Included Novel Storage

- [ ] Task complete

Type: Task  
Priority: High  
Labels: `database`, `storage`, `novels`

Description:

Implement database operations for included novels and their tags.

Acceptance Criteria:

- Upserts included novels by `product_no`.
- Updates `last_seen_at`.
- Stores detail fetch timestamp and hash.
- Upserts tags by name.
- Replaces a novel's tag links after detail refresh.
- Removes stale ignore-list row if a previously ignored novel becomes included.

Dependencies:

- SER-3
- SER-8

### SER-10: Implement Ignored Novel Storage

- [ ] Task complete

Type: Task  
Priority: High  
Labels: `database`, `storage`, `ignore-list`

Description:

Implement database operations for ignored novels that do not contain `로판` or `로맨스`.

Acceptance Criteria:

- Upserts ignored novels by `product_no`.
- Stores observed tags as JSON.
- Stores ignore reason.
- Stores `next_recheck_at`.
- Updates `last_seen_at` when the ignored novel appears again.
- Exposes helper to determine whether an ignored novel should be rechecked.

Dependencies:

- SER-3
- SER-8

### SER-11: Implement Ranking Snapshot Storage

- [ ] Task complete

Type: Task  
Priority: High  
Labels: `database`, `ranking`

Description:

Store ranking snapshots, included ranking entries, and ignored ranking entries.

Acceptance Criteria:

- Creates or replaces a snapshot for date/type/category.
- Stores included ranking entries linked to `novels`.
- Stores ignored ranking entries for ignored novels.
- Prevents duplicate rank/product entries within a snapshot.
- Supports rerunning the same date without duplicating data.

Dependencies:

- SER-3
- SER-9
- SER-10

### SER-12: Implement Daily Crawl Command

- [ ] Task complete

Type: Story  
Priority: High  
Labels: `crawler`, `cli`, `daily`

Description:

Create a CLI command that fetches daily ranking pages, classifies novels, stores included novels, stores ignored novels, and records crawl run status.

Acceptance Criteria:

- Command is available as `pnpm crawl:daily`.
- Fetches daily Top 100 pages 1-5.
- Checks robots policy before crawling.
- Fetches details for new/stale/recheck-due novels only.
- Skips fresh included and ignored cache records.
- Writes `crawl_runs` status.
- Stores included and ignored ranking entries.
- Exits non-zero on failed crawl.

Dependencies:

- SER-5
- SER-6
- SER-7
- SER-9
- SER-10
- SER-11

### SER-13: Implement Weekly Crawl Command

- [ ] Task complete

Type: Story  
Priority: Medium  
Labels: `crawler`, `cli`, `weekly`

Description:

Create a CLI command that fetches Naver weekly ranking pages and stores them using the same inclusion/ignore rules.

Acceptance Criteria:

- Command is available as `pnpm crawl:weekly`.
- Fetches weekly Top 100 pages 1-5.
- Uses the same scope filter as daily crawl.
- Stores snapshot with `ranking_type = WEEKLY`.
- Uses detail and ignore caches.

Dependencies:

- SER-12

### SER-14: Compute Daily Tag Stats

- [ ] Task complete

Type: Task  
Priority: High  
Labels: `analytics`, `tags`, `daily`

Description:

Compute tag statistics for one daily snapshot using only included novels.

Acceptance Criteria:

- Computes novel count per tag.
- Computes weighted score using `101 - rank`.
- Computes average rank.
- Computes best and worst rank.
- Computes rank bucket counts.
- Excludes ignored novels.
- Can persist results to `tag_daily_stats`.

Dependencies:

- SER-11

### SER-15: Compute Daily Trend Deltas

- [ ] Task complete

Type: Task  
Priority: Medium  
Labels: `analytics`, `trends`, `daily`

Description:

Compare today's included tag stats with previous daily snapshots.

Acceptance Criteria:

- Computes count delta vs yesterday.
- Computes score delta vs yesterday.
- Identifies new included tags.
- Identifies dropped included tags.
- Handles missing previous snapshot gracefully.

Dependencies:

- SER-14

### SER-16: Compute Weekly Aggregates

- [ ] Task complete

Type: Task  
Priority: Medium  
Labels: `analytics`, `weekly`

Description:

Aggregate seven daily snapshots into weekly tag metrics.

Acceptance Criteria:

- Aggregates included tag appearances across seven days.
- Aggregates weighted score across seven days.
- Computes days present.
- Computes weekly average count.
- Computes week-over-week changes when prior data exists.
- Excludes ignored novels.

Dependencies:

- SER-14
- SER-15

### SER-17: Generate Daily Markdown Report

- [ ] Task complete

Type: Story  
Priority: High  
Labels: `reports`, `markdown`, `daily`

Description:

Generate a daily markdown report matching the approved example report structure.

Acceptance Criteria:

- Command is available as `pnpm report:daily`.
- Writes report under `reports/daily/`.
- Includes crawl summary.
- Includes scope filter summary.
- Includes top included tags by count.
- Includes top included tags by weighted score.
- Includes risers and fallers.
- Includes new included tags.
- Includes top included novels.
- Includes ignored novel summary.
- Includes data quality notes.
- Stores report content in `reports` table.

Dependencies:

- SER-14
- SER-15

### SER-18: Generate Weekly Markdown Report

- [ ] Task complete

Type: Story  
Priority: Medium  
Labels: `reports`, `markdown`, `weekly`

Description:

Generate a weekly markdown report from seven daily snapshots.

Acceptance Criteria:

- Command is available as `pnpm report:weekly`.
- Writes report under `reports/weekly/`.
- Includes weekly scope filter summary.
- Includes weekly top included tags.
- Includes week-over-week risers and fallers.
- Includes consistent tags by days present.
- Includes top included novels contributing to weekly tag strength.
- Includes ignored novel summary.
- Stores report content in `reports` table.

Dependencies:

- SER-16

### SER-19: Generate JSON Report Output

- [ ] Task complete

Type: Task  
Priority: Low  
Labels: `reports`, `json`

Description:

Generate JSON report files alongside markdown reports for future dashboard or automation use.

Acceptance Criteria:

- Daily report can be emitted as JSON.
- Weekly report can be emitted as JSON.
- JSON schema is documented.
- JSON content matches markdown report data.

Dependencies:

- SER-17
- SER-18

### SER-20: Add Parser Fixtures and Tests

- [ ] Task complete

Type: Task  
Priority: High  
Labels: `tests`, `fixtures`, `parser`

Description:

Add saved HTML fixtures and automated tests for ranking and detail parsers.

Acceptance Criteria:

- Ranking fixture exists.
- Included novel detail fixture exists.
- Ignored novel detail fixture exists.
- Ranking parser test verifies extracted entries.
- Detail parser test verifies extracted tags.
- Inclusion filter test verifies `로판`/`로맨스`.
- Ignore filter test verifies non-matching tags.

Dependencies:

- SER-7
- SER-8

### SER-21: Add Storage Tests

- [ ] Task complete

Type: Task  
Priority: Medium  
Labels: `tests`, `database`

Description:

Add tests around database upsert behavior and rerunnable snapshots.

Acceptance Criteria:

- Included novel upsert test exists.
- Ignored novel upsert test exists.
- Ranking snapshot rerun test exists.
- Tag replacement test exists.
- Ignore recheck TTL test exists.

Dependencies:

- SER-9
- SER-10
- SER-11

### SER-22: Add Crawl Run Logging

- [ ] Task complete

Type: Task  
Priority: Medium  
Labels: `observability`, `logging`

Description:

Add structured logging and crawl run metrics for operational visibility.

Acceptance Criteria:

- Logs crawl start and finish.
- Logs pages fetched.
- Logs detail fetch count.
- Logs included vs ignored counts.
- Logs failures with URL and reason.
- `crawl_runs` is updated on success, partial success, and failure.

Dependencies:

- SER-12

### SER-23: Add Scheduling Documentation

- [ ] Task complete

Type: Task  
Priority: Low  
Labels: `docs`, `operations`

Description:

Document how to schedule daily and weekly crawler/report commands.

Acceptance Criteria:

- README includes cron examples.
- Timezone recommendation is documented.
- Manual run instructions are documented.
- Failure recovery instructions are documented.

Dependencies:

- SER-12
- SER-17

### SER-24: Update README After Implementation

- [ ] Task complete

Type: Task  
Priority: Low  
Labels: `docs`

Description:

Update the README once commands are implemented.

Acceptance Criteria:

- Removes "not implemented yet" note.
- Documents setup steps.
- Documents migration command.
- Documents crawl commands.
- Documents report commands.
- Documents output files and database location.

Dependencies:

- SER-17
- SER-18
