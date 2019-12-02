package context

import (
	"context"

	"github.com/lenslocked/models"
)

const (
	userkey privateKey = "user"
)

type privateKey string

func WithUser(ctx context.Context, user *models.User) context.Context {
	return contex.WithValue(ctx, userkey, user)
}

func User(ctx context.Context) *models.User {
	if temp := ctx.Value(userkey); temp != nil {
		if user, ok := temp.(*models.User); ok {
			return user
		}
	}
	return nil
}
