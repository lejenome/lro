package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lejenome/lro/services/process-executor/lib/process"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type jobT struct {
	ID          uuid.UUID         `gorm:"primary_key;type:uuid;default:gen_random_uuid();<-:false;->"`
	ProcessName string            `gorm:"<-:create;not null;index"`
	Input       string            `gorm:"<-:create;type:text;not null"`
	Output      string            `gorm:"<-:update;type:text"`
	Errors      string            `gorm:"<-:update;type:text"`
	Status      process.JobStatus `gorm:"default:PENDING;not null;index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (jobT) TableName() string {
	return "lro_job"
}

func fromJob(j *process.Job) (*jobT, error) {
	jt := jobT{
		ID:          j.ID,
		ProcessName: j.ProcessName,
		Status:      j.Status,
	}
	if v, err := json.Marshal(j.Input); err != nil {
		return nil, fmt.Errorf("Can not convert Job struct to Job Row: %w", err)
	} else {
		jt.Input = string(v)
	}
	if j.Output == nil {
		jt.Output = ""
	} else if v, err := json.Marshal(j.Output); err != nil {
		return nil, fmt.Errorf("Can not convert Job struct to Job Row: %w", err)
	} else {
		jt.Output = string(v)
	}
	if /*j.Errors == nil ||*/ len(j.Errors) == 0 {
		jt.Errors = ""
	} else if v, err := json.Marshal(j.Errors); err != nil {
		return nil, fmt.Errorf("Can not convert Job struct to Job Row: %w", err)
	} else {
		jt.Errors = string(v)
	}
	return &jt, nil
}
func (jt *jobT) toJob() (*process.Job, error) {
	j := process.Job{
		ID:          jt.ID,
		ProcessName: jt.ProcessName,
		Status:      jt.Status,
		Output:      nil,
		Errors:      nil,
	}
	if err := json.Unmarshal([]byte(jt.Input), &j.Input); err != nil {
		return nil, err
	}
	if jt.Output != "" {
		if err := json.Unmarshal([]byte(jt.Output), &j.Output); err != nil {
			return nil, err
		}
	}
	if jt.Errors != "" {
		if err := json.Unmarshal([]byte(jt.Errors), &j.Errors); err != nil {
			return nil, err
		}
	}
	return &j, nil
}

type dbJobStore struct {
	db *gorm.DB
}

func DBJobStore(dsn string) process.JobStore {
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		// SkipDefaultTransaction: true,
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err.Error())
	}
	DB.AutoMigrate(&jobT{})
	return &dbJobStore{
		db: DB,
	}
}
func (q *dbJobStore) Save(job *process.Job) error {
	if job.ID != uuid.Nil {
		return fmt.Errorf("Can not save a job with already defined id '%s'", job.ID.String())
	}
	jt, err := fromJob(job)
	if err != nil {
		return err
	}
	res := q.db.Omit("ID").Create(jt)
	if res.Error != nil {
		return res.Error
	}
	job.ID = jt.ID
	return nil
}
func (q *dbJobStore) Update(job *process.Job) error {
	if job.ID == uuid.Nil {
		return errors.New("Can not update a job with null id")
	}
	jt, err := fromJob(job)
	if err != nil {
		return err
	}
	res := q.db.Model(jt).
		Where("id = ?", jt.ID).
		Updates(jt)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (q *dbJobStore) Retrive(id uuid.UUID) (*process.Job, error) {
	var jt jobT
	res := q.db.Model(&jt).Take(&jt, id)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("Job '%s' not found in store", id)
	} else if res.Error != nil {
		return nil, res.Error
	}
	return jt.toJob()
}
