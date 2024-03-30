package store

import (
	"context"
	"fmt"
	"github.com/JohnKucharsky/DevGroup/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"strings"
)

type NewsStore struct {
	db *pgxpool.Pool
}

func NewNewsStore(db *pgxpool.Pool) *NewsStore {
	return &NewsStore{
		db: db,
	}
}

func (ns *NewsStore) BulkInsertCategories(newsID int, categories []int) error {
	ctx := context.Background()

	createParams := pgx.NamedArgs{
		"news_id": newsID,
	}

	var valuesStringArr []string

	for idx, category := range categories {
		catString := strconv.Itoa(category)
		idxString := strconv.Itoa(idx)

		valuesStringArr = append(valuesStringArr, fmt.Sprintf("(@news_id, @%d)", idx))
		createParams[idxString] = catString
	}

	sql := fmt.Sprintf(`
		insert into news_categories (news_id, category_id)
		values %s `, strings.Join(valuesStringArr, ", "),
	)

	_, err := ns.db.Exec(ctx, sql, createParams)
	if err != nil {
		return err
	}

	return nil
}

func (ns *NewsStore) GetCategoriesToNews(id int) ([]int, error) {
	ctx := context.Background()

	rows, err := ns.db.Query(
		ctx, `select category_id from news_categories where news_id = @id`, pgx.NamedArgs{"id": id},
	)
	if err != nil {
		return nil, err
	}

	type Category struct {
		ID int `db:"category_id"`
	}

	res, err := pgx.CollectRows(
		rows, pgx.RowToAddrOfStructByName[Category],
	)
	if err != nil {
		return nil, err
	}

	var categoriesIDs []int

	for _, category := range res {
		categoriesIDs = append(categoriesIDs, category.ID)
	}

	return categoriesIDs, nil
}

func (ns *NewsStore) Create(m domain.NewsInput) (
	*domain.News,
	error,
) {
	ctx := context.Background()

	rows, err := ns.db.Query(
		ctx, `
        INSERT INTO news (title, content)
        VALUES (@title, @content)
        RETURNING id, title, content, updated_at`,
		pgx.NamedArgs{
			"title":   m.Title,
			"content": m.Content,
		},
	)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToAddrOfStructByName[domain.NewsDB],
	)
	if err != nil {
		return nil, err
	}

	if m.Categories != nil {
		if err := ns.BulkInsertCategories(res.ID, *m.Categories); err != nil {
			return nil, err
		}
	}

	ids, err := ns.GetCategoriesToNews(res.ID)

	if err != nil {
		return nil, err
	}
	return domain.NewsDBtoNews(res, ids), nil
}

func (ns *NewsStore) GetManyPaginated(pp *domain.ParsedPaginationParams) ([]*domain.News, *domain.Pagination, error) {
	ctx := context.Background()

	limit := 20
	offset := 0

	var pagination domain.Pagination

	if pp != nil {
		limit = pp.Limit
		pagination.Limit = pp.Limit
		if pp.Offset != nil {
			offset = *pp.Offset
			pagination.Offset = *pp.Offset
		}
	}

	rows, err := ns.db.Query(
		ctx, `select * from news limit @limit offset @offset`, pgx.NamedArgs{"limit": limit, "offset": offset},
	)
	if err != nil {
		return nil, nil, err
	}

	res, err := pgx.CollectRows(
		rows, pgx.RowToAddrOfStructByName[domain.NewsDB],
	)
	if err != nil {
		return nil, nil, err
	}

	type count struct {
		Count int `db:"count"`
	}

	row, err := ns.db.Query(
		ctx, `select count(*) as count from news`,
	)
	total, err := pgx.CollectExactlyOneRow(
		row, pgx.RowToAddrOfStructByName[count],
	)
	if err != nil {
		return nil, nil, err
	}
	pagination.Total = total.Count

	var finalArr []*domain.News

	for _, news := range res {
		ids, err := ns.GetCategoriesToNews(news.ID)

		if err != nil {
			return nil, nil, err
		}

		finalArr = append(finalArr, domain.NewsDBtoNews(news, ids))
	}

	return finalArr, &pagination, nil
}

func (ns *NewsStore) Update(m domain.NewsInputUpdate, id int) (
	*domain.News,
	error,
) {
	ctx := context.Background()

	params := pgx.NamedArgs{"id": id}

	var fields []string

	if m.Title != nil {
		fields = append(fields, "title = @title")
		params["title"] = *m.Title
	}
	if m.Content != nil {
		fields = append(fields, "content = @content")
		params["content"] = *m.Content
	}

	sql := fmt.Sprintf(`
		UPDATE news
		SET %s
		WHERE id = @id returning id,
					title,
					content,
					updated_at`,
		strings.Join(fields, ", "),
	)

	rows, err := ns.db.Query(
		ctx,
		sql,
		params,
	)
	if err != nil {
		return nil, err
	}

	res, err := pgx.CollectExactlyOneRow(
		rows, pgx.RowToAddrOfStructByName[domain.NewsDB],
	)
	if err != nil {
		return nil, err
	}

	ids, err := ns.GetCategoriesToNews(res.ID)

	if err != nil {
		return nil, err
	}

	return domain.NewsDBtoNews(res, ids), nil
}
