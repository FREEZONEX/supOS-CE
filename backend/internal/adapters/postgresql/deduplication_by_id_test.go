package postgresql

import (
	"backend/internal/types"
	"encoding/json"
	"testing"
)

func TestQuChong(t *testing.T) {
	/*


		{"@timestamp":"2025-12-31T03:07:18.970Z","caller":"msg_consumer/UnsMessageConsumer.go:262","content":"MergeList[ _opcda_file_22627e78e18847c9be69 ]:
		[map[boolean:false float:19346.25 int1:18042 long:2344752 quality:0 str:String 1 timeStamp:1.76715043755e+12]
		map[dt:1767150436655 quality:0 timeStamp:1.767150437548e+12] map[double:1793.75 quality:0 timeStamp:1.76715043755e+12]]

		, list: [map[float:19346.25 quality:0 timeStamp:1.76715043755e+12] map[long:2344752 quality:0 timeStamp:1.76715043755e+12]
		map[quality:0 str:String 1 timeStamp:1.76715043755e+12] map[boolean:false quality:0 timeStamp:1.76715043755e+12]
		map[int1:18042 quality:0 timeStamp:1.76715043755e+12] map[dt:1767150436655 quality:0 timeStamp:1.767150437548e+12]
		map[double:1793.75 quality:0 timeStamp:1.76715043755e+12]]\n","level":"debug"}


			{"@timestamp":"2025-12-31T03:07:19.298Z","caller":"postgresql/pg_persistence.go:82","content":"insert sql:
			INSERT INTO \"_opcda_file_22627e78e18847c9be69\" AS t(\"timeStamp\",\"int1\",\"long\",\"float\",\"double\",\"boolean\",\"dt\",\"str\",\"quality\")
			VALUES ($1,DEFAULT,DEFAULT,DEFAULT,$2,DEFAULT,NOW(),DEFAULT,$3),($4,DEFAULT,DEFAULT,DEFAULT,DEFAULT,DEFAULT,$5,DEFAULT,$6)
			ON CONFLICT (\"timeStamp\") DO UPDATE SET \"int1\"  = COALESCE(EXCLUDED.\"int1\", t.\"int1\"), \"long\"  =
			COALESCE(EXCLUDED.\"long\", t.\"long\"), \"float\"  = COALESCE(EXCLUDED.\"float\", t.\"float\"),
			\"double\"  = COALESCE(EXCLUDED.\"double\", t.\"double\"), \"boolean\"  = COALESCE(EXCLUDED.\"boolean\", t.\"boolean\"),
			\"dt\"  = COALESCE(EXCLUDED.\"dt\", t.\"dt\"), \"str\"  = COALESCE(EXCLUDED.\"str\", t.\"str\"), \"quality\"
			= COALESCE(EXCLUDED.\"quality\", t.\"quality\"), values: [2025-12-31T03:07:17Z 1793.75 0 2025-12-31T03:07:17Z 2025-12-31T03:07:16Z 0]","level":"debug"}
			{"@timestamp":"2025-12-31T03:07:24.303Z","caller":"timescaledb/TsdbPersistentService.go:60","content":"persistence fail处理完成，但有错误: 批次操作失败[_opcda_file_22627e78e18847c9be69]: ERROR: ON CONFLICT DO UPDATE command cannot affect row a second time (SQLSTATE 21000)  8","level":"error"}

	*/
	var uns = &types.CreateTopicDto{Alias: "_opcda_file_22627e78e18847c9be69"}
	jsonfs := `[{"name":"timeStamp","type":"DATETIME","unique":true,"systemField":true},{"name":"int1","type":"INTEGER"},{"name":"long","type":"LONG"},{"name":"float","type":"FLOAT"},{"name":"double","type":"DOUBLE"},{"name":"boolean","type":"BOOLEAN"},{"name":"dt","type":"DATETIME"},{"name":"str","type":"STRING","maxLen":512},{"name":"quality","type":"LONG","systemField":true}]`
	json.Unmarshal([]byte(jsonfs), &uns.Fields)
	data := []map[string]string{
		{
			"boolean":   "false",
			"float":     "19346.25",
			"int1":      "18042",
			"long":      "2344752",
			"timeStamp": "1.76715043755e+12",
		},
		{
			"dt":        "1767150436655",
			"timeStamp": "1.767150437548e+12",
		},
		{
			"double":    "1793.75",
			"timeStamp": "1.76715043755e+12",
		},
	}
	sql, params := getInsertStatement(uns, data)
	t.Log("sql: ", sql)
	t.Logf("params: %+v", params)
	rs := DeduplicationById(uns, data)
	for i, v := range rs {
		t.Log(i, ": ", v)
	}
	/**
		 INSERT INTO "_opcda_file_22627e78e18847c9be69" AS t("timeStamp","int1","long","float","double","boolean","dt","str","quality")
	VALUES ('2025-12-31T03:07:17.55Z',DEFAULT,DEFAULT,DEFAULT,'1793.75',DEFAULT,NOW(),DEFAULT,DEFAULT),('2025-12-31T03:07:17.548Z',DEFAULT,DEFAULT,DEFAULT,DEFAULT,DEFAULT,'2025-12-31T03:07:16.655Z',DEFAULT,DEFAULT)
	ON CONFLICT ("timeStamp") DO UPDATE SET "int1"  = COALESCE(EXCLUDED."int1", t."int1"), "long"  = COALESCE(EXCLUDED."long", t."long"), "float"  = COALESCE(EXCLUDED."float", t."float"), "double"  = COALESCE(EXCLUDED."double", t."double"), "boolean"  = COALESCE(EXCLUDED."boolean", t."boolean"), "dt"  = COALESCE(EXCLUDED."dt", t."dt"), "str"  = COALESCE(EXCLUDED."str", t."str"), "quality"  = COALESCE(EXCLUDED."quality", t."quality")


	*/

}
