package listeners

import (
	"ferdinand/app/models"
	"ferdinand/util"

	"github.com/caesar-rocks/vexillum"
	"gorm.io/gorm"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/client"
)

type UsersListener struct {
	stripe   *client.API
	db       *gorm.DB
	vexillum *vexillum.Vexillum
}

func NewUsersListener(stripe *client.API, db *gorm.DB, vexillum *vexillum.Vexillum) *UsersListener {
	return &UsersListener{stripe, db, vexillum}
}

func (usersListener *UsersListener) OnCreated(msg *message.Message) ([]*message.Message, error) {
	var user models.User
	if err := util.DecodeJSON(msg.Payload, &user); err != nil {
		return nil, err
	}

	if usersListener.vexillum.IsActive("billing") {
		if err := usersListener.assignStripeCustomer(&user); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (usersListener *UsersListener) assignStripeCustomer(user *models.User) error {
	cus, err := usersListener.stripe.Customers.New(&stripe.CustomerParams{
		Email: stripe.String(user.Email),
		Name:  stripe.String(user.FullName),
	})
	if err != nil {
		return err
	}

	user.StripeCustomerID = cus.ID

	if err := usersListener.db.Model(user).Update("stripe_customer_id", cus.ID).Error; err != nil {
		return err
	}

	return nil
}
