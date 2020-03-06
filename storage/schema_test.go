package storage

import (
	"fmt"
	"testing"
)

func TestCreateTables(t *testing.T) {
	s := OpenDriver()
	defer CloseDriver(s, t)

	s.DropTablesIfExist()

	t.Run("assertion that there is no table initially", func(t *testing.T) {
		if s.driver.HasTable(&Command{}) ||
			s.driver.HasTable(&TaskConfig{}) ||
			s.driver.HasTable(&TaskMetadata{}) ||
			s.driver.HasTable(&Task{}) ||
			s.driver.HasTable(&Job{}) ||
			s.driver.HasTable(&ResourceNode{}) ||
			s.driver.HasTable(&Queue{}) {

			t.Errorf("expected to have no table but has at least one")
		}
	})

	t.Run("assertion that all tables were created", func(t *testing.T) {
		s.CreateTables()
		if !s.driver.HasTable(&Command{}) ||
			!s.driver.HasTable(&TaskConfig{}) ||
			!s.driver.HasTable(&TaskMetadata{}) ||
			!s.driver.HasTable(&Task{}) ||
			!s.driver.HasTable(&Job{}) ||
			!s.driver.HasTable(&ResourceNode{}) ||
			!s.driver.HasTable(&Queue{}) {

			t.Errorf("expected to have no table but has at least one")
		}
	})

	t.Run("assertion that two tables don't will be created with same name", func(t *testing.T) {
		s.DropTablesIfExist()

		t.Run("assertion that the first table are created correctly", func(t *testing.T){
			cmd := Command{}
			err, got := s.CreateTable(cmd)
			var want = fmt.Sprintf("Table %+v correctly created", cmd)
			assertMsg(t, got, want, err)
		})

		t.Run("assertion that the second table don't will be created", func(t *testing.T){
			cmd := Command{}
			err, got := s.CreateTable(cmd)
			var want = fmt.Sprintf("Table %+v already exists", cmd)
			assertMsg(t, got, want, err)
		})
	})

	s.DropTablesIfExist()
}

func TestDropTables(t *testing.T) {
	s := OpenDriver()
	defer CloseDriver(s, t)

	s.CreateTables()

	t.Run("assert drop with tables", func(t *testing.T) {
		s.DropTablesIfExist()
		if s.driver.HasTable(&Command{}) ||
			s.driver.HasTable(&TaskConfig{}) ||
			s.driver.HasTable(&TaskMetadata{}) ||
			s.driver.HasTable(&Task{}) ||
			s.driver.HasTable(&Job{}) ||
			s.driver.HasTable(&ResourceNode{}) ||
			s.driver.HasTable(&Queue{}) {

			t.Errorf("expected to have no table but has at least one")
		}
	})

	t.Run("assert drop without tables", func(t *testing.T) {
		got := s.DropTablesIfExist()
		if  got.Error != nil {
			t.Errorf("expected that nothing changes but an error occurred")
		}
	})

	s.DropTablesIfExist()
}

func TestAutoMigrate(t *testing.T) {
	s := OpenDriver()
	defer CloseDriver(s, t)
	s.DropTablesIfExist()

	cmd := Command{}
	var want = fmt.Sprintf("Table %+v correctly created", cmd)
	err, got := s.CreateTable(cmd)

	if err != nil {
		assertMsg(t, got, want, err)
	}
}

func assertMsg(t *testing.T, got, want string, err error) {
	t.Helper()
	if got != want {
		t.Errorf("want %q but got %q", want, got)
	}
}

func OpenDriver() *Storage {
	return New("127.0.0.1", "5432", "arrebol-admin",
		"arrebol-db", "postgres")
}

func CloseDriver(s *Storage, t *testing.T) {
	err := s.driver.Close()

	if err != nil {
		t.Fail()
	}
}