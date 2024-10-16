package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/calvindc/dpc-tox"
	"github.com/calvindc/dpc-tox/webtox/httpserve"
	"github.com/calvindc/dpc-tox/webtox/server/persistence"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
)

// the global tox instance
var tox *dpc_tox.Tox

// the global connection to the database
var storage *persistence.StorageConn

// the global options for HTTP authentication
var authOptions *httpserve.AuthOptions

type FileTransfer struct {
	fileHandle *os.File
	fileSize   uint64
	fileKind   dpc_tox.ToxFileKind
}

// Map of active file transfers
var transfers = make(map[uint32]FileTransfer)

func main() {
	var newToxInstance bool = false
	var options *dpc_tox.Options
	var databasePath string = filepath.Join(CFG_DATA_DIR, "userdata.db")

	var err error
	storage, err = persistence.Open(databasePath)
	if err != nil {
		log.Panic("DB initialisation failed.")
	}
	defer storage.Close()

	var toxSaveFilepath string
	flag.StringVar(&toxSaveFilepath, "p", filepath.Join(CFG_DATA_DIR, "webtox_save"), "path to save file")
	flag.Parse()
	fmt.Println("ToxData will be saved to", toxSaveFilepath)

	savedata, err := loadData(toxSaveFilepath)
	if err == nil {
		options = &dpc_tox.Options{
			IPv6Enabled:  true,
			UDPEnabled:   true,
			ProxyType:    dpc_tox.TOX_PROXY_TYPE_NONE,
			ProxyHost:    "127.0.0.1",
			ProxyPort:    5555,
			StartPort:    0,
			EndPort:      0,
			TcpPort:      CFG_TCP_PROXY_PORT,
			SaveDataType: dpc_tox.TOX_SAVEDATA_TYPE_TOX_SAVE,
			SaveData:     savedata}
	} else {
		options = &dpc_tox.Options{
			IPv6Enabled:  true,
			UDPEnabled:   true,
			ProxyType:    dpc_tox.TOX_PROXY_TYPE_NONE,
			ProxyHost:    "127.0.0.1",
			ProxyPort:    5555,
			StartPort:    0,
			EndPort:      0,
			TcpPort:      CFG_TCP_PROXY_PORT,
			SaveDataType: dpc_tox.TOX_SAVEDATA_TYPE_NONE,
			SaveData:     nil}
		newToxInstance = true
	}

	tox, err = dpc_tox.New(options)
	if err != nil {
		panic(err)
	}

	var toxid []byte

	toxid, err = tox.SelfGetAddress()
	if err != nil {
		panic(err)
	}
	fmt.Println("Tox ID:", strings.ToUpper(hex.EncodeToString(toxid)))

	if newToxInstance {
		fmt.Println("Setting username to default: WebTox User")
		tox.SelfSetName("WebTox User")
		tox.SelfSetStatusMessage("WebToxing around...")
		tox.SelfSetStatus(dpc_tox.TOX_USERSTATUS_NONE)
	} else {
		name, err := tox.SelfGetName()
		if err != nil {
			fmt.Println("Setting username to default: WebTox User")
			tox.SelfSetName("WebTox User")
		} else {
			fmt.Println("Username:", name)
		}

		if _, err = tox.SelfGetStatusMessage(); err != nil {
			if err = tox.SelfSetStatusMessage("WebToxing around..."); err != nil {
				panic(err)
			}
		}

		if _, err = tox.SelfGetStatus(); err != nil {
			if err = tox.SelfSetStatus(dpc_tox.TOX_USERSTATUS_NONE); err != nil {
				panic(err)
			}
		}
	}

	// Register our callbacks
	tox.CallbackFriendRequest(onFriendRequest)
	tox.CallbackFriendMessage(onFriendMessage)
	tox.CallbackFriendConnectionStatusChanges(onFriendConnectionStatusChanges)
	tox.CallbackFriendNameChanges(onFriendNameChanges)
	tox.CallbackFriendStatusMessageChanges(onFriendStatusMessageChanges)
	tox.CallbackFriendStatusChanges(onFriendStatusChanges)
	tox.CallbackFileRecv(onFileRecv)
	tox.CallbackFileRecvControl(onFileRecvControl)
	tox.CallbackFileRecvChunk(onFileRecvChunk)

	// Connect to the network
	// TODO add more servers (as fallback)
	pubkey, _ := hex.DecodeString("E20ABCF38CDBFFD7D04B29C956B33F7B27A3BB7AF0618101617B036E4AEA402D")
	err = tox.Bootstrap("3.0.24.15", 33445, pubkey)
	if err != nil {
		panic(err)
	}

	// Start the server
	go serveGUI()

	// Main loop
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ticker := time.NewTicker(25 * time.Millisecond)

	for {
		select {
		case <-c:
			fmt.Printf("\nSaving...\n")
			if err := saveData(tox, toxSaveFilepath); err != nil {
				fmt.Println(err)
			}

			fmt.Println("Killing")
			tox.Kill()
			return

		case <-ticker.C:
			tox.Iterate()
		}
	}
}

func serveGUI() {
	var err error
	var user, pass, salt string

	user, err = storage.GetKeyValue("settings_auth_user")
	if err == persistence.KeyNotFound {
		user, pass, salt = storeDefaultHTTPAuth()
	} else if err != nil {
		panic("GUI authentication username could not be determined.")
	}

	pass, err = storage.GetKeyValue("settings_auth_pass")
	if err == persistence.KeyNotFound {
		user, pass, salt = storeDefaultHTTPAuth()
	} else if err != nil {
		panic("GUI authentication password could not be determined.")
	}

	salt, err = storage.GetKeyValue("settings_auth_salt")
	if err == persistence.KeyNotFound {
		user, pass, salt = storeDefaultHTTPAuth()
	} else if err != nil {
		panic("GUI authentication salt could not be determined.")
	}

	authOptions = httpserve.NewAuthOptions(user, pass, salt)
	mux := http.NewServeMux()

	// paths that require authentication
	mux.Handle("/events", httpserve.BasicAuthHandler(handleWS, authOptions))
	mux.Handle("/api/", httpserve.BasicAuthHandler(handleAPI, authOptions))
	mux.Handle("/", httpserve.BasicAuthHandler(http.FileServer(http.Dir(CFG_HTML_DIR)), authOptions))

	// paths that *do not* require authentication
	mux.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir(CFG_IMG_DIR))))

	httpserve.CreateCertificateIfNotExist(CFG_DATA_DIR+CFG_CERT_PREFIX+"cert.pem", CFG_DATA_DIR+CFG_CERT_PREFIX+"key.pem", "localhost", 3072)
	err = httpserve.ListenAndUpgradeTLS(":8080", CFG_DATA_DIR+CFG_CERT_PREFIX+"cert.pem", CFG_DATA_DIR+CFG_CERT_PREFIX+"key.pem", mux)
}
