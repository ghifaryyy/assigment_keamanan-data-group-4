package handler

import (
	"aaa_bekend/internal/db"
	"aaa_bekend/internal/lmodel"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
)

type SubmitProjectReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GradeProjectReq struct {
	Score    int    `json:"score"`
	Feedback string `json:"feedback"`
}

// POST /api/v1/projects

func SubmitProject(w http.ResponseWriter, r *http.Request) {
	userCtx, _ := lmodel.UserFromContext(r.Context())
	// TODO: Decode request body
	var req SubmitProjectReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Simpan proyek ke database
	projectID, err := db.CreateProject(r.Context(), req.Title, req.Description, userCtx.ID)
	if err != nil {
		http.Error(w, "Failed to create project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"project_id": projectID})
}

func ListProjects(w http.ResponseWriter, r *http.Request) {
	// TODO: Ambil semua proyek dari database
	projects, err := db.GetAllProjects(r.Context())
	// TODO: Handle error
	if err != nil {
		http.Error(w, "Failed to fetch projects", http.StatusInternalServerError)
		return
	}
	if projects == nil {
		projects = []lmodel.Project{}
	}
	json.NewEncoder(w).Encode(&projects)
}
func GetAllProjectGrades(w http.ResponseWriter, r *http.Request) {
	// TODO: Ambil semua nilai proyek dari database
	projects, err := db.GetAllProjectGrades(r.Context())
	// TODO: Handle error
	if err != nil {
		http.Error(w, "Failed to fetch project grades", http.StatusInternalServerError)
		return
	}
	if projects == nil {
		projects = []lmodel.ProjectGrade{}
	}
	json.NewEncoder(w).Encode(&projects)
}
func GradeProject(w http.ResponseWriter, r *http.Request) {
	userCtx, _ := lmodel.UserFromContext(r.Context())

	// TODO: Parse project ID dari URL path dan decode request body
	idStr := r.PathValue("id")
	projectID, err := strconv.Atoi(idStr)
	if err != nil || projectID <= 0 {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	var req GradeProjectReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Score < 0 || req.Score > 100 {
		http.Error(w, "Score must be between 0 and 100", http.StatusBadRequest)
		return
	}

	// TODO: Simpan nilai ke database
	if err := db.CreateGrade(r.Context(), projectID, userCtx.ID, req.Score, req.Feedback); err != nil {
		http.Error(w, "Failed to submit grade", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Grade Submitted"})

}

func GetProjectGrades(w http.ResponseWriter, r *http.Request) {
	userCtx, _ := lmodel.UserFromContext(r.Context())

	// TODO: Parse project ID dari URL path
	idStr := r.PathValue("id")
	projectID, err := strconv.Atoi(idStr)
	if err != nil || projectID <= 0 {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	// TODO: Jika role student, periksa kepemilikan proyek
	if userCtx.Role == "student" {
		project, err := db.GetProjectByID(r.Context(), projectID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				http.Error(w, "Project not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Failed to fetch project", http.StatusInternalServerError)
			return
		}

		if project.StudentID != userCtx.ID {
			http.Error(w, "Forbidden: you can only view grades for your own projects", http.StatusForbidden)
			return
		}

	}

	grades, err := db.GetGradesForProject(r.Context(), projectID)
	if err != nil {
		http.Error(w, "Failed to fetch grades", http.StatusInternalServerError)
		return
	}

	if grades == nil {
		grades = []lmodel.Grade{}
	}
	json.NewEncoder(w).Encode(grades)

}
