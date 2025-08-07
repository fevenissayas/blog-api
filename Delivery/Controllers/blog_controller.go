package controllers

import (
	domain "blog-api/Domain"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type createBlogRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type updateBlogRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type GetBlogsRequest struct {
    Page     int    `form:"page" binding:"min=1"`
    Limit    int    `form:"limit" binding:"min=1,max=100"`
    Sort     string `form:"sort"` // "recent", "popular", "views"
    AuthorID string `form:"author_id"`
}

type PaginatedBlogsResponse struct {
    Blogs      []domain.Blog `json:"blogs"`
    Page       int           `json:"page"`
    Limit      int           `json:"limit"`
    Total      int64         `json:"total"`
    TotalPages int           `json:"total_pages"`
    HasNext    bool          `json:"has_next"`
    HasPrev    bool          `json:"has_prev"`
}


type BlogController struct {
	blogUsecase domain.IBlogUsecase
}
func NewBlogController(blogUsecase domain.IBlogUsecase) *BlogController {
	return &BlogController{blogUsecase: blogUsecase}
}

func (bc *BlogController) CreateBlogHandler(ctx *gin.Context) {
	var req createBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}

	userID, ok := getAuthenticatedUserID(ctx)
	if !ok {
		return
	}

	now := time.Now()
	blog := &domain.Blog{
		Title:     req.Title,
		Content:   req.Content,
		Tags:      req.Tags,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
		ViewCount: 0,
	}

	if err := bc.blogUsecase.Create(ctx, blog); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Blog created successfully"})
}

func (bc *BlogController) UpdateBlogHandler(ctx *gin.Context) {

	blogID := ctx.Param("id")

	userID, ok := getAuthenticatedUserID(ctx)
	if !ok {
		return
	}
	var req updateBlogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	input := domain.UpdateBlogInput{
		BlogID:  blogID,
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
		Tags:    req.Tags,
	}

	updatedBlog, err := bc.blogUsecase.Update(ctx.Request.Context(), input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedBlog)
}

func (bc *BlogController) FilterBlogsHandler(ctx *gin.Context) {
	tag := ctx.Query("tag")
	date := ctx.Query("date")
	sort := ctx.DefaultQuery("sort", "recent")

	blogs, err := bc.blogUsecase.FilterBlogs(ctx.Request.Context(), tag, date, sort)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, blogs)
}

func (bc *BlogController) DeleteBlog(ctx *gin.Context) {
	blogID := ctx.Param("id")
	userID, ok := getAuthenticatedUserID(ctx)
	if !ok {
		return
	}
	userRole := ctx.GetString("role")
	if err := bc.blogUsecase.DeleteBlog(ctx, blogID, userID, userRole); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Blog deleted successfully"})

}

func (bc *BlogController) SearchBlogs(ctx *gin.Context) {
    tag := ctx.Query("tag")
    date := ctx.Query("date")
    sort := ctx.Query("sort")
    title := ctx.Query("title")
    userID := ctx.Query("userID")

    blogs, err := bc.blogUsecase.SearchBlogs(ctx.Request.Context(), tag, date, sort, title, userID)
    if err != nil {
        ctx.JSON(500, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(200, blogs)
}
func (bc *BlogController) AiSuggestion(ctx *gin.Context){
   var req updateBlogRequest
   if err := ctx.ShouldBindJSON(&req); err != nil{
	   ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error":"invalid request"})
	   return
   }
   input := domain.AiSuggestionRequest{
	   Title: req.Title,
	   Content: req.Content,
	   Tags: req.Tags,
   }
   text, err := bc.blogUsecase.GetSuggestion(input)
   if err != nil {
	   ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error":err.Error()})
	   return
   }
   ctx.IndentedJSON(http.StatusOK, gin.H{"Suggestion":text})
}


// handler for getting paginated blogs
func (bc *BlogController) GetBlogsHandler(ctx *gin.Context) {
    var req GetBlogsRequest
    
    // Set defaults
    req.Page = 1
    req.Limit = 10
    req.Sort = "recent"
    
    if err := ctx.ShouldBindQuery(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
        return
    }

    result, err := bc.blogUsecase.GetBlogs(ctx.Request.Context(), req.Page, req.Limit, req.Sort, req.AuthorID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, result)
}

// Handler for getting single blog (with view increment)
func (bc *BlogController) GetBlogByIDHandler(ctx *gin.Context) {
    blogID := ctx.Param("id")
    
    blog, err := bc.blogUsecase.GetByIDAndIncrementViews(ctx.Request.Context(), blogID)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
        return
    }
    
    ctx.JSON(http.StatusOK, blog)
}
