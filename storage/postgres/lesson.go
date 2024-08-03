package postgres

import (
	"context"
	"database/sql"
	"fmt"
	ls "go_schedule_service/genproto/lesson_service"
	"go_schedule_service/pkg"
	"go_schedule_service/storage"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type lessonRepo struct {
	db *pgxpool.Pool
}

func NewLessonRepo(db *pgxpool.Pool) storage.LessonRepoI {
	return &lessonRepo{
		db: db,
	}
}

func (c *lessonRepo) Create(ctx context.Context, req *ls.CreateLesson) (*ls.GetLesson, error) {

	id := uuid.NewString()

	comtag, err := c.db.Exec(ctx, `
		INSERT INTO lessons (
			id,
			theme,
			links,
			type_of_group
		) VALUES ($1,$2,$3,$4
		)`,
		id,
		req.Theme,
		req.Links,
		req.TypeOfGroup,
	)
	if err != nil {
		log.Println("error while creating lesson", comtag)
		return nil, err
	}

	lesson, err := c.GetById(ctx, &ls.LessonPrimaryKey{Id: id})
	if err != nil {
		log.Println("error while getting lesson by id")
		return nil, err
	}
	return lesson, nil
}

func (c *lessonRepo) Update(ctx context.Context, req *ls.UpdateLesson) (*ls.GetLesson, error) {

	_, err := c.db.Exec(ctx, `
		UPDATE lessons SET
		theme = $1,
		links = $2,
		type_of_group = $3,
		updated_at = NOW()
		WHERE id = $4
		`,
		req.Theme,
		req.Links,
		req.TypeOfGroup,
		req.Id)
	if err != nil {
		log.Println("error while updating lesson")
		return nil, err
	}

	lesson, err := c.GetById(ctx, &ls.LessonPrimaryKey{Id: req.Id})
	if err != nil {
		log.Println("error while getting lesson by id")
		return nil, err
	}
	return lesson, nil
}

func (c *lessonRepo) GetAll(ctx context.Context, req *ls.GetListLessonRequest) (*ls.GetListLessonResponse, error) {
	lessons := ls.GetListLessonResponse{}
	var (
		created_at sql.NullString
		updated_at sql.NullString
	)
	filter_by_name := ""
	offset := (req.Offset - 1) * req.Limit
	if req.Search != "" {
		filter_by_name = fmt.Sprintf(`AND theme ILIKE '%%%v%%'`, req.Search)
	}
	query := `SELECT
				id,
				theme,
				links,
				type_of_group,
				created_at,
				updated_at
			FROM lessons
			WHERE TRUE AND deleted_at is null ` + filter_by_name + `
			OFFSET $1 LIMIT $2
`
	rows, err := c.db.Query(ctx, query, offset, req.Limit)

	if err != nil {
		log.Println("error while getting all lessons")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			lesson ls.GetLesson
		)
		if err = rows.Scan(
			&lesson.Id,
			&lesson.Theme,
			&lesson.Links,
			&lesson.TypeOfGroup,
			&created_at,
			&updated_at,
		); err != nil {
			return &lessons, err
		}
		lesson.CreatedAt = pkg.NullStringToString(created_at)
		lesson.UpdatedAt = pkg.NullStringToString(updated_at)

		lessons.Lessons = append(lessons.Lessons, &lesson)
	}

	err = c.db.QueryRow(ctx, `SELECT count(*) from lessons WHERE TRUE AND deleted_at is null `+filter_by_name+``).Scan(&lessons.Count)
	if err != nil {
		return &lessons, err
	}

	return &lessons, nil
}

func (c *lessonRepo) GetById(ctx context.Context, id *ls.LessonPrimaryKey) (*ls.GetLesson, error) {
	var (
		lesson     ls.GetLesson
		created_at sql.NullString
		updated_at sql.NullString
	)

	query := `SELECT
				id,
				theme,
				links,
				type_of_group,
				created_at,
				updated_at
			FROM lessons
			WHERE id = $1 AND deleted_at IS NULL`

	rows := c.db.QueryRow(ctx, query, id.Id)

	if err := rows.Scan(
		&lesson.Id,
		&lesson.Theme,
		&lesson.Links,
		&lesson.TypeOfGroup,
		&created_at,
		&updated_at); err != nil {
		return &lesson, err
	}
	lesson.CreatedAt = pkg.NullStringToString(created_at)
	lesson.UpdatedAt = pkg.NullStringToString(updated_at)

	return &lesson, nil
}

func (c *lessonRepo) Delete(ctx context.Context, id *ls.LessonPrimaryKey) (emptypb.Empty, error) {

	_, err := c.db.Exec(ctx, `
		UPDATE lessons SET
		deleted_at = NOW()
		WHERE id = $1
		`,
		id.Id)

	if err != nil {
		return emptypb.Empty{}, err
	}
	return emptypb.Empty{}, nil
}
