package handler

import (
	"github.com/labstack/echo/v4"
	"reflect"
	"testing"
)

func TestNewUserHTTPHandler(t *testing.T) {
	type args struct {
		repository userRepository
	}
	tests := []struct {
		name string
		args args
		want *UserHTTPHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserHTTPHandler(tt.args.repository); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserHTTPHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserHTTPHandler_Create(t *testing.T) {
	type fields struct {
		userRepository userRepository
	}
	type args struct {
		ec echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UserHTTPHandler{
				userRepository: tt.fields.userRepository,
			}
			if err := h.Create(tt.args.ec); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserHTTPHandler_Delete(t *testing.T) {
	type fields struct {
		userRepository userRepository
	}
	type args struct {
		ec echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UserHTTPHandler{
				userRepository: tt.fields.userRepository,
			}
			if err := h.Delete(tt.args.ec); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserHTTPHandler_GetByID(t *testing.T) {
	type fields struct {
		userRepository userRepository
	}
	type args struct {
		ec echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UserHTTPHandler{
				userRepository: tt.fields.userRepository,
			}
			if err := h.GetByID(tt.args.ec); (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserHTTPHandler_GetMany(t *testing.T) {
	type fields struct {
		userRepository userRepository
	}
	type args struct {
		ec echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UserHTTPHandler{
				userRepository: tt.fields.userRepository,
			}
			if err := h.GetMany(tt.args.ec); (err != nil) != tt.wantErr {
				t.Errorf("GetMany() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserHTTPHandler_Update(t *testing.T) {
	type fields struct {
		userRepository userRepository
	}
	type args struct {
		ec echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UserHTTPHandler{
				userRepository: tt.fields.userRepository,
			}
			if err := h.Update(tt.args.ec); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserHTTPHandler_UpdatePassword(t *testing.T) {
	type fields struct {
		userRepository userRepository
	}
	type args struct {
		ec echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UserHTTPHandler{
				userRepository: tt.fields.userRepository,
			}
			if err := h.UpdatePassword(tt.args.ec); (err != nil) != tt.wantErr {
				t.Errorf("UpdatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
