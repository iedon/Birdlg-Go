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
	servers      []string
	domain       string
	proxyPort    int
	whoisServer  string
	listen       string
	dnsInterface string
}

var setting settingType

func main() {
	var settingDefault = settingType{
		[]string{""},
		"",
		8000,
		"whois.verisign-grs.com",
		":5000",
		"asn.cymru.com",
	}

	if env := os.Getenv("BIRDLG_SERVERS"); env != "" {
		settingDefault.servers = strings.Split(env, ",")
	}
	if env := os.Getenv("BIRDLG_DOMAIN"); env != "" {
		settingDefault.domain = env
	}
	if env := os.Getenv("BIRDLG_PROXY_PORT"); env != "" {
		var err error
		if settingDefault.proxyPort, err = strconv.Atoi(env); err != nil {
			panic(err)
		}
	}
	if env := os.Getenv("BIRDLG_WHOIS"); env != "" {
		settingDefault.whoisServer = env
	}
	if env := os.Getenv("BIRDLG_LISTEN"); env != "" {
		settingDefault.listen = env
	}
	if env := os.Getenv("BIRDLG_DNS_INTERFACE"); env != "" {
		settingDefault.dnsInterface = env
	}

	serversPtr := flag.String("servers", strings.Join(settingDefault.servers, ","), "server name prefixes, separated by comma")
	domainPtr := flag.String("domain", settingDefault.domain, "server name domain suffixes")
	proxyPortPtr := flag.Int("proxy-port", settingDefault.proxyPort, "port bird-lgproxy is running on")
	whoisPtr := flag.String("whois", settingDefault.whoisServer, "whois server for queries")
	listenPtr := flag.String("listen", settingDefault.listen, "address bird-lg is listening on")
	dnsInterfacePtr := flag.String("dns-interface", settingDefault.dnsInterface, "dns zone to query ASN information")
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
		*dnsInterfacePtr,
	}

	webServerStart()
}
