package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/imakheri/notifications-thch/internal/domain/usecase"
	"github.com/imakheri/notifications-thch/internal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDeleteNotificationHandler(t *testing.T) {
	type fields struct {
		deleteNotificationUseCase func(ctrl *gomock.Controller) usecase.DeleteNotificationUseCase
	}
	type args struct {
		notificationID string
		authUserID     uint
		hasAuthContext bool
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantStatus int
		wantBody   string
	}{
		{
			name: "Notification Deleted Successfully",
			fields: fields{
				deleteNotificationUseCase: func(ctrl *gomock.Controller) usecase.DeleteNotificationUseCase {
					m := mocks.NewMockDeleteNotificationUseCase(ctrl)
					m.EXPECT().Exec(uint(1), uint(10)).Return(uint(10), nil)
					return m
				},
			},
			args: args{
				notificationID: "10",
				authUserID:     1,
				hasAuthContext: true,
			},
			wantStatus: http.StatusOK,
			wantBody:   `{"id":10, "message":"Notification deleted successfully"}`,
		},
		{
			name: "Invalid Notification ID Format",
			fields: fields{
				deleteNotificationUseCase: func(ctrl *gomock.Controller) usecase.DeleteNotificationUseCase {
					m := mocks.NewMockDeleteNotificationUseCase(ctrl)
					return m
				},
			},
			args: args{
				notificationID: "abc",
				authUserID:     1,
				hasAuthContext: false,
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"strconv.Atoi: parsing \"abc\": invalid syntax"}`,
		},
		{
			name: "Invalid Notification ID Length",
			fields: fields{
				deleteNotificationUseCase: func(ctrl *gomock.Controller) usecase.DeleteNotificationUseCase {
					m := mocks.NewMockDeleteNotificationUseCase(ctrl)
					return m
				},
			},
			args: args{
				notificationID: "",
				authUserID:     1,
				hasAuthContext: false,
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"must enter a valid id"}`,
		},
		{
			name: "Error getting UserID from token",
			fields: fields{
				deleteNotificationUseCase: func(ctrl *gomock.Controller) usecase.DeleteNotificationUseCase {
					m := mocks.NewMockDeleteNotificationUseCase(ctrl)
					return m
				},
			},
			args: args{
				notificationID: "1",
				authUserID:     0,
				hasAuthContext: false,
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error":"can not get user id"}`,
		},
		{
			name: "Error getting UserID from token",
			fields: fields{
				deleteNotificationUseCase: func(ctrl *gomock.Controller) usecase.DeleteNotificationUseCase {
					m := mocks.NewMockDeleteNotificationUseCase(ctrl)
					m.EXPECT().Exec(uint(1), uint(999)).Return(uint(0), errors.New("notification not found"))
					return m
				},
			},
			args: args{
				notificationID: "999",
				authUserID:     1,
				hasAuthContext: true,
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error":"notification not found"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			ctx.Params = []gin.Param{{Key: "id", Value: tt.args.notificationID}}

			if tt.args.hasAuthContext {
				ctx.Set("user_id", tt.args.authUserID)
			}

			mockUseCase := tt.fields.deleteNotificationUseCase(ctrl)
			handler := DeleteNotificationHandler(mockUseCase)
			handler(ctx)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.JSONEq(t, tt.wantBody, w.Body.String())
		})
	}
}
