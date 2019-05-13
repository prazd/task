package main

import (
	. "github.com/prazd/task/server/handlers"
	. "github.com/prazd/task/server/lib"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/inv/up", PostOnly(InvSignUp))          // Sign up (inv)
	http.HandleFunc("/inv/in", PostOnly(InvSignIn))          // Sign in (inv), get inv info
	http.HandleFunc("/vol/up", PostOnly(VolSignUp))          // Sign up (vol)
	http.HandleFunc("/vol/in", PostOnly(VolSignIn))          // Sign in (vol), get vol info
	http.HandleFunc("/vol/ex", PostOnly(VolExit))            // Exit (vol)
	http.HandleFunc("/inv/ex", PostOnly(InvExit))            // Exit (inv)
	http.HandleFunc("/vol/geolist", GetOnly(VGeoList))       // Get geolocation and user info (vol)
	http.HandleFunc("/inv/geolist", GetOnly(IGeoList))       // Get geolocation and user info (inb)
	http.HandleFunc("/vol/getrev", PostOnly(GetVRev))        // Get reviews about vol
	http.HandleFunc("/vol/ch", PostOnly(VolHelp))            // (Vol) canHelp - set state(1) and geolocation
	http.HandleFunc("/inv/nh", PostOnly(InvHelp))            // (Inv) needHelp - set state(1) and geolocation
	http.HandleFunc("/vol/gp", PostOnly(VolGP))              // Set Vol geo
	http.HandleFunc("/inv/gp", PostOnly(InvGP))              // Set Inv geo
	http.HandleFunc("/vol/help", PostOnly(VolHelpInv))       // Set state(2) vol and in, get inv info
	http.HandleFunc("/vol/renouncement", PostOnly(VolRen))   // If vol can't help, but put conid
	http.HandleFunc("/inv/stophelp", PostOnly(IStop))        // Set vol and inv state(0)
	http.HandleFunc("/inv/helperinfo", PostOnly(HelperInfo)) // on inv side get geo and info about helper
	http.HandleFunc("/inv/volgeo", PostOnly(HelperGeo))      // get helper geo

	log.Fatal(http.ListenAndServe(":3000", nil))
}
