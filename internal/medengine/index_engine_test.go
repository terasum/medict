package medengine

import "testing"

func TestIndexEngine_AddRecord(t *testing.T) {
	type fields struct {
		indexFilePath string
	}
	type args struct {
		record *IndexRecord
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "t1",
			fields: struct{ indexFilePath string }{indexFilePath: "./testdata/testidx.meidx"},
			args: args{record: &IndexRecord{
				keyWord:           "0test",
				keyBlockIndex:     0,
				recordStartOffset: 1231,
				recordEndOffset:   1233,
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine, err := NewEngine(tt.fields.indexFilePath)
			if err != nil {
				t.Fatal(err)
			}
			db, err := engine.Acquire()
			if err != nil {
				t.Fatal(err)
			}

			defer engine.Release(db)

			if err = engine.AddRecord(tt.args.record); (err != nil) != tt.wantErr {
				t.Errorf("AddRecord() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}
