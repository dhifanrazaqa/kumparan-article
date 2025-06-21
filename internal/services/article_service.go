package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dhifanrazaqa/kumparan-article/internal/models"
	"github.com/dhifanrazaqa/kumparan-article/internal/repositories"
	"github.com/go-redis/redis"
	"golang.org/x/sync/errgroup"
)

type ArticleService interface {
	CreateArticle(ctx context.Context, req models.CreateArticleRequest, authorID string) (*models.Article, error)
	GetArticles(ctx context.Context, params models.ListArticlesParams) (*models.PaginatedArticles, error)
	GetArticleByID(ctx context.Context, id string) (*models.Article, error)
	UpdateArticle(ctx context.Context, id string, req models.UpdateArticleRequest, currentUserID string) (*models.Article, error)
	DeleteArticle(ctx context.Context, id string, currentUserID string) error
}

type articleService struct {
	repo        repositories.ArticleRepository
	redisClient *redis.Client
}

func NewArticleService(repo repositories.ArticleRepository, redisClient *redis.Client) ArticleService {
	return &articleService{repo: repo, redisClient: redisClient}
}

func (s *articleService) CreateArticle(ctx context.Context, req models.CreateArticleRequest, authorID string) (*models.Article, error) {
	article := &models.Article{
		Title:    req.Title,
		Body:     req.Body,
		AuthorID: authorID,
	}
	if err := s.repo.Create(ctx, article); err != nil {
		return nil, err
	}
	s.clearArticleCache()
	return article, nil
}

func (s *articleService) GetArticles(ctx context.Context, params models.ListArticlesParams) (*models.PaginatedArticles, error) {
	g, ctx := errgroup.WithContext(ctx)

	var articles []models.Article
	var total int64

	g.Go(func() error {
		var err error
		articles, err = s.repo.FindAll(ctx, params)
		return err
	})
	g.Go(func() error {
		var err error
		total, err = s.repo.CountAll(ctx, params)
		return err
	})
	
	if err := g.Wait(); err != nil {
		return nil, err
	}
	
	totalPages := 0
	if total > 0 && params.Limit > 0 {
		totalPages = int((total + int64(params.Limit) - 1) / int64(params.Limit))
	}
	
	currentPage := 1
	if params.Limit > 0 {
		currentPage = (params.Offset / params.Limit) + 1
	}
	
	fmt.Println(articles)
	return &models.PaginatedArticles{
		Data:       articles,
		Total:      total,
		Page:       currentPage,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

func (s *articleService) GetArticleByID(ctx context.Context, id string) (*models.Article, error) {
	cacheKey := "article:" + id

	val, err := s.redisClient.Get(cacheKey).Result()
	if err == nil {
		var article models.Article
		if json.Unmarshal([]byte(val), &article) == nil {
			log.Printf("Cache HIT for article ID: %s", id)
			return &article, nil
		}
	}

	log.Printf("Cache MISS for article ID: %s", id)
	article, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	jsonData, _ := json.Marshal(article)
	s.redisClient.Set(cacheKey, jsonData, 5*time.Minute)

	return article, nil
}

func (s *articleService) UpdateArticle(ctx context.Context, id string, req models.UpdateArticleRequest, currentUserID string) (*models.Article, error) {
	article, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if article.AuthorID != currentUserID {
		return nil, ErrForbidden
	}

	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Body != "" {
		article.Body = req.Body
	}

	if err := s.repo.Update(ctx, article); err != nil {
		return nil, err
	}

	s.clearArticleCache()
	return article, nil
}

func (s *articleService) DeleteArticle(ctx context.Context, id string, currentUserID string) error {
	article, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if article.AuthorID != currentUserID {
		return ErrForbidden
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.clearArticleCache()
	return nil
}

func (s *articleService) clearArticleCache() {
	iter := s.redisClient.Scan(0, "article:*", 0).Iterator()
	for iter.Next() {
		s.redisClient.Del(iter.Val())
	}
}
