package typescript

import (
	"reflect"
	"testing"
	"time"
)

type TestCustomStruct struct {
	ID            uint64 `json:"id"`
	secret        string
	AnotherSecret string    `json:"-"`
	Slice         []float32 `json:"slice"`
}

type TestUser struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	Age              int
	Score            int    `json:"score,string"`
	Admin            bool   `json:"admin,omitempty"`
	Password         string `json:"-"`
	Roles            []string
	VoteCount        uint64
	CreatedAt        time.Time
	DeletedAt        *time.Time
	secret           string
	Data             interface{}
	NumberSlice      []uint64 `json:"numberSlice,omitempty"`
	TestCustomStruct TestCustomStruct
	MapAny           map[string]any
	MapNumber        map[uint64]float32 `json:",omitempty"`
	MapNumberSlice   map[string][]int32
	MapPointer       map[string]*string
}

func TestGenerate(t *testing.T) {
	// Unexported fields should not be used
	_ = TestUser{}.secret
	_ = TestUser{}.TestCustomStruct.secret

	type args struct {
		goStructs map[string]interface{}
	}

	tests := []struct {
		name string
		args args
		want Interfaces
	}{
		{
			name: "basic test",
			args: args{
				goStructs: map[string]interface{}{
					"TestUserCustomizedByMap": TestUser{},
				},
			},
			want: Interfaces{
				{
					Name: "TestUserCustomizedByMap",
					Fields: []Field{
						{
							Name: "id",
							Type: "number",
						},
						{
							Name: "name",
							Type: "string",
						},
						{
							Name: "Age",
							Type: "number",
						},
						{
							Name: "score",
							Type: "string",
						},
						{
							Name:     "admin",
							Type:     "boolean",
							Optional: true,
						},
						{
							Name: "Roles",
							Type: "string[]",
						},
						{
							Name: "VoteCount",
							Type: "number",
						},
						{
							Name: "CreatedAt",
							Type: "string",
						},
						{
							Name: "DeletedAt",
							Type: "string",
							Null: true,
						},
						{
							Name: "Data",
							Type: "any",
						},
						{
							Name:     "numberSlice",
							Type:     "number[]",
							Optional: true,
						},
						{
							Name: "TestCustomStruct",
							Type: "TestCustomStruct",
						},
						{
							Name: "MapAny",
							Type: "Record<string, any>",
						},
						{
							Name:     "MapNumber",
							Type:     "Record<number, number>",
							Optional: true,
						},
						{
							Name: "MapNumberSlice",
							Type: "Record<string, number[]>",
						},
						{
							Name: "MapPointer",
							Type: "Record<string, string | null>",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Generate(tt.args.goStructs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}
