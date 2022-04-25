package repository

import (
	"context"
	"github.com/reaganiwadha/arah/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
	"time"
)

type MongoTestSuite struct {
	suite.Suite
	db *mongo.Database
}

func (s *MongoTestSuite) SetupSuite() {
	mClient, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("TESTING_MONGO_URI")))

	if err != nil {
		s.Error(err)
		return
	}

	connectCtx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err = mClient.Connect(connectCtx)

	s.db = mClient.Database("arah_testing")
}

func (s *MongoTestSuite) TearDownTest() {
	col := s.db.Collection("links")
	err := col.Drop(context.Background())

	require.NoError(s.T(), err)
}

func (s *MongoTestSuite) TearDownSuite() {
	err := s.db.Drop(context.Background())

	require.NoError(s.T(), err)
}

func (s *MongoTestSuite) TestLinkMongoRepo_GetLink_Success() {
	slug := "gh"
	link := "https://github.com"

	col := s.db.Collection("links")

	_, err := col.InsertOne(context.Background(), bson.M{
		"slug": slug,
		"link": link,
	})
	require.NoError(s.T(), err)

	repo, err := NewLinkRepository(s.db)
	require.NoError(s.T(), err)

	res, err := repo.GetLink(context.Background(), slug)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), res.Link, link)
	assert.Equal(s.T(), res.Slug, slug)
}

func (s *MongoTestSuite) TestLinkMongoRepo_GetLink_ErrorOnNoSlug() {
	repo, err := NewLinkRepository(s.db)
	require.NoError(s.T(), err)

	res, err := repo.GetLink(context.Background(), "nonexistent")
	assert.Nil(s.T(), res)
	assert.ErrorIs(s.T(), err, domain.ErrDataNotFound)
}

func (s *MongoTestSuite) TestLinkMongoRepo_CreateLink_Success() {
	repo, err := NewLinkRepository(s.db)

	require.NoError(s.T(), err)

	slug := "gfq"
	link := "https://genericfilter.quest"

	res, err := repo.CreateLink(context.Background(), slug, link)

	require.NoError(s.T(), err)

	assert.Equal(s.T(), res.Link, link)
	assert.Equal(s.T(), res.Slug, slug)

	var result shortenedLinkModel

	col := s.db.Collection("links")
	err = col.FindOne(context.Background(), bson.M{"slug": slug}).Decode(&result)

	require.NoError(s.T(), err)

	assert.Equal(s.T(), result.Link, link)
	assert.Equal(s.T(), result.Slug, slug)
}

func (s *MongoTestSuite) TestLinkMongoRepo_CreateLink_ErrorOnDuplicateSlug() {
	repo, err := NewLinkRepository(s.db)

	require.NoError(s.T(), err)

	slug := "gfq"
	link := "https://genericfilter.quest"

	res, err := repo.CreateLink(context.Background(), slug, link)

	require.NoError(s.T(), err)

	assert.Equal(s.T(), res.Link, link)
	assert.Equal(s.T(), res.Slug, slug)

	res, err = repo.CreateLink(context.Background(), slug, link)
	assert.Nil(s.T(), res)
	assert.ErrorIs(s.T(), err, domain.ErrDataExists)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(MongoTestSuite))
}
