package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrpostgres/builder"
	"github.com/mondegor/go-storage/mrtests/infra"
	"github.com/mondegor/go-webcore/mrtests/helpers"
	"github.com/mondegor/go-webcore/mrtype"
	"github.com/stretchr/testify/suite"

	"github.com/mondegor/go-sample/internal/catalog/category/section/pub/entity"
	"github.com/mondegor/go-sample/internal/catalog/category/section/pub/repository"
	"github.com/mondegor/go-sample/tests"
)

type CategoryPostgresTestSuite struct {
	suite.Suite

	ctx  context.Context
	pgt  *infra.PostgresTester
	repo *repository.CategoryPostgres
}

func TestCategoryPostgresTestSuite(t *testing.T) {
	suite.Run(t, new(CategoryPostgresTestSuite))
}

func (ts *CategoryPostgresTestSuite) SetupSuite() {
	ts.ctx = helpers.ContextWithNopLogger()
	ts.pgt = infra.NewPostgresTester(ts.T(), tests.DBSchemas(), tests.ExcludedDBTables())
	ts.pgt.ApplyMigrations(tests.AppMigrationsDir())

	ts.repo = repository.NewCategoryPostgres(
		ts.pgt.ConnManager(),
		builder.NewSQL(
			builder.WithSQLOrderByDefaultSort(mrtype.SortParams{}),
			builder.WithSQLLimitMaxSize(100),
		),
	)
}

func (ts *CategoryPostgresTestSuite) TearDownSuite() {
	ts.pgt.Destroy(ts.ctx)
}

func (ts *CategoryPostgresTestSuite) SetupTest() {
	ts.pgt.TruncateTables(ts.ctx)
}

func (ts *CategoryPostgresTestSuite) Test_Fetch() {
	ts.pgt.ApplyFixtures("testdata/Fetch")

	expected := []entity.Category{
		{
			ID:      uuid.MustParse("166a72b5-b9fa-499c-8140-3627b34acbbe"),
			Caption: "Бытовая техника",
		},
		{
			ID:      uuid.MustParse("386555ab-9320-4680-b62d-1ea449550fff"),
			Caption: "Электроника",
		},
	}

	ctx := context.Background()
	got, gotCount, err := ts.repo.FetchWithTotal(ctx, entity.CategoryParams{})

	ts.Require().NoError(err)
	ts.Equal(expected, got)
	ts.Equal(uint64(len(expected)), gotCount)
}
