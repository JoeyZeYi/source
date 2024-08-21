package excel

import (
	"github.com/xuri/excelize/v2"
	"reflect"
	"testing"
)

type ExcelDemo struct {
	Uid      string `excel:"用户ID"`
	NickName string `excel:"昵称"`
}

func TestCreateExcel(t *testing.T) {
	type args[T any] struct {
		rows []*T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want *excelize.File
	}
	var tests []testCase[ExcelDemo]
	tests = append(tests, testCase[ExcelDemo]{
		name: "",
		args: args[ExcelDemo]{
			rows: []*ExcelDemo{},
		},
		want: nil,
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateExcel(tt.args.rows); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateExcel() = %v, want %v", got, tt.want)
			}
		})
	}
}
