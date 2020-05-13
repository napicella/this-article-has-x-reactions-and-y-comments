package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"

	"net/http"
)

const (
	devToEndpoint = "https://dev.to/api"
)

var securityProviderConfig = struct {
	method     string
	headerName string
}{
	method:     "header",
	headerName: "api-key",
}

type programEnv struct {
	DevToApiKey string `required:"true" split_words:"true"`
	ArticleID   int32  `required:"true" split_words:"true"`
}

func main() {
	log.SetLevel(log.DEBUG)

	env, e := loadEnvVariables()
	if e != nil {
		log.Fatal(e)
	}
	client, e := getHttpClient(env.DevToApiKey)
	if e != nil {
		log.Fatal("Unable to create the http client")
	}
	run(context.Background(), env, client)
}

func run(ctx context.Context, env *programEnv, client *ClientWithResponses) {
	r, e := getArticleActivity(ctx, client, env.ArticleID)
	if e != nil {
		log.Fatal(e)
	}
	log.Debugf("Article reactions: %[1]d - Article comments: %[2]d",
		r.ReactionsCount,
		r.CommentsCount)

	newTitle := generateNewTitle(r.ReactionsCount, r.CommentsCount)
	fmEditor, e := newFrontMatterEditor(r.BodyMarkdown)
	if e != nil {
		log.Fatal(e)
	}
	if fmEditor.title == newTitle {
		log.Debug("Article title is already up to date with the number of reactions and comments")
		return
	}

	log.Debug("Updating the title of the article")
	fmEditor.updateTitle(newTitle)
	e = saveArticle(ctx, client, env.ArticleID, fmEditor.markdown)
	if e != nil {
		log.Fatal(e)
	}
	log.Info("Article with the new title has been saved")
}

func loadEnvVariables() (*programEnv, error) {
	var env programEnv
	e := envconfig.Process("", &env)
	if e != nil {
		return nil, e
	}
	return &env, e
}

func getHttpClient(apiKey string) (*ClientWithResponses, error) {
	apiKeyProvider, e := securityprovider.NewSecurityProviderApiKey(
		securityProviderConfig.method,
		securityProviderConfig.headerName,
		apiKey)
	if e != nil {
		return nil, errors.Wrap(e, "Unable to create the security provider")
	}
	client, e := NewClientWithResponses(devToEndpoint,
		WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			return apiKeyProvider.Intercept(ctx, req)
		}),
	)
	return client, nil
}

func generateNewTitle(reactions, comments int) string {
	return fmt.Sprintf(
		"My personal blog has %[1]d reactions and %[2]d comments",
		reactions,
		comments)
}

func getArticleActivity(ctx context.Context, client *ClientWithResponses, articleId int32) (
	*article, error) {

	r, e := client.GetArticleByIdWithResponse(ctx, articleId)
	if e != nil {
		return nil, e
	}
	var activity article
	e = json.Unmarshal(r.Body, &activity)
	if e != nil {
		return nil, e
	}
	return &activity, nil
}

func saveArticle(ctx context.Context, client *ClientWithResponses, articleID int32, markdown string) error {
	articleUpdate := new(ArticleUpdate)
	articleUpdate.Article.BodyMarkdown = &markdown
	updateResponse, e := client.UpdateArticleWithResponse(ctx, articleID, UpdateArticleJSONRequestBody(*articleUpdate))
	if e != nil {
		return e
	}
	if updateResponse.StatusCode() != http.StatusOK {
		e := errors.Errorf("Update article failed, status code: %d", updateResponse.StatusCode())
		if len(updateResponse.Body) != 0 {
			return errors.Wrap(e, string(updateResponse.Body))
		}
		return e
	}
	return nil
}

type article struct {
	ReactionsCount int    `json:"positive_reactions_count"`
	CommentsCount  int    `json:"comments_count"`
	Title          string `json:"title"`
	BodyMarkdown   string `json:"body_markdown"`
}
