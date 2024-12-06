package config

type Cfg struct {
	MysqlInfo  MysqlCfg  `toml:"mysql"`
	BucketInfo BucketCfg `toml:"bucket"`
	FilterInfo FilterCfg `toml:"filter"`
}

type MysqlCfg struct {
	Host        string `toml:"host"`
	User        string `toml:"User"`
	Password    string `toml:"password"`
	DB          string `toml:"db"`
	Port        int64  `toml:"port"`
	MaxOpenConn int    `toml:"maxOpenConn"`
	MaxIdleConn int    `toml:"maxIdleConn"`
}

type BucketCfg struct {
	Region        string `toml:"region"`
	AccessKey     string `toml:"accessKey"`
	SecretKey     string `toml:"secretKey"`
	BucketName    string `toml:"bucketName"`
	OldBucketName string `toml:"oldBucketName"`
	FilePath      string `toml:"filePath"`
	ApiId         string `toml:"apiId"`
}

type FilterCfg struct {
	WorkNum     int    `toml:"workNum"`
	TableName   string `toml:"tableName"`
	ColumnName  string `toml:"columnName"`
	SplitFilter string `toml:"splitFilter"`
}
