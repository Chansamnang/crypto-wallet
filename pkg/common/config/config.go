package config

var Config struct {
	App struct {
		Timezone string `yaml:"timezone"`
		Env      string `yaml:"env"`
		ProxyURL string `yaml:"proxyURL"`
		Port     int    `yaml:"port"`
	} `yaml:"app"`
	MysqlMaster struct {
		Address       *[]string `yaml:"address"`
		Username      *string   `yaml:"username"`
		Password      *string   `yaml:"password"`
		Database      *string   `yaml:"database"`
		MaxOpenConn   *int      `yaml:"maxOpenConn"`
		MaxIdleConn   *int      `yaml:"maxIdleConn"`
		MaxLifeTime   *int      `yaml:"maxLifeTime"`
		LogLevel      *int      `yaml:"logLevel"`
		SlowThreshold *int      `yaml:"slowThreshold"`
	} `yaml:"mysql_master"`
	MysqlSlave struct {
		Address       *[]string `yaml:"address"`
		Username      *string   `yaml:"username"`
		Password      *string   `yaml:"password"`
		Database      *string   `yaml:"database"`
		MaxOpenConn   *int      `yaml:"maxOpenConn"`
		MaxIdleConn   *int      `yaml:"maxIdleConn"`
		MaxLifeTime   *int      `yaml:"maxLifeTime"`
		LogLevel      *int      `yaml:"logLevel"`
		SlowThreshold *int      `yaml:"slowThreshold"`
	} `yaml:"mysql_slave"`
	Redis struct {
		ClusterMode bool     `yaml:"clusterMode"`
		Address     []string `yaml:"address"`
		Username    string   `yaml:"username"`
		Password    string   `yaml:"password"`
	} `yaml:"redis"`
	Log struct {
		LogLevel string `yaml:"logLevel"`
	} `yaml:"log"`
	Jwt struct {
		Key string `yaml:"key"`
	} `yaml:"jwt"`
	Blockchain struct {
		MnemonicPhrase   string `yaml:"mnemonicPhrase"`
		TronAlchemy      string `yaml:"tronAlchemy"`
		TronUSDTContract string `yaml:"tronUsdtContract"`
		TronGrpc         string `yaml:"tronGrpc"`
		TronGasFee       string `yaml:"tronGasFee"`
		EthAlchemy       string `yaml:"ethAlchemy"`
		EthUSDTContract  string `yaml:"ethUsdtContract"`
	} `yaml:"blockchain"`
}
