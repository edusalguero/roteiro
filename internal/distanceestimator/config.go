package distanceestimator

type Config struct {
	GoogleMaps GoogleMapsConf
}

type GoogleMapsConf struct {
	Enabled bool   `default:"false"`
	APIKey  string `required:"true"`
}
