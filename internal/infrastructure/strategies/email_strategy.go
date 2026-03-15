package strategies

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/domain/gateway"
)

type EmailStrategy struct {
	simulatedApiService gateway.SimulatedApiService
}

func NewEmailStrategy(simulatedApiService gateway.SimulatedApiService) entities.NotificationStrategy {
	return &EmailStrategy{
		simulatedApiService: simulatedApiService,
	}
}

func (es *EmailStrategy) validate(recipient entities.User) error {
	if !strings.Contains(recipient.Email, "@") || !strings.Contains(recipient.Email, ".") {
		return errors.New("invalid email structure")
	}
	return nil
}

func (es *EmailStrategy) Send(sender string, recipient entities.User, notification entities.Notification) (int, error) {
	err := es.validate(recipient)
	if err != nil {
		return http.StatusBadRequest, err
	}
	emailTemplateDTO := toEmailTemplateDTO(sender, notification)
	_, err = generateEmailTemplate(emailTemplateDTO)
	status, err := es.simulatedApiService.RandomizeHTTPStatus()
	if err != nil {
		return status, errors.New(fmt.Sprintf("%v %v", err, "notification via email"))
	}
	return status, nil
}

func generateEmailTemplate(emailTemplateDTO EmailTemplateDTO) (string, error) {
	path := filepath.Join("templates", "email_template.html")

	tpl, err := template.ParseFiles(path)
	if err != nil {
		return "", err
	}

	var tmpl bytes.Buffer
	if err := tpl.Execute(&tmpl, emailTemplateDTO); err != nil {
		return "", err
	}

	return tmpl.String(), nil
}

func toEmailTemplateDTO(sender string, notification entities.Notification) EmailTemplateDTO {
	return EmailTemplateDTO{
		Title:   notification.Title,
		Content: notification.Content,
		Sender:  sender,
	}
}

type EmailTemplateDTO struct {
	Title   string
	Content string
	Sender  string
}
