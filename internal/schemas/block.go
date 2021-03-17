package schemas

type BlockType string

const (
	BlockTypeScheme BlockType = "scheme"
	BlockTypeDomain BlockType = "domain"
	BlockTypeRegex  BlockType = "regex"
)

type Block struct {
	Type  BlockType `gorm:"type:enum('scheme', 'domain', 'regex');primaryKey"`
	Value string    `gorm:"type:varchar(255);primaryKey"`
}
