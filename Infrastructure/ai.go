package infrastructure
import (
    "blog-api/Domain"
    "context"
    "fmt"
    "google.golang.org/genai"
)
type AiService struct{

}
func NewAiService() domain.AiService{
   return &AiService{} 
}
func(as *AiService) Getsuggestion(req domain.AiSuggestionRequest)(string, error){
    ctx := context.Background()
    client, err := genai.NewClient(ctx, &genai.ClientConfig{
        APIKey: Env.API_Key,
        Backend: genai.BackendGeminiAPI,
    })
    if err != nil {
		return "", err
    }

    result, err := client.Models.GenerateContent(
        ctx,
        "gemini-2.5-flash",
		genai.Text(fmt.Sprintf("Generate a Suggestion to the blog i am writing you can give me suggestion in grammer structure the flow of the blog and other aspects that can make my blog great. the title of the blog is %v, the content of the blog is %v, and the tags i am putting on the blog is %v ",req.Title,req.Content, req.Tags)),
        nil,
    )
    if err != nil {
        return "", err
    }
    return result.Text(), nil
}
