package repository

import (
	"context"
	"github.com/reaganiwadha/arah/domain"
	lr "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const linkCollectionName = "links"

type linkRepository struct {
	db *mongo.Database
}

type shortenedLinkModel struct {
	ID   string `bson:"_id,omitempty"`
	Slug string `bson:"slug"`
	Link string `bson:"link"`
}

func NewLinkRepository(db *mongo.Database) (res domain.LinkRepository, err error) {
	linkCollection := db.Collection(linkCollectionName)

	_, err = linkCollection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{
			"slug": 1,
		},
		Options: options.Index().SetUnique(true),
	})

	if err != nil {
		return nil, err
	}

	return &linkRepository{
		db: db,
	}, nil
}

func (l linkRepository) CreateLink(ctx context.Context, slug string, link string) (res *domain.ShortenedLink, err error) {
	c := l.db.Collection(linkCollectionName)

	insertRes, err := c.InsertOne(ctx, bson.M{
		"slug": slug,
		"link": link,
	})

	if err != nil {
		if we, ok := err.(mongo.WriteException); ok {
			if len(we.WriteErrors) > 0 && we.WriteErrors[0].Code == 11000 {
				return nil, domain.ErrDataExists
			} else {
				lr.Errorf("unknown error : %s", err)
				return nil, domain.ErrUnknown
			}
		}
		return
	}

	res = &domain.ShortenedLink{
		Slug: slug,
		Link: link,
		ID:   insertRes.InsertedID.(primitive.ObjectID).String(),
	}

	return
}

func (l linkRepository) GetLink(ctx context.Context, slug string) (res *domain.ShortenedLink, err error) {
	c := l.db.Collection(linkCollectionName)

	var result shortenedLinkModel

	err = c.FindOne(ctx, bson.M{
		"slug": slug,
	}).Decode(&result)

	if err == mongo.ErrNoDocuments {
		err = domain.ErrDataNotFound
		return
	}

	res = &domain.ShortenedLink{
		ID:   result.ID,
		Slug: result.Slug,
		Link: result.Link,
	}

	return
}
