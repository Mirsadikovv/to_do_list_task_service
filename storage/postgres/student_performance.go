package postgres

import (
	"context"
	"database/sql"
	"fmt"
	stp "go_schedule_service/genproto/perfomance_service"
	"go_schedule_service/pkg"
	"go_schedule_service/pkg/check"
	"go_schedule_service/storage"
	"log"
	"strconv"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/cast"
)

type perfomanceRepo struct {
	db *pgxpool.Pool
}

func NewPerfomanceRepo(db *pgxpool.Pool) storage.PerfomanceRepoI {
	return &perfomanceRepo{
		db: db,
	}
}

func (c *perfomanceRepo) Create(ctx context.Context, req *stp.CreatePerfomance) (*stp.GetPerfomance, error) {
	var deadline sql.NullString
	var task_score float64
	string_score := cast.ToString(req.TaskScore)
	amount, err := strconv.ParseFloat(string_score, 64)
	if err != nil {
		fmt.Println("error while parsing task_score", err)
		return nil, err
	}
	id := uuid.NewString()
	query := `
			SELECT
			deadline
		FROM schedules
		WHERE id = $1 and deleted_at is null`

	rows := c.db.QueryRow(ctx, query, req.ScheduleId)

	if err := rows.Scan(
		&deadline); err != nil {
		return nil, err
	}
	hoursUntil, err := check.CheckDeadline(pkg.NullStringToString(deadline))
	if err != nil {
		log.Println("error while creating perfomance")
		return nil, err
	}
	if hoursUntil == 0 {
		task_score = 0
	} else if hoursUntil == 1 {
		task_score = amount * 0.5
	} else {
		task_score = amount
	}

	comtag, err := c.db.Exec(ctx, `
		INSERT INTO student_performance (
			id,
			student_id,
			schedule_id,
			attended,
			task_score
		) VALUES ($1,$2,$3,$4,$5
		)`,
		id,
		req.StudentId,
		req.ScheduleId,
		req.Attended,
		task_score)
	if err != nil {
		log.Println("error while creating perfomance", comtag)
		return nil, err
	}

	performance, err := c.GetById(ctx, &stp.PerfomancePrimaryKey{Id: id})
	if err != nil {
		log.Println("error while getting perfomance by id")
		return nil, err
	}
	return performance, nil
}

func (c *perfomanceRepo) Update(ctx context.Context, req *stp.UpdatePerfomance) (*stp.GetPerfomance, error) {

	_, err := c.db.Exec(ctx, `
		UPDATE student_performance SET
		student_id = $1,
		schedule_id = $2,
		attended = $3,
		task_score = $4,
		updated_at = NOW()
		WHERE id = $5
		`,
		req.StudentId,
		req.ScheduleId,
		req.Attended,
		req.TaskScore,
		req.Id)
	if err != nil {
		log.Println("error while updating perfomance")
		return nil, err
	}

	perfomance, err := c.GetById(ctx, &stp.PerfomancePrimaryKey{Id: req.Id})
	if err != nil {
		log.Println("error while getting perfomance by id")
		return nil, err
	}
	return perfomance, nil
}

func (c *perfomanceRepo) GetAll(ctx context.Context, req *stp.GetListPerfomanceRequest) (*stp.GetListPerfomanceResponse, error) {
	perfomances := stp.GetListPerfomanceResponse{}
	var (
		created_at sql.NullString
		updated_at sql.NullString
	)
	filter_by_name := ""
	offest := (req.Offset - 1) * req.Limit
	if req.Search != "" {
		filter_by_name = fmt.Sprintf(`AND fullname ILIKE '%%%v%%'`, req.Search)
	}
	query := `SELECT
				id,
				student_id,
				schedule_id,
				attended,
				task_score,
				created_at,
				updated_at
			FROM student_performance
			WHERE TRUE AND deleted_at is null ` + filter_by_name + `
			OFFSET $1 LIMIT $2
`
	rows, err := c.db.Query(ctx, query, offest, req.Limit)

	if err != nil {
		log.Println("error while getting all perfomances")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			perfomance stp.GetPerfomance
		)
		if err = rows.Scan(
			&perfomance.Id,
			&perfomance.StudentId,
			&perfomance.ScheduleId,
			&perfomance.Attended,
			&perfomance.TaskScore,
			&created_at,
			&updated_at,
		); err != nil {
			return &perfomances, err
		}
		perfomance.CreatedAt = pkg.NullStringToString(created_at)
		perfomance.UpdatedAt = pkg.NullStringToString(updated_at)

		perfomances.Perfomances = append(perfomances.Perfomances, &perfomance)
	}

	err = c.db.QueryRow(ctx, `SELECT count(*) from student_performance WHERE TRUE AND deleted_at is null `+filter_by_name+``).Scan(&perfomances.Count)
	if err != nil {
		return &perfomances, err
	}

	return &perfomances, nil
}

func (c *perfomanceRepo) GetById(ctx context.Context, id *stp.PerfomancePrimaryKey) (*stp.GetPerfomance, error) {
	var (
		perfomance stp.GetPerfomance
		created_at sql.NullString
		updated_at sql.NullString
	)

	query := `SELECT
				id,
				student_id,
				schedule_id,
				attended,
				task_score,
				created_at,
				updated_at
				FROM student_performance
			WHERE id = $1 AND deleted_at IS NULL`

	rows := c.db.QueryRow(ctx, query, id.Id)

	if err := rows.Scan(
		&perfomance.Id,
		&perfomance.StudentId,
		&perfomance.ScheduleId,
		&perfomance.Attended,
		&perfomance.TaskScore,
		&created_at,
		&updated_at); err != nil {
		return &perfomance, err
	}
	perfomance.CreatedAt = pkg.NullStringToString(created_at)
	perfomance.UpdatedAt = pkg.NullStringToString(updated_at)

	return &perfomance, nil
}

func (c *perfomanceRepo) Delete(ctx context.Context, id *stp.PerfomancePrimaryKey) (emptypb.Empty, error) {

	_, err := c.db.Exec(ctx, `
		UPDATE student_performance SET
		deleted_at = NOW()
		WHERE id = $1
		`,
		id.Id)

	if err != nil {
		return emptypb.Empty{}, err
	}
	return emptypb.Empty{}, nil
}
