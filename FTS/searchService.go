package fts

import (
	"log"
	"os"
	"something/model"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/v2/analysis/lang/ru"
	"github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	"github.com/blevesearch/bleve/v2/analysis/tokenizer/unicode"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/search/query"
	"github.com/google/uuid"
)

type SearchService struct {
	index bleve.Index
}

func NewSearchService(indexPath string) (*SearchService, error) {
    var index bleve.Index

    if _, err := os.Stat(indexPath); os.IsNotExist(err) {
        mapping := createIndexMapping()
        index, err = bleve.New(indexPath, mapping)
        if err != nil {
            return nil, err
        }
    } else {
        index, err = bleve.Open(indexPath)
        if err != nil {
            return nil, err
        }
    }

    return &SearchService{index: index}, nil
}

func createIndexMapping() *mapping.IndexMappingImpl {
	mapping := bleve.NewIndexMapping()
	err := mapping.AddCustomAnalyzer("ru_en", map[string]interface{}{
		"type":      custom.Name,
		"tokenizer": unicode.Name,
		"token_filters": []string{
			lowercase.Name,
			ru.StopName,
			ru.SnowballStemmerName,
		},
	})
	if err != nil {
		log.Fatalf("Failed to add custom analyzer: %v", err)
	}

    textFieldMapping := bleve.NewTextFieldMapping()
    textFieldMapping.Analyzer = "ru_en"

    keywordFieldMapping := bleve.NewTextFieldMapping()
    keywordFieldMapping.Analyzer = "keyword"

    userIDFieldMapping := bleve.NewTextFieldMapping()
    userIDFieldMapping.Analyzer = "keyword"

    dateFieldMapping := bleve.NewDateTimeFieldMapping()

    bookmarkMapping := bleve.NewDocumentMapping()
    
	bookmarkMapping.AddFieldMappingsAt("title", textFieldMapping)
	bookmarkMapping.AddFieldMappingsAt("content", textFieldMapping)
	bookmarkMapping.AddFieldMappingsAt("description", textFieldMapping)
	bookmarkMapping.AddFieldMappingsAt("url", keywordFieldMapping)
	bookmarkMapping.AddFieldMappingsAt("tag_ids", keywordFieldMapping)
	bookmarkMapping.AddFieldMappingsAt("tag_names", keywordFieldMapping)
	bookmarkMapping.AddFieldMappingsAt("user_id", userIDFieldMapping)
	bookmarkMapping.AddFieldMappingsAt("created_at", dateFieldMapping)

    indexMapping := bleve.NewIndexMapping()
    indexMapping.DefaultMapping = bookmarkMapping
    indexMapping.DefaultAnalyzer = "ru_en"

    return indexMapping
}

func (s *SearchService) BookmarkSearch(userID uuid.UUID, query string, limit int) ([]*model.Bookmark, error) {
	userQuery := bleve.NewTermQuery(userID.String())
	userQuery.SetField("user_id")

	var searchRequest *bleve.SearchRequest
	if query != "" {
		textQuery := bleve.NewDisjunctionQuery(
			makeFieldQuery("title", query, 2.0),
			makeFieldQuery("description", query, 1.0),
		)

		searchRequest = bleve.NewSearchRequest(bleve.NewConjunctionQuery(userQuery, textQuery))
	} else {
		searchRequest = bleve.NewSearchRequest(userQuery)
	}

	searchRequest.Size = limit
	searchRequest.Fields = []string{"title", "description", "url", "created_at"}

	searchRequest.SortBy([]string{"-created_at"})

	result, err := s.index.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	var bookmarks []*model.Bookmark

	for _, hit := range result.Hits {
		var bookmark model.Bookmark
		fields := hit.Fields

		bookmark.ID, _ = uuid.Parse(hit.ID)
		bookmark.URL = fields["url"].(string)
		bookmark.Title = fields["title"].(string)
		bookmark.Description = fields["description"].(string)
		tagIDs := fields["tag_ids"].([]interface{})
		tagNames := fields["tag_names"].([]interface{})
		for i := range tagIDs {
			id, _ := uuid.Parse(tagIDs[i].(string))
			bookmark.Tags = append(bookmark.Tags, model.Tag{ID: id, Name: tagNames[i].(string)})
		}
		bookmark.UserID, _ = uuid.Parse(fields["user_id"].(string))
		bookmark.CreatedAt, _ = time.Parse(time.RFC3339, fields["created_at"].(string))

		bookmarks = append(bookmarks, &bookmark)
	}

	return bookmarks, nil
}

func (s *SearchService) NoteSearch(userID uuid.UUID, query string, limit int) ([]*model.Note, error) {
	userQuery := bleve.NewTermQuery(userID.String())
	userQuery.SetField("user_id")

	var searchRequest *bleve.SearchRequest
	if query != "" {
		textQuery := bleve.NewDisjunctionQuery(
			makeFieldQuery("title", query, 2.0),
			makeFieldQuery("content", query, 1.0),
		)

		searchRequest = bleve.NewSearchRequest(bleve.NewConjunctionQuery(userQuery, textQuery))
	} else {
		searchRequest = bleve.NewSearchRequest(userQuery)
	}

	searchRequest.Size = limit
	searchRequest.Fields = []string{"title", "content", "created_at"}

	searchRequest.SortBy([]string{"-created_at"})

	result, err := s.index.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	var notes []*model.Note

	for _, hit := range result.Hits {
		var note model.Note
		fields := hit.Fields

		note.ID, _ = uuid.Parse(hit.ID)
		note.Title = fields["title"].(string)
		note.Content = fields["content"].(string)
		tagIDs := fields["tag_ids"].([]interface{})
		tagNames := fields["tag_names"].([]interface{})
		for i := range tagIDs {
			id, _ := uuid.Parse(tagIDs[i].(string))
			note.Tags = append(note.Tags, model.Tag{ID: id, Name: tagNames[i].(string)})
		}
		note.UserID, _ = uuid.Parse(fields["user_id"].(string))
		note.CreatedAt, _ = time.Parse(time.RFC3339, fields["created_at"].(string))

		notes = append(notes, &note)
	}

	return notes, nil
}

type indexBookmark struct {
	ID          string    `json:"id"`
	URL         string    `json:"url"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	TagNames   	[]string  `json:"tag_names"`
	TagIds	  	[]string  `json:"tag_ids"`
	UserID      string    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}

func (s *SearchService) IndexBookmark(bookmark model.Bookmark) error {
	return s.index.Index(bookmark.ID.String(), indexBookmark{
		ID:          bookmark.ID.String(),
		URL:         bookmark.URL,
		Title:       bookmark.Title,
		Description: bookmark.Description,
		TagNames:    bookmark.Tags.Names(),
		TagIds: 	 bookmark.Tags.IDs(),
		UserID:      bookmark.UserID.String(),
		CreatedAt:   bookmark.CreatedAt,
	})
}

type indexNote struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	TagNames  []string  `json:"tag_names"`
	TagIds	  []string  `json:"tag_ids"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *SearchService) IndexNote(note model.Note) error {
	return s.index.Index(note.ID.String(), indexNote{
		ID:        note.ID.String(),
		Title:     note.Title,
		Content:   note.Content,
		TagNames:  note.Tags.Names(),
		TagIds:    note.Tags.IDs(),
		UserID:    note.UserID.String(),
		CreatedAt: note.CreatedAt,
	})
}

func makeFieldQuery(field, query string, boost float64) *query.MatchQuery {
    q := bleve.NewMatchQuery(query)
    q.SetField(field)
    q.SetBoost(boost)
    return q
}