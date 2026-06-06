package dao

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/GoSimplicity/AI-CloudOps/internal/model"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestWorkorderInstanceDAO(t *testing.T) WorkorderInstanceDAO {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite database: %v", err)
	}

	if err := db.AutoMigrate(&model.WorkorderInstance{}); err != nil {
		t.Fatalf("migrate workorder instance table: %v", err)
	}

	return NewWorkorderInstanceDAO(db, zap.NewNop())
}

func TestGenerateSerialNumberStartsAtOneWhenNoInstanceExistsToday(t *testing.T) {
	instanceDAO := newTestWorkorderInstanceDAO(t)
	prefix := "WO" + time.Now().Format("20060102")

	serialNumber, err := instanceDAO.GenerateSerialNumber(context.Background())
	if err != nil {
		t.Fatalf("GenerateSerialNumber returned error: %v", err)
	}

	if !strings.HasPrefix(serialNumber, prefix) {
		t.Fatalf("expected serial number prefix %q, got %q", prefix, serialNumber)
	}
	if serialNumber != prefix+"0001" {
		t.Fatalf("expected first serial number %q, got %q", prefix+"0001", serialNumber)
	}
}
