package test_tool

import (
	"database/sql"
	"fmt"
	"testing"

	gormV1 "github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	gormV2 "gorm.io/gorm"

	gormutil "github.com/geeksmy/go-libs/gormutil_v2"
)

func MockedGORMV1DBForTest(t *testing.T, sqlDB *sql.DB) *gormV1.DB {
	gormDB, err := gormV1.Open("postgres", sqlDB)
	if err != nil {
		t.Error(err)
	}

	return gormDB
}

// Deprecated: Use MockedGORMV1DBForTest instead.
func MockedGORMDBForTest(t *testing.T, sqlDB *sql.DB) *gormV1.DB {
	gormDB, err := gormV1.Open("postgres", sqlDB)
	if err != nil {
		t.Error(err)
	}

	return gormDB
}

func MockedGORMV1DBFuncForTest(t *testing.T, sqlDB *sql.DB) func() *gormV1.DB {
	gormDB := MockedGORMV1DBForTest(t, sqlDB)
	return func() *gormV1.DB { return gormDB }
}

// Deprecated: Use MockedGORMV1DBFuncForTest instead.
func MockedGORMDBFuncForTest(t *testing.T, sqlDB *sql.DB) func() *gormV1.DB {
	gormDB := MockedGORMV1DBForTest(t, sqlDB)
	return func() *gormV1.DB { return gormDB }
}

var GormV2DbCnnForTest = func(t *testing.T, envDsnKey string) *gormV2.DB {
	return DbCnnForTest(t, envDsnKey)
}

// Deprecated: Use GormV2DbCnnForTest instead.
func DbCnnForTest(t *testing.T, envDsnKey string) *gormV2.DB {
	viper.AutomaticEnv()
	dsn := viper.GetString(envDsnKey)
	db, err := gormutil.ConnectWithDSN(dsn, gormutil.Conf{LogMode: true})
	assert.NoError(t, err)
	return db.Debug()

}

var GromV2CleanTableForTest = func(t *testing.T, db *gormV2.DB, tableNames []string) {
	CleanTableForTest(t, db, tableNames)
}

// 清空表数据
// Deprecated: Use GromV2CleanTableForTest instead.
func CleanTableForTest(t *testing.T, db *gormV2.DB, tableNames []string) {
	sql := "TRUNCATE TABLE "

	for k, v := range tableNames {
		if len(tableNames)-1 == k {
			// 最后一条数据,以分号结尾
			sql += fmt.Sprintf("%s;", v)
		} else {
			sql += fmt.Sprintf("%s, ", v)
		}
	}
	err := db.Exec(sql).Error
	assert.NoError(t, err)
}
