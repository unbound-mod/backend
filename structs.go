package main

import (
	"gorm.io/gorm"
)

/* Users */
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt int            `json:"created_at"`
	UpdatedAt int            `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Discord DiscordUser `gorm:"embedded;embeddedPrefix:discord_" json:"discord"`
	Tokens  UserTokens  `json:"-" gorm:"embedded"`
}

type UserTokens struct {
	AccessToken  string `gorm:"column:access_token"`
	RefreshToken string `gorm:"column:refresh_token"`
}

type DiscordUser struct {
	ID          string  `json:"id"`
	Username    string  `json:"username"`
	Avatar      *string `json:"avatar"`
	Flags       int     `json:"flags"`
	Banner      *string `json:"banner"`
	AccentColor *string `json:"accent_color"`
	DisplayName string  `json:"global_name"`
	Email       string  `json:"email"`
}

/* Plugins */
type Plugin struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt int            `json:"created_at"`
	UpdatedAt int            `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Manifest Manifest `json:"manifest" gorm:"type:json"`
}

type Manifest struct {
	ID string `json:"id" gorm:"column:id"`
}

/* Types */
type AuthorizeErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type AuthorizeSuccessResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}
