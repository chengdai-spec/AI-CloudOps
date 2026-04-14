/*
 * MIT License
 *
 * Copyright (c) 2024 Bamboo
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 *
 */

package di

import (
	"fmt"
	"time"

	"github.com/GoSimplicity/AI-CloudOps/internal/config"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config, logger *zap.Logger) (*gorm.DB, error) {
	if cfg == nil {
		return nil, fmt.Errorf("配置不能为空")
	}
	if logger == nil {
		logger = zap.NewNop()
	}

	addr := cfg.MySQL.Addr
	db, err := gorm.Open(mysql.Open(addr), &gorm.Config{})
	if err != nil {
		logger.Error("数据库连接失败", zap.Error(err))
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("获取sql.DB失败", zap.Error(err))
		return nil, fmt.Errorf("获取sql.DB失败: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		logger.Error("数据库ping失败", zap.Error(err))
		return nil, fmt.Errorf("数据库ping失败: %w", err)
	}

	if err := InitTables(db); err != nil {
		logger.Error("初始化数据库表失败", zap.Error(err))
		return nil, fmt.Errorf("初始化数据库表失败: %w", err)
	}

	logger.Info("数据库初始化成功")
	return db, nil
}

func CheckDBHealth(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("数据库连接为空")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %v", err)
	}

	var pingErr error
	for i := 0; i < 3; i++ {
		pingErr = sqlDB.Ping()
		if pingErr == nil {
			return nil
		}
		if i < 2 { // 只在前两次失败后等待
			time.Sleep(10 * time.Second)
		}
	}

	if pingErr != nil {
		return fmt.Errorf("数据库ping失败: %v", pingErr)
	}
	return nil
}

func IsDBAvailable(db *gorm.DB) bool {
	return CheckDBHealth(db) == nil
}
