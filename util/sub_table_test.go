package util

import (
	"gorm.io/gorm"
	"testing"
)

type UserWealthLog struct {
	Uid    int `gorm:"column:Uid"`
	Wealth int `gorm:"column:Wealth"`
}

func TestBranchQuery(t *testing.T) {
	type args struct {
		db          *gorm.DB
		startTime   any
		endTime     any
		tablePrefix string
		timeColumn  string
		pageNo      int32
		pageSize    int32
		column      []string
		columnValue []any
		total       *int64
		result      any
		tableSuffix []string
		sort        string
	}
	column := make([]string, 0)
	columnValue := make([]any, 0)
	column = append(column, "Uid = ?")
	columnValue = append(columnValue, 20087376)
	total := int64(0)
	result := make([]*UserWealthLog, 0)
	tableSuffix := GetSectionMonth("2023-11-25 00:00:00", "2023-12-25 23:59:59")

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{

			name: "",
			args: args{
				db:          nil,
				startTime:   1700841600,
				endTime:     1703951999,
				tablePrefix: "user_wealth_log_",
				timeColumn:  "LogTime",
				pageNo:      1,
				pageSize:    20,
				column:      column,
				columnValue: columnValue,
				total:       &total,
				result:      result,
				tableSuffix: tableSuffix,
				sort:        "desc",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BranchQuery(tt.args.db, tt.args.startTime, tt.args.endTime, tt.args.tablePrefix, tt.args.timeColumn, tt.args.pageNo, tt.args.pageSize, tt.args.column, tt.args.columnValue, tt.args.total, tt.args.result, tt.args.tableSuffix, tt.args.sort); (err != nil) != tt.wantErr {
				t.Errorf("BranchQuery() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
