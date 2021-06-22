package mailer

import (
	"testing"

	"github.com/my-companies-be/models"
)

func TestSendActiveAccount(t *testing.T) {
	var u models.User
	u.RegisterToken = "ihPVUFsZU4-cRQF-nPpQ"
	u.Email = "329176418@qq.com"
	SendActiveAccount(&u)
}
