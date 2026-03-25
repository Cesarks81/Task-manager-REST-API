package db

import (
	"app/models"
	"database/sql"
	"fmt"
	"time"
)

func GetTasks(db *sql.DB, id string) ([]models.Task, error) {
	var rows *sql.Rows
	var err error

	if id == "" {
		rows, err = db.Query("select * from task")
	} else {
		rows, err = db.Query("select * from task where id = ?", id)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.CompletedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func UpdateTask(db *sql.DB, id, title, description, status string) error {
	var count int
	db.QueryRow("select count(id) from task where id = ?", id).Scan(&count)
	if count != 1 {
		return sql.ErrNoRows
	}

	if !models.AcceptedTaskStatus[status] {
		return fmt.Errorf("invalid status: %s", status)
	}

	if status == "completed" {
		_, err := db.Exec(
			"update task set title = ?, description = ?, status = ?, completedat = ? WHERE id = ?",
			title, description, status, time.Now(), id)
		return err
	}

	_, err := db.Exec(
		"update task set title = ?, description = ?, status = ? where id = ?",
		title, description, status, id)
	return err
}

func DeleteTask(db *sql.DB, id string) error {
	_, err := db.Exec("delete from task where id = ?", id)
	return err
}
