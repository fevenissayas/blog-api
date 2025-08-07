package domain

type AiSuggestionRequest struct{
	Title string
	Content string
    Tags    []string
}
type AiService interface {
	Getsuggestion(req AiSuggestionRequest)(string, error)
}
