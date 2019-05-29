package mongodbatlas

//Config ...
type Config struct {
	Username string
	APIKey   string
}

//NewClient ...
func (c *Config) NewClient() interface{} {
	return nil
}
