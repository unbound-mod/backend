package Structures

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Addon struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt int            `json:"created_at"`
	UpdatedAt int            `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Type    string `json:"type" gorm:"column:type"`
	AddonID string `json:"addon_id" gorm:"column:addon_id"`

	Name        string         `json:"name" gorm:"column:name"`
	Description string         `json:"description" gorm:"column:description"`
	Icon        string         `json:"icon" gorm:"column:icon"`
	Images      pq.StringArray `json:"images" gorm:"type:text[];column:images"`
	Changelog   string         `json:"changelog" gorm:"column:changelog"`
	// README      string   `json:"readme" gorm:"column:readme"`
	// "readme": "https://raw.githubusercontent.com/unbound-mod/sources/main/test/Plumpy/README.md",
	// "manifest": "https://raw.githubusercontent.com/unbound-mod/icons/main/packs/Plumpy/manifest.json"
	// Manifest Manifest `json:"manifest" gorm:"type:json"`
}

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

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt int            `json:"created_at"`
	UpdatedAt int            `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Developer bool        `gorm:"column:developer" json:"developer"`
	Discord   DiscordUser `gorm:"embedded;embeddedPrefix:discord_" json:"discord"`
	Tokens    UserTokens  `json:"-" gorm:"embedded;embeddedPrefix:tokens_"`
}

type UserTokens struct {
	AccessToken  string `gorm:"column:access_token"`
	RefreshToken string `gorm:"column:refresh_token"`
}

type DiscordUser struct {
	ID          string `gorm:"column:id" json:"id"`
	Username    string `gorm:"column:username" json:"username"`
	Avatar      string `gorm:"column:avatar" json:"avatar"`
	DisplayName string `gorm:"column:global_name" json:"global_name"`
	Email       string `gorm:"column:email" json:"email"`
}
