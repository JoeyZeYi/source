package util

import (
	"fmt"
	"github.com/JoeyZeYi/source/log"
	"github.com/JoeyZeYi/source/log/zap"
	"gorm.io/gorm"
)

type IBranchQueryTable interface {
	TablePrefix() string
}

// BranchQuery 分表查询
/*
分表查询
采用年份加月份分表查询例如202001 202002
startTime 开始时间 2020-07-01 00:00:00
endTime   结束时间 2020-09-01 00:00:00
tablePrefix 表前缀  user_
timeColumn  表中的创建时间字段
pageNo    pageSize  分页字段
column    动态查询的字段名
columnValue  动态查询的字段值
total 总数
result    结果集
tableSuffix 表后缀数组
sort 排序  DESC降序或ASC升序
*/
func BranchQuery(db *gorm.DB, startTime, endTime any, tablePrefix, timeColumn string, pageNo, pageSize int32, column []string, columnValue []any, total *int64, result any, tableSuffix []string, sort string) error {
	log.Debug("BranchQuery", zap.Any("startTime", startTime),
		zap.Any("endTime", endTime), zap.Any("tablePrefix", tablePrefix),
		zap.Any("timeColumn", timeColumn), zap.Any("pageNo", pageNo),
		zap.Any("pageSize", pageSize), zap.Any("column", column), zap.Any("columnValue", columnValue),
		zap.Any("tableSuffix", tableSuffix), zap.Any("sort", sort))

	sql := ""
	totalSql := ""
	args := make([]interface{}, 0)
	where := "  where 1=1 "
	//动态添加查询条件 目前不支持模糊查询,如需模糊查询。可在value!=""判断后在判断列名的方法拼接 args = append(args, "%"+value+"%")
	for i, value := range columnValue {
		if value != nil {
			where += " and " + column[i]
			args = append(args, value)
		}
	}
	for i, v := range tableSuffix {
		for _, value := range columnValue {
			if i > 0 {
				//动态拼接子查询中的查询条件
				if value != nil {
					args = append(args, value)
				}
			}
		}
		timeWhere := ""
		if startTime != nil && endTime != nil {
			args = append(args, startTime)
			args = append(args, endTime)
			timeWhere = fmt.Sprintf(" and %s >= ? and  %s <= ?", timeColumn, timeColumn)
		}
		tableSql := "select * from " + tablePrefix + v
		countSql := "select count(1) count from " + tablePrefix + v
		if i > 0 {
			sql = sql + " union " + tableSql + where + timeWhere
			totalSql = totalSql + " union all " + countSql + where + timeWhere
		} else {
			sql = sql + tableSql + where + timeWhere
			totalSql = totalSql + countSql + where + timeWhere
		}
	}
	if len(tableSuffix) > 1 {
		sql = fmt.Sprintf("select * from ( %s ) a order by a.%s %s  limit ?,?; ", sql, timeColumn, sort)
		totalSql = "select sum(count) count from (" + totalSql + ") a "
	} else {
		sql += fmt.Sprintf(" order by %s %s  limit ?,?; ", timeColumn, sort)
	}
	err := db.Raw(totalSql, args...).Count(total).Error
	if err != nil {
		log.Error("BranchQuery-Count", zap.Error(err))
		return err
	}
	args = append(args, getPage(pageNo, pageSize), pageSize)
	err = db.Raw(sql, args...).Scan(result).Error
	if err != nil {
		log.Error("BranchQuery-Scan", zap.Error(err))
	}
	log.Debug("BranchQuery", zap.Any("result", result))
	return err
}

// BranchQuerySum
/*
用于计算分表查询中的SUM总数查询、不加分页
采用年份加月份分表查询例如202001 202002
startTime 开始时间 2020-07-01 00:00:00
endTime   结束时间 2020-09-01 00:00:00
tablePrefix 表前缀  user_
timeColumn  表中的创建时间字段
pageNo    pageSize  分页字段
column    动态查询的字段名
columnValue  动态查询的字段值
sumColumn 需要被合计的字段名
result    结果集
tableSuffix 表后缀数组
*/
func BranchQuerySum(db *gorm.DB, startTime, endTime, tablePrefix, timeColumn,
	sumColumn string, column []string,
	columnValue []interface{}, result int64,
	tableSuffix []string) error {

	log.Debug("BranchQuery", zap.Any("startTime", startTime),
		zap.Any("endTime", endTime), zap.Any("tablePrefix", tablePrefix),
		zap.Any("timeColumn", timeColumn), zap.Any("column", column), zap.Any("columnValue", columnValue),
		zap.Any("tableSuffix", tableSuffix))

	totalSql := ""
	args := make([]interface{}, 0)
	where := "  where 1=1 "
	//动态添加查询条件 目前不支持模糊查询,如需模糊查询。可在value!=""判断后在判断列名的方法拼接 args = append(args, "%"+value+"%")
	for i, value := range columnValue {
		if value != nil {
			where += " and " + column[i]
			args = append(args, value)
		}
	}
	for i, v := range tableSuffix {
		for _, value := range columnValue {
			if i > 0 {
				//动态拼接子查询中的查询条件
				if value != nil {
					args = append(args, value)
				}
			}
		}
		timeWhere := ""
		if len(startTime) > 0 && len(endTime) > 0 {
			args = append(args, startTime)
			args = append(args, endTime)
			timeWhere = fmt.Sprintf(" and %s >= ? and  %s <= ?", timeColumn, timeColumn)
		}
		countSql := "select IFNULL(sum(" + sumColumn + "),0)  count from " + tablePrefix + v
		if i > 0 {
			totalSql = totalSql + " union all " + countSql + where + timeWhere
		} else {
			totalSql = totalSql + countSql + where + timeWhere
		}
	}
	if len(tableSuffix) > 1 {
		totalSql = fmt.Sprintf("select IFNULL(sum(count),0) count from (" + totalSql + ") a ")
	}
	err := db.Raw(totalSql, args...).Count(&result).Error
	if err != nil {
		log.Error("BranchQuery-Count", zap.Error(err))
	}
	log.Debug("BranchQuery", zap.Any("result", result))
	return err
}

func getPage(pageNo, pageSize int32) int32 {
	if pageNo == 1 {
		pageNo = pageNo - 1
	} else {
		pageNo = (pageNo - 1) * pageSize
	}
	return pageNo
}
