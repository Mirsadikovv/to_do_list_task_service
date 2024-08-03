package postgres

import (
	"context"
	"database/sql"
	"fmt"
	jo "go_schedule_service/genproto/journal_service"
	"go_schedule_service/pkg"
	"go_schedule_service/storage"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type journalRepo struct {
	db *pgxpool.Pool
}

func NewJournalRepo(db *pgxpool.Pool) storage.JournalRepoI {
	return &journalRepo{
		db: db,
	}
}

func (c *journalRepo) Create(ctx context.Context, req *jo.CreateJournal) (*jo.GetJournal, error) {

	id := uuid.NewString()

	comtag, err := c.db.Exec(ctx, `
		INSERT INTO journals (
			id,
			schedule_id,
			date_of_lesson
		) VALUES ($1,$2,$3
		)`,
		id,
		req.ScheduleId,
		req.DateOfLesson,
	)
	if err != nil {
		log.Println("error while creating journal", comtag)
		return nil, err
	}

	journal, err := c.GetById(ctx, &jo.JournalPrimaryKey{Id: id})
	if err != nil {
		log.Println("error while getting journal by id")
		return nil, err
	}
	return journal, nil
}

func (c *journalRepo) Update(ctx context.Context, req *jo.UpdateJournal) (*jo.GetJournal, error) {

	_, err := c.db.Exec(ctx, `
		UPDATE journals SET
		schedule_id = $1,
		date_of_lesson = $2,
		updated_at = NOW()
		WHERE id = $3
		`,
		req.ScheduleId,
		req.DateOfLesson,
		req.Id)
	if err != nil {
		log.Println("error while updating journal")
		return nil, err
	}

	journal, err := c.GetById(ctx, &jo.JournalPrimaryKey{Id: req.Id})
	if err != nil {
		log.Println("error while getting journal by id")
		return nil, err
	}
	return journal, nil
}

func (c *journalRepo) GetAll(ctx context.Context, req *jo.GetListJournalRequest) (*jo.GetListJournalResponse, error) {
	journals := jo.GetListJournalResponse{}
	var (
		created_at sql.NullString
		updated_at sql.NullString
	)
	filter_by_date := ""
	offest := (req.Offset - 1) * req.Limit
	if req.Search != "" {
		filter_by_date = fmt.Sprintf(`AND date_of_lesson = '%%%v%%'`, req.Search)
	}
	query := `SELECT
				id,
				schedule_id,
				date_of_lesson,
				created_at,
				updated_at
			FROM journals
			WHERE TRUE AND deleted_at is null ` + filter_by_date + `
			OFFSET $1 LIMIT $2
`
	rows, err := c.db.Query(ctx, query, offest, req.Limit)

	if err != nil {
		log.Println("error while getting all journals")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			journal jo.GetJournal
		)
		if err = rows.Scan(
			&journal.Id,
			&journal.ScheduleId,
			&journal.DateOfLesson,
			&created_at,
			&updated_at,
		); err != nil {
			return &journals, err
		}
		journal.CreatedAt = pkg.NullStringToString(created_at)
		journal.UpdatedAt = pkg.NullStringToString(updated_at)

		journals.Journals = append(journals.Journals, &journal)
	}

	err = c.db.QueryRow(ctx, `SELECT count(*) from journals WHERE TRUE AND deleted_at is null `+filter_by_date+``).Scan(&journals.Count)
	if err != nil {
		return &journals, err
	}

	return &journals, nil
}

func (c *journalRepo) GetById(ctx context.Context, id *jo.JournalPrimaryKey) (*jo.GetJournal, error) {
	var (
		journal     jo.GetJournal
		created_at  sql.NullString
		updated_at  sql.NullString
		date_lesson sql.NullString
	)

	query := `SELECT
				id,
				schedule_id,
				date_of_lesson,
				created_at,
				updated_at
			FROM journals
			WHERE id = $1 AND deleted_at IS NULL`

	rows := c.db.QueryRow(ctx, query, id.Id)

	if err := rows.Scan(
		&journal.Id,
		&journal.ScheduleId,
		&date_lesson,
		&created_at,
		&updated_at); err != nil {
		return &journal, err
	}
	journal.DateOfLesson = pkg.NullStringToString(date_lesson)
	journal.CreatedAt = pkg.NullStringToString(created_at)
	journal.UpdatedAt = pkg.NullStringToString(updated_at)

	return &journal, nil
}

func (c *journalRepo) Delete(ctx context.Context, id *jo.JournalPrimaryKey) (emptypb.Empty, error) {

	_, err := c.db.Exec(ctx, `
		UPDATE journals SET
		deleted_at = NOW()
		WHERE id = $1
		`,
		id.Id)

	if err != nil {
		return emptypb.Empty{}, err
	}
	return emptypb.Empty{}, nil
}
