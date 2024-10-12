package dbconn

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"

	"wallet/pkg/common/config"
)

func NewMysqlGormDB() (*gorm.DB, error) {
	var (
		replicasList []gorm.Dialector
		sourcesList  []gorm.Dialector
	)

	logConfig := &gorm.Config{
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	if config.Config.App.Env == "dev" {
		logConfig.Logger = gormlogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			gormlogger.Config{
				Colorful: true,
				LogLevel: gormlogger.Info,
			})
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
			*config.Config.MysqlMaster.Username,
			*config.Config.MysqlMaster.Password,
			(*config.Config.MysqlMaster.Address)[0],
			*config.Config.MysqlMaster.Database),
		DefaultStringSize: 255,
	}), logConfig)

	if err != nil {
		return nil, err
	}

	for _, item := range *config.Config.MysqlSlave.Address {
		dbDsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			*config.Config.MysqlSlave.Username, *config.Config.MysqlSlave.Password, item, *config.Config.MysqlSlave.Database)
		replicasList = append(replicasList, mysql.Open(dbDsn))
	}

	for _, item := range *config.Config.MysqlMaster.Address {
		dbDsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			*config.Config.MysqlMaster.Username, *config.Config.MysqlMaster.Password, item, *config.Config.MysqlMaster.Database)
		sourcesList = append(sourcesList, mysql.Open(dbDsn))
	}

	err = db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  sourcesList,
		Replicas: replicasList,
		Policy:   dbresolver.RandomPolicy{},
		// TraceResolverMode: true,
		TraceResolverMode: false,
	}).SetConnMaxLifetime(time.Second * time.Duration(*config.Config.MysqlMaster.MaxLifeTime)).
		SetMaxOpenConns(*config.Config.MysqlMaster.MaxOpenConn).
		SetMaxIdleConns(*config.Config.MysqlMaster.MaxIdleConn))

	return db, nil
}

func NewGormDB() (*gorm.DB, error) {
	return NewMysqlGormDB()
}
