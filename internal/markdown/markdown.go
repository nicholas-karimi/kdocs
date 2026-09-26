package markdown

import (
	"bytes"

	"github.com/yuin/goldmark"
)

type Renderer struct {
	markdown goldmark.Markdown
}

func NewRenderer() *Renderer {
	md := goldmark.New()

	return &Renderer{
		markdown: md,
	}
}

func (r *Renderer) Render(source string) (string, error) {
	var buf bytes.Buffer

	err := r.markdown.Convert([]byte(source), &buf)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
