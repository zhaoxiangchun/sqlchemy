// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clickhouse

import (
	"testing"

	"yunion.io/x/pkg/errors"
	"yunion.io/x/sqlchemy"
)

func insertSqlPrep(v interface{}, update bool) (string, []interface{}, error) {
	sqlchemy.SetDBWithNameBackend(nil, sqlchemy.DefaultDB, sqlchemy.ClickhouseBackend)
	ts := sqlchemy.NewTableSpecFromStruct(v, "vv")
	results, err := ts.InsertSqlPrep(v, update)
	if err != nil {
		return "", nil, errors.Wrap(err, "InsertSqlPrep")
	}
	return results.Sql, results.Values, err
}

func TestInsertAutoIncrement(t *testing.T) {
	cases := []struct {
		value   interface{}
		update  bool
		wantSQL string
		wantVar int
	}{
		{
			value: &struct {
				RowId int `auto_increment:"true"`
			}{
				RowId: 12345,
			},
			update:  false,
			wantSQL: "INSERT INTO `vv` (`row_id`) VALUES (?)",
			wantVar: 1,
		},
		{
			value: &struct {
				RowId int    `primary:"true"`
				Name  string `width:"24"`
			}{
				RowId: 1,
				Name:  "a",
			},
			update:  true,
			wantSQL: "",
			wantVar: 3,
		},
	}
	for _, c := range cases {
		sql, vals, err := insertSqlPrep(c.value, c.update)
		if err != nil {
			t.Errorf("prepare sql failed: %s", err)
		} else {
			if sql != c.wantSQL {
				t.Errorf("sql want %s got %s", c.wantSQL, sql)
			} else {
				if len(vals) != c.wantVar {
					t.Errorf("vars want %d got %d", c.wantVar, len(vals))
				}
			}
		}
	}
}

func TestInsertWithPointerValue(t *testing.T) {
	sql, vals, err := insertSqlPrep(&struct {
		RowId int `auto_increment:"true"`
		ColT1 *int
		ColT2 int
		ColT3 string
		ColT4 *string
	}{}, false)
	if err != nil {
		t.Errorf("prepare sql failed: %s", err)
		return
	}
	t.Logf("%s values: %v", sql, vals)
}
