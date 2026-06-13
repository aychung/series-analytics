package store

import (
	"context"
	"path/filepath"
	"testing"

	"series-analytics/internal/model"
)

func TestStoreNovelTagsReusesExistingTags(t *testing.T) {
	ctx := context.Background()

	s, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	if err := RunMigration(s.DB, filepath.Join("..", "..", "migration")); err != nil {
		t.Fatal(err)
	}

	novelID, err := s.StoreNovelDetail(ctx, model.NovelDetail{
		PrdNo:     "1",
		Title:     "Title",
		Author:    "Author",
		Publisher: "Publisher",
		Category:  "Category",
	})
	if err != nil {
		t.Fatal(err)
	}

	tags := []string{"fantasy", "romance"}
	if err := s.StoreNovelTags(ctx, novelID, tags); err != nil {
		t.Fatal(err)
	}

	if err := s.StoreNovelTags(ctx, novelID, tags); err != nil {
		t.Fatal(err)
	}

	var tagCount int
	if err := s.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM tags").Scan(&tagCount); err != nil {
		t.Fatal(err)
	}
	if tagCount != len(tags) {
		t.Fatalf("tag count = %d, want %d", tagCount, len(tags))
	}

	var novelTagCount int
	if err := s.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM novel_tags").Scan(&novelTagCount); err != nil {
		t.Fatal(err)
	}
	if novelTagCount != len(tags) {
		t.Fatalf("novel tag count = %d, want %d", novelTagCount, len(tags))
	}
}
