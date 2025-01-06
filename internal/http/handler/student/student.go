package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Wasid786/student-info/internal/storage"
	"github.com/Wasid786/student-info/internal/types"
	"github.com/Wasid786/student-info/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return

		}
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return

		}

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)
		slog.Info("user created successfully ", slog.String("userId", fmt.Sprint(lastId)))

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		slog.Info("Creating a Student ")
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
		w.Write([]byte("Welcome to Students API"))
	}

}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Getting a student", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("Invalid student ID", slog.String("id", id), slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid student ID: %s", id)))
			return
		}

		student, err := storage.GetStudentById(intId)
		if err != nil {
			// Assuming storage.ErrNotFound is used for "not found"
			slog.Warn("Student not found", slog.String("id", id))
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("student with ID %d not found", intId)))
			return

			// Handle other errors
			slog.Error("Error retrieving student", slog.String("id", id), slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, student)
	}
}
