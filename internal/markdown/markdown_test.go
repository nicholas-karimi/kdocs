package markdown

import (
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	renderer := NewRenderer()

	input := "# Hello KDocs"

	output, err := renderer.Render(input)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if !strings.Contains(output, "<h1>Hello KDocs</h1>") {
		t.Fatalf("Render() output = %q", output)
	}
}

func TestRenderCommonMarkdown(t *testing.T) {
	renderer := NewRenderer()

	input := `# PostgreSQL Migration

This is **important**.

- Backup the database
- Run the migration
- Verify the application

## Commands

` + "```bash\nmigrate -path internal/database/migrations -database \"$KDOCS_DB_DSN\" up\n```" + `

`

	output, err := renderer.Render(input)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	expected := []string{
		"<h1>PostgreSQL Migration</h1>",
		"<strong>important</strong>",
		"<li>Backup the database</li>",
		"<li>Run the migration</li>",
		"<li>Verify the application</li>",
		"<h2>Commands</h2>",
		"<code",
	}

	for _, value := range expected {
		if !strings.Contains(output, value) {
			t.Errorf("Render() output does not contain %q\nOutput:\n%s", value, output)
		}
	}
}
