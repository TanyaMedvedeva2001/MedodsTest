package mail

import (
	"authTestMedods/auth/types"
	"fmt"
	"net/smtp"
)

// Manager - структура email менеджера
type Manager struct {
	email    string
	password string
}

var mailManager *Manager

// InitMailManager - инициализация email менеджера
func InitMailManager(email string, password string) {
	if mailManager == nil {
		mailManager = &Manager{email, password}
	}
}

// GetMailManager - геттер email менеджера
func GetMailManager() *Manager {
	return mailManager
}

// SendMailNotification - отправка письма на почту
func (mm *Manager) SendMailNotification(identity types.UserSessionIdentity) error {
	to := []string{identity.Email}
	smptHost := "smtp.gmail.com"
	smptPort := ":587"
	message := []byte(fmt.Sprintf("Уважаемый пользователь (guid: %v), замечена активность на вашей странице с другого ip адреса: %v", identity.GUID, identity.Ip))
	auth := smtp.PlainAuth("", mm.email, mm.password, smptHost)
	err := smtp.SendMail(smptHost+smptPort, auth, mm.email, to, message)
	if err != nil {
		return fmt.Errorf("не удалось отправить письмо: %v", err)
	}
	return nil
}
