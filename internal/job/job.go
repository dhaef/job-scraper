package job

import (
	"encoding/json"

	"github.com/dhaef/job-scraper/internal/dom"
	"github.com/tdewolff/minify/v2"
	h "github.com/tdewolff/minify/v2/html"
	"golang.org/x/net/html"
)

type Job struct {
	Title          string            `json:"title"`
	Link           string            `json:"link"`
	Description    []*html.Node      `json:"-"`
	DescriptionStr string            `json:"description"`
	Metadata       map[string]string `json:"metadata"`
}

func NewJob() *Job {
	return &Job{
		Metadata: map[string]string{},
	}
}

func (j *Job) SetTitle(title string) {
	j.Title = title
}

func (j *Job) SetLink(link string) {
	j.Link = link
}

func (j *Job) SetDescription(description []*html.Node) {
	j.Description = description
}

func (j *Job) SetMetadata(key string, value string) {
	j.Metadata[key] = value
}

func (j *Job) descriptionNodesToHTMLString() (string, error) {
	description := "<div>"

	for _, node := range j.Description {
		node = dom.CleanNodes(node)

		s, err := dom.RenderNode(node)
		if err != nil {
			return "", err
		}

		description += s
	}

	description += "</div>"

	m := minify.New()
	m.AddFunc("text/html", h.Minify)
	return m.String("text/html", description)
}

func (j *Job) ToJSON() (string, error) {
	description, err := j.descriptionNodesToHTMLString()
	if err != nil {
		return "", err
	}

	j.DescriptionStr = description

	jsonData, err := json.Marshal(j)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
