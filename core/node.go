package core

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/NetSepio/beacon/api"
	grpc "github.com/NetSepio/beacon/gRPC"
	"github.com/NetSepio/beacon/p2p"
	"github.com/NetSepio/beacon/util"
	"github.com/NetSepio/beacon/web3"
	"github.com/blocto/solana-go-sdk/pkg/hdwallet"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/sha3"
	helmet "github.com/danielkov/gin-helmet"
)

var wg sync.WaitGroup

// These variables will be set at build time
var (
	NodeID      string
	NodeName    string
	NodeSpec    string
	NodeConfig  string
	NodeAccess  string
	NodeRegion  string
	NodeIP      string
	NodeVersion string
)

// LoadNodeDetails loads the node details from environment variables
func LoadNodeDetails() {
	NodeID = os.Getenv("NODE_ID")
	NodeName = os.Getenv("NODE_NAME")
	NodeSpec = os.Getenv("NODE_SPEC")
	NodeConfig = os.Getenv("NODE_CONFIG")
	NodeAccess = os.Getenv("NODE_ACCESS")
	NodeRegion = os.Getenv("NODE_REGION")
	NodeIP = os.Getenv("NODE_IP")
	NodeVersion = os.Getenv("NODE_VERSION")

	if NodeID == "" {
		log.Fatal("NODE_ID environment variable is not set")
	}
	if NodeName == "" {
		log.Fatal("NODE_NAME environment variable is not set")
	}
	if NodeSpec == "" {
		log.Fatal("NODE_SPEC environment variable is not set")
	}
	if NodeConfig == "" {
		log.Fatal("NODE_CONFIG environment variable is not set")
	}
	if NodeAccess == "" {
		log.Fatal("NODE_ACCESS environment variable is not set")
	}
	if NodeRegion == "" {
		log.Fatal("NODE_REGION environment variable is not set")
	}
	if NodeIP == "" {
		log.Fatal("NODE_IP environment variable is not set")
	}
	if NodeVersion == "" {
		log.Fatal("NODE_VERSION environment variable is not set")
	}
}

func RungRPCServer() {
	grpc_server := grpc.Initialize()
	port := os.Getenv("GRPC_PORT")

	log.Printf("Starting gRPC Api, Listening on Port : %s", port)

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		wg.Done()
		log.Fatal("Unable to listen on port", port)
	}

	//Server GRPC
	if err := grpc_server.Serve(listener); err != nil {
		wg.Done()
		log.Fatal("Failed to create GRPC server!")
	}
	wg.Done()
}

// RunBeaconNode starts the Beacon node
func RunBeaconNode() {
	log.Printf("Starting NetSepio - Erebrus Version: %s", util.Version)

	// check directories or create it
	if !util.DirectoryExists(filepath.Join(os.Getenv("WG_CONF_DIR"))) {
		err := os.Mkdir(os.Getenv("WG_CONF_DIR"), 0755)
		if err != nil {
			log.Fatalf("failed to create wireguard configuration directory: %v", err)
		}
	}

	// check directories or create it
	if !util.DirectoryExists(filepath.Join(os.Getenv("WG_CLIENTS_DIR"))) {
		err := os.Mkdir(os.Getenv("WG_CLIENTS_DIR"), 0755)
		if err != nil {
			log.Fatalf("failed to create wireguard clients directory: %v", err)
		}
	}

	// check if server.json exists otherwise create it with default values
	if !util.FileExists(filepath.Join(os.Getenv("WG_CONF_DIR"), "server.json")) {
		_, err := ReadServer()
		if err != nil {
			log.Fatal("server.json does not exist and unable to open")
		}
	}

	if os.Getenv("RUNTYPE") == "debug" {
		// set gin release debug
		gin.SetMode(gin.DebugMode)
	} else {
		// set gin release mode
		gin.SetMode(gin.ReleaseMode)
		// disable console color
		gin.DisableConsoleColor()
		// log level info
		log.SetLevel(log.InfoLevel)
	}

	// dump wg config file
	err := UpdateServerConfigWg()
	util.CheckError("Error while creating WireGuard config file: ", err)

	LoadNodeDetails()

	// Register node on chain if configured
	if err := RegisterNodeOnChain(); err != nil {
		log.Printf("Failed to register node on %s: %v", os.Getenv("CHAIN_NAME"), err)
	}

	go p2p.Init()
	wg.Add(1)

	if os.Getenv("GRPC_PORT") != "" {
		wg.Add(1)
		go RungRPCServer()
	}

	if os.Getenv("HTTP_PORT") != "" {
		ginApp := gin.Default()
		config := cors.DefaultConfig()
		config.AllowOrigins = []string{os.Getenv("GATEWAY_DOMAIN")}
		ginApp.Use(cors.New(config))
		ginApp.Use(helmet.Default())

		ginApp.Use(func(ctx *gin.Context) {
			ctx.Set("cache", cache.New(60*time.Minute, 10*time.Minute))
			ctx.Next()
		})

		ginApp.Use(static.Serve("/", static.LocalFile("./webapp", false)))
		ginApp.NoRoute(func(c *gin.Context) {
			c.JSON(404, gin.H{"status": 404, "message": "Invalid Endpoint Request"})
		})

		api.ApplyRoutes(ginApp)
		err = ginApp.Run(fmt.Sprintf("%s:%s", os.Getenv("SERVER"), os.Getenv("HTTP_PORT")))
		util.CheckError("Failed to Start HTTP Server: ", err)
	}

	wg.Wait()
}

// GenerateEthereumWalletAddress generates an Ethereum wallet address from a mnemonic
func GenerateEthereumWalletAddress(mnemonic string) (string, *ecdsa.PrivateKey, error) {
	// Generate seed from mnemonic
	seed := bip39.NewSeed(mnemonic, "")

	// Generate master key
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate master key: %v", err)
	}

	// Derive Ethereum path
	path := []uint32{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0, 0}
	key := masterKey
	for _, i := range path {
		key, err = key.NewChildKey(i)
		if err != nil {
			return "", nil, fmt.Errorf("failed to derive key: %v", err)
		}
	}

	// Generate private key
	privateKey, err := crypto.ToECDSA(key.Key)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate private key: %v", err)
	}

	// Get public key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", nil, fmt.Errorf("failed to get public key")
	}

	// Get address
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return address, privateKey, nil
}

// toChecksumAddress converts an Ethereum address to checksum address
func toChecksumAddress(address string) string {
	// Remove 0x prefix if present
	addr := strings.ToLower(address)
	if strings.HasPrefix(addr, "0x") {
		addr = addr[2:]
	}

	// Calculate hash
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(addr))
	hashBytes := hash.Sum(nil)

	// Convert to checksum address
	result := make([]byte, len(addr))
	for i := 0; i < len(addr); i++ {
		if hashBytes[i/2]&(1<<(4*(1-i%2))) != 0 {
			result[i] = strings.ToUpper(string(addr[i]))[0]
		} else {
			result[i] = addr[i]
		}
	}

	return "0x" + string(result)
}

// GenerateWalletAddressSolanaAndEclipse generates a Solana wallet address from a mnemonic
func GenerateWalletAddressSolanaAndEclipse(mnemonic string) {
	// Generate seed from mnemonic
	seed := bip39.NewSeed(mnemonic, "")

	// Generate master key
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		log.Fatalf("Failed to generate master key: %v", err)
	}

	// Derive Solana path
	path := []uint32{0x80000000 + 44, 0x80000000 + 501, 0x80000000 + 0, 0, 0}
	key := masterKey
	for _, i := range path {
		key, err = key.NewChildKey(i)
		if err != nil {
			log.Fatalf("Failed to derive key: %v", err)
		}
	}

	// Generate keypair
	keypair, err := types.AccountFromSeed(key.Key)
	if err != nil {
		log.Fatalf("Failed to generate keypair: %v", err)
	}

	// Print address
	fmt.Printf("Solana Address: %s\n", keypair.PublicKey.ToBase58())
}

// GenerateWalletAddressSui generates a Sui wallet address from a mnemonic
func GenerateWalletAddressSui(mnemonic string) {
	// Generate seed from mnemonic
	seed := bip39.NewSeed(mnemonic, "")

	// Generate master key
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		log.Fatalf("Failed to generate master key: %v", err)
	}

	// Derive Sui path
	path := []uint32{0x80000000 + 44, 0x80000000 + 784, 0x80000000 + 0, 0, 0}
	key := masterKey
	for _, i := range path {
		key, err = key.NewChildKey(i)
		if err != nil {
			log.Fatalf("Failed to derive key: %v", err)
		}
	}

	// Generate keypair
	keypair, err := ed25519.GenerateKey(key.Key)
	if err != nil {
		log.Fatalf("Failed to generate keypair: %v", err)
	}

	// Print address
	fmt.Printf("Sui Address: %s\n", hex.EncodeToString(keypair.Public().(ed25519.PublicKey)))
}

// GenerateWalletAddressAptos generates an Aptos wallet address from a mnemonic
func GenerateWalletAddressAptos(mnemonic string) {
	// Generate seed from mnemonic
	seed := bip39.NewSeed(mnemonic, "")

	// Generate master key
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		log.Fatalf("Failed to generate master key: %v", err)
	}

	// Derive Aptos path
	path := []uint32{0x80000000 + 44, 0x80000000 + 637, 0x80000000 + 0, 0, 0}
	key := masterKey
	for _, i := range path {
		key, err = key.NewChildKey(i)
		if err != nil {
			log.Fatalf("Failed to derive key: %v", err)
		}
	}

	// Generate keypair
	keypair, err := ed25519.GenerateKey(key.Key)
	if err != nil {
		log.Fatalf("Failed to generate keypair: %v", err)
	}

	// Print address
	fmt.Printf("Aptos Address: %s\n", hex.EncodeToString(keypair.Public().(ed25519.PublicKey)))
}

// RegisterNodeOnChain registers the node on the blockchain
func RegisterNodeOnChain() error {
	return web3.RegisterNodeOnChain()
}

// DeactivateNode deactivates the node on the blockchain
func DeactivateNode() error {
	return web3.DeactivateNode()
}

// GetNodeStatus returns the current status of the node from the blockchain
func GetNodeStatus() (*web3.NodeStatus, error) {
	return web3.GetNodeStatus()
}

