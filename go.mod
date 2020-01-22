module gitlab.com/ttpcodes/prismriver

go 1.11

require (
	github.com/adrg/libvlc-go v0.0.0-20191105210939-8fd26894baa1
	github.com/fsnotify/fsnotify v1.4.7
	github.com/gorilla/context v1.1.1
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/websocket v1.4.1
	github.com/hashicorp/hcl v1.0.0
	github.com/jinzhu/gorm v1.9.12
	github.com/jinzhu/inflection v1.0.0
	github.com/konsorten/go-windows-terminal-sequences v1.0.2
	github.com/lib/pq v1.3.0
	github.com/magiconair/properties v1.8.1
	github.com/mitchellh/mapstructure v1.1.2
	github.com/pelletier/go-toml v1.6.0
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/afero v1.2.2
	github.com/spf13/cast v1.3.1
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.1
	github.com/xfrr/goffmpeg v0.0.0-20191120110122-53b0a69281d4
	gitlab.com/ttpcodes/youtube-dl-go v0.0.0-20200122072320-c1e4ae6e5050
	golang.org/x/crypto v0.0.0-20200108215511-5d647ca15757
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553
	golang.org/x/sys v0.0.0-20200107162124-548cf772de50
	golang.org/x/text v0.3.2
	gopkg.in/ini.v1 v1.51.1 // indirect
	gopkg.in/yaml.v2 v2.2.7
)

replace github.com/xfrr/goffmpeg => github.com/ttpcodes/goffmpeg v0.0.0-20200110192638-089dcbcc69a5
