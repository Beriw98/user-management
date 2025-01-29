package entity

import (
	"entgo.io/ent"
	"reflect"
	"testing"
)

func TestUser_Fields(t *testing.T) {
	type fields struct {
		Schema ent.Schema
	}
	tests := []struct {
		name   string
		fields fields
		want   []ent.Field
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := User{
				Schema: tt.fields.Schema,
			}
			if got := us.Fields(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fields() = %v, want %v", got, tt.want)
			}
		})
	}
}
