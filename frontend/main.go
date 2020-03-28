package main

import (
	"flag"
	"os"
	"strconv"
	"strings"
)

const NETWORK_NAME string = "iEdon-Net"
const PAGE_TITLE_PREFIX string = NETWORK_NAME + " Looking Glass: "
const LG_HOME_PAGE_URL string = "http://lg.iedon.dn42/"
const HOME_PAGE_URL string = "http://iedon.dn42/"
const FOOTER_HTML string = `
	<p>Bringing to you with ❤ by iEdon.</p>
	<p>Thanks for using iEdon-Net services.&nbsp;&nbsp;Have questions? <a href="mailto:iedonami@gmail.com">Contact me</a></p>
`

type settingType struct {
	servers     []string
	domain      string
	proxyPort   int
	whoisServer string
	listen      string
}

var setting settingType

func main() {
	var settingDefault = settingType{
		[]string{""}, "", 8000, "whois.verisign-grs.com", ":5000",
	}

	if serversEnv := os.Getenv("BIRDLG_SERVERS"); serversEnv != "" {
		settingDefault.servers = strings.Split(serversEnv, ",")
	}
	if domainEnv := os.Getenv("BIRDLG_DOMAIN"); domainEnv != "" {
		settingDefault.domain = domainEnv
	}
	if proxyPortEnv := os.Getenv("BIRDLG_PROXY_PORT"); proxyPortEnv != "" {
		var err error
		if settingDefault.proxyPort, err = strconv.Atoi(proxyPortEnv); err != nil {
			panic(err)
		}
	}
	if whoisEnv := os.Getenv("BIRDLG_WHOIS"); whoisEnv != "" {
		settingDefault.whoisServer = whoisEnv
	}
	if listenEnv := os.Getenv("BIRDLG_LISTEN"); listenEnv != "" {
		settingDefault.listen = listenEnv
	}

	serversPtr := flag.String("servers", strings.Join(settingDefault.servers, ","), "server name prefixes, separated by comma")
	domainPtr := flag.String("domain", settingDefault.domain, "server name domain suffixes")
	proxyPortPtr := flag.Int("proxy-port", settingDefault.proxyPort, "port bird-lgproxy is running on")
	whoisPtr := flag.String("whois", settingDefault.whoisServer, "whois server for queries")
	listenPtr := flag.String("listen", settingDefault.listen, "address bird-lg is listening on")
	flag.Parse()

	if *serversPtr == "" {
		flag.Usage()
		panic("no server set")
	} else if *domainPtr == "" {
		flag.Usage()
		panic("no base domain set")
	}

	setting = settingType{
		strings.Split(*serversPtr, ","),
		*domainPtr,
		*proxyPortPtr,
		*whoisPtr,
		*listenPtr,
	}

	webServerStart()
}
