package handlers

import (
	"app/db"
	"app/models"
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterTaskRoutes(r *gin.Engine, dbConn *sql.DB) {
	r.GET("/api/tasks", func(c *gin.Context) {
		tasks, err := db.GetTasks(dbConn, "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(tasks) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No tasks yet"})
			return
		}
		c.JSON(http.StatusOK, tasks)
	})

	r.GET("/api/tasks/search", func(c *gin.Context) {
		title := strings.ToLower(c.Query("title"))
		status := strings.ToLower(c.Query("status"))

		var rows *sql.Rows
		var err error

		switch {
		case title != "":
			rows, err = dbConn.Query("select * from task where lower(title) = ?", title)
		case status != "":
			rows, err = dbConn.Query("select * from task where lower(status) = ?", status)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter required"})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Query error"})
			return
		}
		defer rows.Close()

		var tasks []models.Task
		for rows.Next() {
			var t models.Task
			if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.CompletedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read results"})
				return
			}
			tasks = append(tasks, t)
		}
		c.JSON(http.StatusOK, tasks)
	})

	r.GET("/api/tasks/:id", func(c *gin.Context) {
		tasks, err := db.GetTasks(dbConn, c.Param("id"))
		if err != nil || len(tasks) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusOK, tasks)
	})

	r.POST("/api/tasks", func(c *gin.Context) {
		var newTask models.Task
		if err := c.ShouldBindJSON(&newTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Json"})
			return
		}

		if !models.AcceptedTaskStatus[newTask.Status] || newTask.Title == "" || newTask.Description == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status, title, or description"})
			return
		}

		_, err := dbConn.Exec("insert into task values (default, ?, ?, ?, ?, default)",
			newTask.Title, newTask.Description, newTask.Status, time.Now())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Insert failed"})
			return
		}
		c.Status(http.StatusCreated)
	})

	r.PUT("/api/tasks/:id", func(c *gin.Context) {
		var updatedTask models.Task
		if err := c.ShouldBindJSON(&updatedTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
		err := db.UpdateTask(dbConn, c.Param("id"), updatedTask.Title, updatedTask.Description, updatedTask.Status)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "ID not found"})
			return
		}
		c.Status(http.StatusOK)
	})

	r.DELETE("/api/tasks/:id", func(c *gin.Context) {
		db.DeleteTask(dbConn, c.Param("id"))
		c.Status(http.StatusNoContent)
	})
}
