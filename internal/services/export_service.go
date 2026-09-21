package services

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"employee-management/internal/models"
	"employee-management/internal/repositories"
	"employee-management/internal/utils"
)

const exportTimeout = 10 * time.Second

type ExportService struct {
	employees repositories.EmployeeRepository
	outputDir string
}

func NewExportService(employees repositories.EmployeeRepository, outputDir string) *ExportService {
	return &ExportService{employees: employees, outputDir: outputDir}
}

type ExportResult struct {
	JSONPath string `json:"jsonPath"`
	CSVPath  string `json:"csvPath"`
}

func (s *ExportService) ExportAll(ctx context.Context) (*ExportResult, error) {
	ctx, cancel := context.WithTimeout(ctx, exportTimeout)
	defer cancel()

	employees, _, err := s.employees.List(ctx, models.EmployeeListFilter{Limit: 1_000_000, Offset: 0})
	if err != nil {
		return nil, utils.AsAppError(err)
	}

	if err := os.MkdirAll(s.outputDir, 0o755); err != nil {
		return nil, models.ErrInternal(fmt.Errorf("create output dir: %w", err))
	}

	jsonPath := filepath.Join(s.outputDir, "employees.json")
	csvPath := filepath.Join(s.outputDir, "employees.csv")

	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		errs []error
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := writeEmployeesJSON(s.outputDir, jsonPath, employees); err != nil {
			mu.Lock()
			errs = append(errs, fmt.Errorf("json export: %w", err))
			mu.Unlock()
		}
	}()
	go func() {
		defer wg.Done()
		if err := writeEmployeesCSV(s.outputDir, csvPath, employees); err != nil {
			mu.Lock()
			errs = append(errs, fmt.Errorf("csv export: %w", err))
			mu.Unlock()
		}
	}()
	wg.Wait()

	if len(errs) > 0 {
		return nil, models.ErrInternal(errors.Join(errs...))
	}

	return &ExportResult{JSONPath: jsonPath, CSVPath: csvPath}, nil
}

func writeEmployeesJSON(dir, path string, employees []*models.Employee) error {
	return writeFileAtomic(dir, path, func(f *os.File) error {
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		return enc.Encode(employees)
	})
}

func writeEmployeesCSV(dir, path string, employees []*models.Employee) error {
	return writeFileAtomic(dir, path, func(f *os.File) error {
		w := csv.NewWriter(f)

		header := []string{"id", "name", "age", "position", "department_id", "salary", "created_at", "updated_at"}
		if err := w.Write(header); err != nil {
			return err
		}

		for _, e := range employees {
			record := []string{
				strconv.FormatInt(e.ID, 10),
				e.Name,
				strconv.Itoa(e.Age),
				e.Position,
				strconv.FormatInt(e.DepartmentID, 10),
				e.Salary.String(),
				e.CreatedAt.Format(time.RFC3339),
				e.UpdatedAt.Format(time.RFC3339),
			}
			if err := w.Write(record); err != nil {
				return err
			}
		}
		w.Flush()
		return w.Error()
	})
}

// writeFileAtomic writes into a temp file in dir and renames it onto path
// once fully written. Concurrent exports (e.g. overlapping POST
// /employees/export calls) would otherwise both os.Create the same path and
// interleave/truncate each other's writes; rename(2) within the same
// filesystem is atomic, so readers and other writers only ever see a
// complete old or new file, never a half-written one.
func writeFileAtomic(dir, path string, write func(*os.File) error) (err error) {
	tmp, err := os.CreateTemp(dir, ".export-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmp.Name()
	defer func() {
		if err != nil {
			os.Remove(tmpPath)
		}
	}()

	if err = write(tmp); err != nil {
		tmp.Close()
		return fmt.Errorf("write temp file: %w", err)
	}
	if err = tmp.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}
	if err = os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("rename temp file: %w", err)
	}
	return nil
}
