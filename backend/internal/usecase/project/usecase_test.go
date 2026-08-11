package project

import (
	"strings"
	"testing"

	"github.com/lib/pq"

	"siyue-blog/internal/entity"
)

func TestNormalizeStackTrimsAndDeduplicates(t *testing.T) {
	result, err := normalizeStack([]string{" Go ", "go", "PostgreSQL"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 || result[0] != "Go" || result[1] != "PostgreSQL" {
		t.Fatalf("unexpected normalized stack: %#v", result)
	}
}

func TestValidateProjectRejectsOversizedURL(t *testing.T) {
	project := &entity.Project{Title: "Project", RepoURL: "https://example.com/" + strings.Repeat("a", 500), TechStack: pq.StringArray{}}
	if err := validateProject(project); err == nil {
		t.Fatal("expected oversized URL validation error")
	}
}
