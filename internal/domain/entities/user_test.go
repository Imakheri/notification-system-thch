package entities

import (
	"reflect"
	"testing"
)

func TestNewUser(t *testing.T) {
	type args struct {
		name        string
		password    string
		email       string
		phone       string
		deviceToken string
	}
	tests := []struct {
		name           string
		args           args
		want           User
		wantErr        bool
		ExpectedErrMsg string
	}{
		{
			name: "Valid New User",
			args: args{
				name:        "Juan Pérez",
				password:    "SecurePassword123!",
				email:       "juanperez@example.com",
				phone:       "1234567890",
				deviceToken: "thisIsAValidToken",
			},
			want: User{
				ID:            0,
				Name:          "Juan Pérez",
				Password:      "SecurePassword123!",
				Email:         "juanperez@example.com",
				Phone:         "1234567890",
				DeviceToken:   "thisIsAValidToken",
				Notifications: nil,
			},
			wantErr: false,
		},
		{
			name: "Invalid New User (Invalid Name)",
			args: args{
				name:        "12345 67890",
				password:    "isnotasecurepass",
				email:       "juanperezexamplecom",
				phone:       "123456789012345",
				deviceToken: "",
			},
			want:           User{},
			wantErr:        true,
			ExpectedErrMsg: "name should only contain letters",
		},
		{
			name: "Invalid New User (Invalid Password)",
			args: args{
				name:        "Juan Pérez",
				password:    "isnotasecurepass",
				email:       "juanperezexamplecom",
				phone:       "123456789012345",
				deviceToken: "",
			},
			want:           User{},
			wantErr:        true,
			ExpectedErrMsg: "password must contain at least one upper case and one lower case character, one digit and one special character",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(tt.args.name, tt.args.password, tt.args.email, tt.args.phone, tt.args.deviceToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.ExpectedErrMsg {
				t.Errorf("NewUser() error [%v], wantErr [%v]", err.Error(), tt.ExpectedErrMsg)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	type args struct {
		name        string
		password    string
		phone       string
		deviceToken string
	}
	tests := []struct {
		name           string
		args           args
		want           User
		wantErr        bool
		ExpectedErrMsg string
	}{
		{
			name: "Valid User To Update",
			args: args{
				name:     "",
				password: "PasswordIsSecure1!",
				phone:    "0987654321",
			},
			want: User{
				Password: "PasswordIsSecure1!",
				Phone:    "0987654321",
			},
			wantErr: false,
		},
		{
			name: "Invalid User To Update (Invalid Password)",
			args: args{
				password: "0192837465",
			},
			want:           User{},
			wantErr:        true,
			ExpectedErrMsg: "password must contain at least one upper case and one lower case character, one digit and one special character",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpdateUser(tt.args.name, tt.args.password, tt.args.phone, tt.args.deviceToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.ExpectedErrMsg {
				t.Errorf("NewUser() mensaje de error = [%v], se esperaba [%v]", err.Error(), tt.ExpectedErrMsg)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
