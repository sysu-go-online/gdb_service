package types

// DefaultLanguage defines the default language
var DefaultLanguage = "cpp"

// UserConf stores conf read from user file
type UserConf struct {
	Language    string
	Username    string
	ProjectName string
	Environment []string
}

// ResponseData contains output data and type
type ResponseData struct {
	Type string                 `json:"type"`
	Msg  map[string]interface{} `json:"msg"`
}

// SetDefault set default value
func (c *UserConf) SetDefault() {
	if c != nil {
		if c.Language == "" {
			c.Language = DefaultLanguage
		}
		if c.ProjectName == "" {
			c.ProjectName = "main"
		}
	}
}
