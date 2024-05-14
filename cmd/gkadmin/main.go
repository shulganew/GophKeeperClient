package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/shulganew/GophKeeperClient/internal/app"
	"github.com/shulganew/GophKeeperClient/internal/client"
	"github.com/shulganew/GophKeeperClient/internal/client/oapi"
)

const MakeKey = "/admin/key"
const MakeMaster = "/admin/master"
const Shema = "https://"

func main() {

	log.Println("Start app")
	serverAddress := flag.String("a", "localhost:8443", "Service GKeeper address")
	_ = flag.Bool("n", false, "Create new ephemeral key (admin mode)")
	sertPath := flag.String("s", "cert/server.crt", "Service GKeeper address")
	login := flag.String("l", "admin", "Admin login")
	pw := flag.String("p", "123", "Admin pw")
	master := flag.String("m", "MasterKey:NewMasterKey", "Change master passwor old:new format. Don't use septarator in pass simbols! (admin mode)")
	flag.Parse()

	// Client with TLS session.
	c, err := oapi.NewClient(Shema+*serverAddress, oapi.WithHTTPClient(app.GetTLSClietn(*sertPath)))
	if err != nil {
		log.Fatal(err)
	}

	// Login admin
	_, jwt, status, err := client.UserLogin(c, *login, *pw)
	if status != http.StatusOK || err != nil {
		log.Println("Can't login.", err, status)
	}

	// Make request to create new eKey
	if isFlagPassed("n") {
		status, err := client.CreateNewEKey(c, jwt)
		if err != nil {
			log.Fatal(err)
		}
		if status == http.StatusCreated {
			log.Println("New Key created.")
		} else {
			log.Println("New ephemeral key not created. ")
		}

		return
	}

	// Change master key
	if isFlagPassed("m") {
		keys := strings.Split(*master, ":")
		status, err := client.CrateNewMaster(c, jwt, keys[0], keys[1])
		if err != nil {
			log.Fatal(err)
		}
		if status == http.StatusCreated {
			log.Println("New Key created.")
		} else {
			log.Println("Master key not set.")
		}
		return
	}
}

// check if flag set
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
