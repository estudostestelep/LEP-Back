package models

type User2 struct {
	Id            int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string `json:"name"`
	Email         string `gorm:"unique" json:"email"`
	EmailVerified string `json:"email_verified"`
	Image         string `json:"image"`
	Access        string `json:"access"`
	GroupMember   string `json:"group_member"`
	Password      string `json:"password"`
}

type BannedLists struct {
	BannedListId int    `gorm:"primaryKey;autoIncrement" json:"banned_list_id"`
	Token        string `gorm:"type:varchar(300)" json:"token"`
	Date         string `gorm:"type:varchar(300)" json:"date"`
}

type Groups struct {
	GroupId int `gorm:"primaryKey;autoIncrement" json:"group_id"`
}

type LoggedLists struct {
	LoggedListId int    `gorm:"primaryKey;autoIncrement" json:"logged_list_id"`
	Token        string `gorm:"type:varchar(300)" json:"token"`
	UserEmail    string `gorm:"type:varchar(300)" json:"user_email"`
	UserId       int    `json:"user_id"`
}

type Products struct {
	ProductId    int    `gorm:"primaryKey;autoIncrement" json:"product_id"`
	Name         string `json:"name"`
	Participants string `json:"participants"`
	Quantity     string `json:"quantity"`
	Price        string `json:"price"`
	Purchase     string `json:"purchase"`
	GroupMember  string `json:"group_member"`
}

type Purchases struct {
	PurchasesId int    `gorm:"primaryKey;autoIncrement" json:"purchases_id"`
	Name        string `json:"name"`
	Timestamp   string `json:"timestamp"`
	Active      bool   `json:"active"`
	GroupMember string `json:"group_member"`
}
