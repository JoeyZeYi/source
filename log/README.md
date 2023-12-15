### 使用gormLog输出sql
```go
gormLogger := log.NewGormLogger(log.GormLoggerSlowThreshold(time.Second), log.GormLoggerLevel(gormlogger.Info))
	db, err := gorm.Open(directory, &gorm.Config{
		Logger: gormLogger,
	})
```
