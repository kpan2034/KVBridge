package config

import (
	. "KVBridge/types"
	"KVBridge/utils"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mitchellh/mapstructure"
	flag "github.com/spf13/pflag"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
)

const (
	defaultTimeoutMs = 1000 * time.Millisecond
)

// Wraps configuration of a KVNode
type Config struct {
	// Address at which node listens to client
	Address string `koanf:"address"`
	// Address at which node listens to other nodes
	Grpc_address string `koanf:"grpc_address"`
	// path to logfile
	LogPath string `koanf:"log_path"`
	// path to persistent storage file
	DataPath string `koanf:"data_path"`
	// other nodes in the cluster to connect to
	BootstrapServers []string `koanf:"bootstrap_servers"`
	// number of nodes each data point is replicated to
	ReplicationFactor int `koanf:"replication_factor"`
	// Read preference -- decided when the node boots
	ReadPreference OpPreference `koanf:"read_pref"`
	// Write preference -- decided when the node boots
	WritePreference OpPreference `koanf:"write_pref"`
	// How long to wait during request
	Timeout time.Duration `koanf:"timeout"`

	// For internal use
	ReadThreshold  int
	WriteThreshold int
}

func DefaultConfig() *Config {
	return &Config{
		Address:           ":6379",
		Grpc_address:      "localhost:50051",
		LogPath:           "./tmp/log",
		DataPath:          "./tmp/storage",
		BootstrapServers:  []string{"localhost:50051", "localhost:50052", "localhost:50053"},
		ReplicationFactor: 3,
		ReadPreference:    OpMajority,
		WritePreference:   OpMajority,
		Timeout:           defaultTimeoutMs,
		ReadThreshold:     2,
		WriteThreshold:    2,
	}
}

func NewConfigFromEnv() *Config {
	// Create new koanf instance
	k := koanf.New(".")

	// Set up POSIX-compliant flag libg
	f := flag.NewFlagSet("config", flag.ContinueOnError)
	f.Usage = func() {
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}

	// Path to one or more config files to load into koanf along with some config params.
	f.StringSlice("conf", []string{"./config.yaml"}, "path to one or more .yaml config files to load")
	f.String("address", ":6379", "address for server to listen for client")
	f.String("grpc_address", "localhost:50051", "address for server to listen for other nodes")
	f.String("log_path", "./tmp/log", "path to logfile")
	f.String("data_path", "./tmp/storage", "path to persistent storage")
	f.String("bootstrap_servers", "localhost:50051,localhost:50052", "bootstrap servers in the cluster")
	f.Int("replication_factor", 3, "number of nodes to replicate each data point")
	f.String("read_pref", "majority", "read preference for cluster")
	f.String("write_pref", "majority", "write preference for cluster")
	f.Int64("timeout", 10000000000, "default timeout between nodes in ns")
	f.Parse(os.Args[1:])

	// Load the config files provided in the commandline.
	cFiles, _ := f.GetStringSlice("conf")
	for _, c := range cFiles {
		if err := k.Load(file.Provider(c), yaml.Parser()); err != nil {
			log.Fatalf("could not load config file: %v", err)
		}
	}

	// Overwrite values in config file with ones provided on the command line
	if err := k.Load(posflag.Provider(f, ".", k), nil); err != nil {
		log.Fatalf("could not read command-line arguments: %v", err)
	}

	config := &Config{}
	decodeHook := mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToSliceHookFunc(","),
		utils.OpStrToPrefHookFunc(),
	)
	unmarshalConf := koanf.UnmarshalConf{
		DecoderConfig: &mapstructure.DecoderConfig{
			DecodeHook: decodeHook,
			Result:     config,
		},
	}
	err := k.UnmarshalWithConf("", config, unmarshalConf)
	if err != nil {
		log.Fatalf("could not unmarshall config: %v", err)
	}

	config.SetReadThreshold()
	config.SetWriteThreshold()

	return config
}

func (config *Config) GetReadThreshold() int {
	return config.ReadThreshold
}

func (config *Config) GetWriteThreshold() int {
	return config.WriteThreshold
}

func (config *Config) SetReadThreshold() {
	threshold := 0
	N := len(config.BootstrapServers)
	// There is a more optimal way of doing this -- but using bit manips are out of the purview of our deadline
	switch config.ReadPreference {
	case OpLocal:
		threshold = 1
	case OpMajority:
		threshold = (N / 2) + 1
	case OpAll:
		threshold = N
	default:
		log.Fatalf("invalid read preference specified: %v", config.ReadPreference)
	}
	config.ReadThreshold = threshold
}

func (config *Config) SetWriteThreshold() {
	threshold := 0
	N := len(config.BootstrapServers)
	switch config.WritePreference {
	case OpLocal:
		threshold = 1
	case OpMajority:
		threshold = (N / 2) + 1
	case OpAll:
		threshold = N
	default:
		log.Fatalf("invalid write preference specified: %v", config.ReadPreference)
	}
	config.WriteThreshold = threshold
}
