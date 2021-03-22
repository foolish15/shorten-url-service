package gorm_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/foolish15/shorten-url-service/internal/repositories/token"
	"github.com/foolish15/shorten-url-service/internal/schemas"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func connectGormDB() *gorm.DB {
	connectString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	logrus.Debugf(connectString)
	for {
		db, err := gorm.Open(
			mysql.Open(connectString),
			&gorm.Config{
				SkipDefaultTransaction: true,
			},
		)

		if err != nil {
			time.Sleep(5 * time.Second)
			logrus.Errorf("ConnectDB error: %+v", err)
			continue
		}
		sqlDB, err := db.DB()
		if err != nil {
			time.Sleep(5 * time.Second)
			logrus.Errorf("ConnectDB error: %+v", err)
			continue
		}

		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(500)
		sqlDB.SetConnMaxLifetime(time.Duration(30) * time.Minute)
		return db
	}
}

func TestFind(t *testing.T) {
	err := godotenv.Load("../../../../.env")
	assert.Nil(t, err)
	db := connectGormDB()

	repo := token.New(db.Debug())
	_, err = repo.Find(token.WhereSub{Sub: "test"}, token.WhereExpLessThan{Exp: 500})
	assert.Nil(t, err)
	_, err = repo.Find(token.WhereExpLessThan{Exp: 500})
	assert.Nil(t, err)
	_, err = repo.Find(token.WhereSub{Sub: "test"})
	assert.Nil(t, err)
}

func TestCount(t *testing.T) {
	err := godotenv.Load("../../../../.env")
	assert.Nil(t, err)
	db := connectGormDB()

	repo := token.New(db.Debug())
	_, err = repo.Count(token.WhereSub{Sub: "test"}, token.WhereExpLessThan{Exp: 500})
	assert.Nil(t, err)
	_, err = repo.Count(token.WhereExpLessThan{Exp: 500})
	assert.Nil(t, err)
	_, err = repo.Count(token.WhereSub{Sub: "test"})
	assert.Nil(t, err)
}

func TestDelete(t *testing.T) {
	err := godotenv.Load("../../../../.env")
	assert.Nil(t, err)
	db := connectGormDB()

	repo := token.New(db.Debug())
	_, err = repo.Count(token.WhereSub{Sub: "test"}, token.WhereExpLessThan{Exp: 500})
	assert.Nil(t, err)
	_, err = repo.Count(token.WhereExpLessThan{Exp: 500})
	assert.Nil(t, err)
	_, err = repo.Count(token.WhereSub{Sub: "test"})
	assert.Nil(t, err)

	err = repo.Deletes()
	assert.Equal(t, err, gorm.ErrMissingWhereClause)
}

func TestTransaction(t *testing.T) {
	err := godotenv.Load("../../../../.env")
	assert.Nil(t, err)
	db := connectGormDB()

	repo := token.New(db.Debug())

	tx, txRepo, err := repo.StartTransaction()
	assert.Nil(t, err)
	tkn := schemas.Token{
		Sub: "test rollback",
		Exp: 6000,
	}
	txRepo.Create(&tkn)
	cnt, err := txRepo.Count()
	assert.Nil(t, err)
	fmt.Printf("txRepo.Count(): %d", cnt)
	assert.Equal(t, int64(1), cnt)

	cnt, err = repo.Count()
	assert.Nil(t, err)
	fmt.Printf("repo.Count(): %d", cnt)
	assert.Equal(t, int64(0), cnt)

	err = tx.Rollback()
	assert.Nil(t, err)

	cnt, err = repo.Count()
	assert.Nil(t, err)
	fmt.Printf("repo.Count(): %d", cnt)
	assert.Equal(t, int64(0), cnt)

	tx, txRepo, err = repo.StartTransaction()
	assert.Nil(t, err)
	tkn = schemas.Token{
		Sub: "test commit",
		Exp: 700,
	}
	txRepo.Create(&tkn)

	err = tx.Commit()
	assert.Nil(t, err)
}
