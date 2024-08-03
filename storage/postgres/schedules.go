package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"go_schedule_service/genproto/group_service"
	sched "go_schedule_service/genproto/schedule_service"
	"go_schedule_service/grpc/client"
	"go_schedule_service/pkg"
	"go_schedule_service/storage"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type scheduleRepo struct {
	db       *pgxpool.Pool
	services client.GrpcClientI
}

func NewScheduleRepo(db *pgxpool.Pool, srvs client.GrpcClientI) storage.ScheduleRepoI {
	return &scheduleRepo{
		db:       db,
		services: srvs,
	}
}

func (c *scheduleRepo) Create(ctx context.Context, req *sched.CreateSchedule) (*sched.GetSchedule, error) {

	id := uuid.NewString()

	comtag, err := c.db.Exec(ctx, `
		INSERT INTO schedules (
			id,
			group_id,
			lesson_id,
			classroom,
			type_of_group,
			task,
			deadline,
			score,
			start_time,
			end_time
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10
		)`,
		id,
		req.GroupId,
		req.LessonId,
		req.Classroom,
		req.TypeOfGroup,
		req.Task,
		req.Deadline,
		req.Score,
		req.StartTime,
		req.EndTime,
	)
	if err != nil {
		log.Println("error while creating schedules", comtag)
		return nil, err
	}

	schedules, err := c.GetById(ctx, &sched.SchedulePrimaryKey{Id: id})
	if err != nil {
		log.Println("error while getting schedules by id")
		return nil, err
	}
	return schedules, nil
}

func (c *scheduleRepo) Update(ctx context.Context, req *sched.UpdateSchedule) (*sched.GetSchedule, error) {

	_, err := c.db.Exec(ctx, `
		UPDATE schedules SET
		group_id = $1,
		lesson_id = $2,
		classroom = $3,
		type_of_group = $4,
		task = $5,
		deadline = $6,
		score = $7,
		start_time = $8,
		end_time = $9,
		updated_at = NOW()
		WHERE id = $10
		`,
		req.GroupId,
		req.LessonId,
		req.Classroom,
		req.TypeOfGroup,
		req.Task,
		req.Deadline,
		req.Score,
		req.StartTime,
		req.EndTime,
		req.Id)
	if err != nil {
		log.Println("error while updating schedules")
		return nil, err
	}

	schedules, err := c.GetById(ctx, &sched.SchedulePrimaryKey{Id: req.Id})
	if err != nil {
		log.Println("error while getting schedules by id")
		return nil, err
	}
	return schedules, nil
}

func (c *scheduleRepo) GetAll(ctx context.Context, req *sched.GetListScheduleRequest) (*sched.GetListScheduleResponse, error) {
	schedules := sched.GetListScheduleResponse{}
	var (
		deadline   sql.NullString
		start_time sql.NullString
		end_time   sql.NullString
		created_at sql.NullString
		updated_at sql.NullString
	)
	filter_by_name := ""
	offset := (req.Offset - 1) * req.Limit
	if req.Search != "" {
		filter_by_name = fmt.Sprintf(`AND classroom ILIKE '%%%v%%'`, req.Search)
	}
	query := `SELECT
				id,
				group_id,
				lesson_id,
				classroom,
				type_of_group,
				task,
				deadline,
				score,
				start_time,
				end_time,
				created_at,
				updated_at
			FROM schedules
			WHERE TRUE AND deleted_at is null ` + filter_by_name + `
			OFFSET $1 LIMIT $2
`
	rows, err := c.db.Query(ctx, query, offset, req.Limit)

	if err != nil {
		log.Println("error while getting all schedules")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			schedule sched.GetSchedule
		)
		if err = rows.Scan(
			&schedule.Id,
			&schedule.GroupId,
			&schedule.LessonId,
			&schedule.Classroom,
			&schedule.TypeOfGroup,
			&schedule.Task,
			&deadline,
			&schedule.Score,
			&start_time,
			&end_time,
			&created_at,
			&updated_at,
		); err != nil {
			return &schedules, err
		}
		schedule.Deadline = pkg.NullStringToString(deadline)
		schedule.StartTime = pkg.NullStringToString(start_time)
		schedule.EndTime = pkg.NullStringToString(end_time)
		schedule.CreatedAt = pkg.NullStringToString(created_at)
		schedule.UpdatedAt = pkg.NullStringToString(updated_at)

		schedules.Schedules = append(schedules.Schedules, &schedule)
	}

	err = c.db.QueryRow(ctx, `SELECT count(*) from schedules WHERE TRUE AND deleted_at is null `+filter_by_name+``).Scan(&schedules.Count)
	if err != nil {
		return &schedules, err
	}

	return &schedules, nil
}

func (c *scheduleRepo) GetById(ctx context.Context, id *sched.SchedulePrimaryKey) (*sched.GetSchedule, error) {
	var (
		deadline   sql.NullString
		schedule   sched.GetSchedule
		start_time sql.NullString
		end_time   sql.NullString
		created_at sql.NullString
		updated_at sql.NullString
	)

	query := `SELECT
				id,
				group_id,
				lesson_id,
				classroom,
				type_of_group,
				task,
				deadline,
				score,
				start_time,
				end_time,
				created_at,
				updated_at
			FROM schedules
			WHERE id = $1 AND deleted_at IS NULL`

	rows := c.db.QueryRow(ctx, query, id.Id)

	if err := rows.Scan(
		&schedule.Id,
		&schedule.GroupId,
		&schedule.LessonId,
		&schedule.Classroom,
		&schedule.TypeOfGroup,
		&schedule.Task,
		&deadline,
		&schedule.Score,
		&start_time,
		&end_time,
		&created_at,
		&updated_at); err != nil {
		return &schedule, err
	}
	schedule.Deadline = pkg.NullStringToString(deadline)
	schedule.StartTime = pkg.NullStringToString(start_time)
	schedule.EndTime = pkg.NullStringToString(end_time)
	schedule.CreatedAt = pkg.NullStringToString(created_at)
	schedule.UpdatedAt = pkg.NullStringToString(updated_at)

	return &schedule, nil
}

func (c *scheduleRepo) Delete(ctx context.Context, id *sched.SchedulePrimaryKey) (emptypb.Empty, error) {

	_, err := c.db.Exec(ctx, `
		UPDATE schedules SET
		deleted_at = NOW()
		WHERE id = $1
		`,
		id.Id)

	if err != nil {
		return emptypb.Empty{}, err
	}
	return emptypb.Empty{}, nil
}

func (c *scheduleRepo) GetStudentSchedule(ctx context.Context, id *sched.SchedulePrimaryKey) (*sched.GetStudentSchedules, error) {
	var (
		deadline   sql.NullString
		schedule   sched.GetStudentSchedules
		start_time sql.NullString
		end_time   sql.NullString
	)

	query := `SELECT
				sc.classroom,
				sc.type_of_group,
				sc.task,
				sc.deadline,
				sc.score,
				sc.start_time,
				sc.end_time,
				sp.attended,
				sp.task_score
			FROM student_performance sp 
			JOIN schedules sc ON sp.schedule_id = sc.id 
			WHERE TRUE AND sc.deleted_at is null 
`
	rows := c.db.QueryRow(ctx, query)

	if err := rows.Scan(
		&schedule.Classroom,
		&schedule.TypeOfGroup,
		&schedule.Task,
		&deadline,
		&schedule.Score,
		&start_time,
		&end_time,
		&schedule.Attended,
		&schedule.TaskScore); err != nil {
		return &schedule, err
	}
	schedule.Deadline = pkg.NullStringToString(deadline)
	schedule.StartTime = pkg.NullStringToString(start_time)
	schedule.EndTime = pkg.NullStringToString(end_time)

	return &schedule, nil
}

func (c *scheduleRepo) GetTeacherSchedule(ctx context.Context, id *sched.SchedulePrimaryKey) (*sched.GetStudentSchedules, error) {
	var (
		deadline   sql.NullString
		schedule   sched.GetStudentSchedules
		start_time sql.NullString
		end_time   sql.NullString
	)

	query := `SELECT
				sc.classroom,
				sc.type_of_group,
				sc.task,
				sc.deadline,
				sc.score,
				sc.start_time,
				sc.end_time,
				sp.attended,
				sp.task_score
			FROM student_performance sp 
			JOIN schedules sc ON sp.schedule_id = sc.id 
			WHERE TRUE AND sc.deleted_at is null 
`
	rows := c.db.QueryRow(ctx, query)

	if err := rows.Scan(
		&schedule.Classroom,
		&schedule.TypeOfGroup,
		&schedule.Task,
		&deadline,
		&schedule.Score,
		&start_time,
		&end_time,
		&schedule.Attended,
		&schedule.TaskScore); err != nil {
		return &schedule, err
	}
	schedule.Deadline = pkg.NullStringToString(deadline)
	schedule.StartTime = pkg.NullStringToString(start_time)
	schedule.EndTime = pkg.NullStringToString(end_time)

	return &schedule, nil
}

func (c *scheduleRepo) GetForWeek(ctx context.Context, req *sched.GetWeekScheduleRequest) (*sched.GetWeekScheduleResponse, error) {
	schedules := sched.GetWeekScheduleResponse{}
	var (
		start_time sql.NullString
		end_time   sql.NullString
	)

	filter_by_date := ""

	if req.FromDate != "" {
		filter_by_date = fmt.Sprintln(` AND jo.date_of_lesson BETWEEN '`, req.FromDate, `' AND '`, req.ToDate, `'`)
	}
	query := `SELECT
				sc.group_id,
				sc.type_of_group,
				sc.classroom,
				sc.start_time,
				sc.end_time
			FROM schedules sc 
			JOIN journals jo ON jo.schedule_id = sc.id
			WHERE TRUE AND sc.deleted_at is null ` + filter_by_date + ``

	rows, err := c.db.Query(ctx, query)

	if err != nil {
		log.Println("error while getting all schedules")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			schedule sched.GetScheduleWeek
		)

		if err = rows.Scan(
			&schedule.GroupId,
			&schedule.TypeOfGroup,
			&schedule.Classroom,
			&start_time,
			&end_time,
		); err != nil {
			return &schedules, err
		}
		id := &group_service.GroupPrimaryKey{Id: schedule.GroupId}
		resp, err := c.services.Group().GetTBS(ctx, id)
		if err != nil {
			log.Println("error while gettinf groupinfo", err)
			return nil, err
		}
		schedule.GroupName = resp.GroupName
		schedule.BranchName = resp.BranchName
		schedule.TeacherName = resp.TeacherName
		schedule.SupportTeacherName = resp.SupportTeacherName
		schedule.StudentsCount = resp.StudentCount
		schedule.StartTime = pkg.NullStringToString(start_time)
		schedule.EndTime = pkg.NullStringToString(end_time)

		if err := rows.Scan(
			&start_time); err != nil {
			return nil, err
		}

		schedules.Schedules = append(schedules.Schedules, &schedule)
	}

	err = c.db.QueryRow(ctx, `  SELECT count(*) FROM schedules sc 
								JOIN journals jo ON jo.schedule_id = sc.id
								WHERE TRUE AND sc.deleted_at is null  `+filter_by_date+``).Scan(&schedules.Count)
	if err != nil {
		return &schedules, err
	}

	return &schedules, nil
}
