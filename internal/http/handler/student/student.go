package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ananddub/students-api/internal/storage"
	"github.com/ananddub/students-api/internal/types"
	"github.com/ananddub/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		//---------------------------------------------------

		if errors.Is(err, io.EOF) {
			slog.Info("Empty Body")
			response.WriteJson(w, http.StatusBadRequest,
				response.GeneralError(fmt.Errorf("empty body")),
			)
			return
		}
		if err := validator.New().Struct(student); err != nil {
			slog.Info("Invalid Body")
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest,
				response.ValidationError(validateErrs),
			)
			return
		}
		lastID, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			slog.Info("Failed to create student")
			response.WriteJson(w, http.StatusInternalServerError,
				response.GeneralError(err),
			)
			return
		}
		//---------------------------------------------------
		log.Printf("%#v", student)
		slog.Info("user created id", slog.Int64("id", lastID))
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastID})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		intid, err := strconv.ParseInt(id, 10, 64)
		student, err := storage.GetStudentById(int64(intid))
		if err != nil {
			slog.Info("Failed to get student")
			response.WriteJson(w, http.StatusInternalServerError,
				response.GeneralError(err),
			)
			return
		}
		response.WriteJson(w, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		students, err := storage.GetStudents()
		if err != nil {
			slog.Info("Failed to get students")
			response.WriteJson(w, http.StatusInternalServerError,
				response.GeneralError(err),
			)
			return
		}
		response.WriteJson(w, http.StatusOK, students)
	}
}
