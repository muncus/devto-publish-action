// Types used by the dev.to json api.

package devto

// Article represents a published or unpublished article.
// Note when creating new articles, the BodyMarkdown may contain frontmatter
// that takes precendence over values included here.
type Article struct {
	Title        string `json:"title"`
	TypeOf       string `json:"type_of"`
	ID           int    `json:"id"`
	Description  string `json:"description"`
	Published    bool   `json:"published"`
	BodyMarkdown string `json:"body_markdown"`
}
