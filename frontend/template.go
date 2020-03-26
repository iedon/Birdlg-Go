package main

import (
	"net"
	"net/http"
	"strconv"
	"strings"
)

// Helper to check if the IP is valid
func isIP(s string) bool {
	return nil != net.ParseIP(s)
}

// Helper to check if the number is valid
func isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return nil == err
}

// Print HTML header to the given http response
func templateHeader(w http.ResponseWriter, r *http.Request, title string) {
	path := r.URL.Path[1:]
	split := strings.Split(path, "/")

	// Mark if the URL is for a whois query
	var isWhois bool = split[0] == "whois"
	var whoisTarget string = strings.Join(split[1:], "/")

	// Use a default URL if the request URL is too short
	// The URL is for return to IPv4 summary page
	if len(split) < 3 {
		path = "ipv4/summary/" + strings.Join(settingServers[:], "+") + "/"
	} else if len(split) == 3 {
		path += "/"
	}

	split = strings.Split(path, "/")

	// Compose URLs for link in navbar
	ipv4URL := "/" + strings.Join([]string{"ipv4", split[1], split[2], strings.Join(split[3:], "/")}, "/")
	ipv6URL := "/" + strings.Join([]string{"ipv6", split[1], split[2], strings.Join(split[3:], "/")}, "/")
	allURL := "/" + strings.Join([]string{split[0], split[1], strings.Join(settingServers[:], "+"), strings.Join(split[3:], "/")}, "/")

	// Check if the "All Server" link should be marked as active
	var serverAllActive bool = strings.ToLower(split[2]) == strings.ToLower(strings.Join(settingServers[:], "+"))

	// Print the IPv4, IPv6, All Servers link in navbar
	var serverNavigation string = `
		<li class="nav-item"><a class="nav-link` + (map[bool]string{true: " active"})[strings.ToLower(split[0]) == "ipv4"] + `" href="` + ipv4URL + `"` + `> IPv4 </a></li>
		<li class="nav-item"><a class="nav-link` + (map[bool]string{true: " active"})[strings.ToLower(split[0]) == "ipv6"] + `" href="` + ipv6URL + `"` + `> IPv6 </a></li>
        <span class="navbar-text">&nbsp;&nbsp;|&nbsp;&nbsp;</span>
        <li class="nav-item">
			<a class="nav-link` + (map[bool]string{true: " active"})[serverAllActive] + `" href="` + allURL + `"> ALL Nodes </a>
        </li>
    `

	// Add a link for each of the servers
	for _, server := range settingServers {
		var serverActive string
		if split[2] == server {
			serverActive = " active"
		}
		serverURL := "/" + strings.Join([]string{split[0], split[1], server, strings.Join(split[3:], "/")}, "/")

		serverNavigation += `
			<li class="nav-item">
				<a class="nav-link` + serverActive + `" href="` + serverURL + `">` + server + `</a>
            </li>`
	}

	// Add the options in navbar form, and check if they are active
	var optionKeys = []string{
		"summary",
		"detail",
		"route",
		"route_all",
		"route_bgpmap",
		"route_where",
		"route_where_all",
		"route_where_bgpmap",
		"whois",
		"traceroute",
	}
	var optionDisplays = []string{
		"show protocol",
		"show protocol all",
		"show route for ...",
		"show route for ... all",
		"show route for ... (bgpmap)",
		"show route where net ~ [ ... ]",
		"show route where net ~ [ ... ] all",
		"show route where net ~ [ ... ] (bgpmap)",
		"whois ...",
		"traceroute ...",
	}

	var options string
	for optionKeyID, optionKey := range optionKeys {
		options += "<option value=\"" + optionKey + "\""
		if (optionKey == "whois" && isWhois) || optionKey == split[1] {
			options += " selected"
		}
		options += ">" + optionDisplays[optionKeyID] + "</option>"
	}

	var target string
	if isWhois {
		// This is a whois request, use original path URL instead of the modified one
		// and extract the target
		target = whoisTarget
	} else if len(split) >= 4 {
		// This is a normal request, just extract the target
		target = strings.Join(split[3:], "/")
	}

	w.Write([]byte(`
<!DOCTYPE html>
<html>
    <head>
        <meta charset="UTF-8" />
        <meta name="viewport" content="width=device-width,initial-scale=1,shrink-to-fit=no"/>
        <title>` + title + `</title>
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@4.2.1/dist/css/bootstrap.min.css" rel="stylesheet">
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@4.4.1/dist/css/bootstrap.min.css" rel="stylesheet">
		<script src="https://cdn.jsdelivr.net/npm/viz.js@2.1.2/viz.min.js" crossorigin="anonymous"></script>
		<script src="https://cdn.jsdelivr.net/npm/viz.js@2.1.2/lite.render.js" crossorigin="anonymous"></script>
        <style>
			<!--
				html,body{height:100%}
                .container{margin-top:80px;min-height:100%}
                .table td{padding:0.5rem}
				.table td pre{margin-top:0;margin-bottom:0;font-size:100%}
				.table th:nth-child(2),.table th:nth-child(3),.table th:nth-child(4),.table th:nth-child(5),.table td:nth-child(2),.table td:nth-child(3),.table td:nth-child(4),.table td:nth-child(5){text-align:center}
                .badge{font-size:100%}
				.bd-footer{margin-top:30px;font-size:.875rem;text-align:center;background-color:#f7f7f7;width:100%}
                .bd-footer-links{padding-left:0;margin-bottom:1rem}
                .bd-footer p{margin-bottom:0}
                .bd-footer-links li{display:inline-block}
                .bd-footer-links li+li{margin-left:1rem}
                .bd-footer a{font-weight:600;color:#495057}
				.bd-footer a:focus,.bd-footer a:hover{color:#007bff}
				.nav-link.active{font-weight:bold}
                @media(min-width:576px){.bd-footer{text-align:left}}
            -->
        </style>
    </head>
    <body>

    <nav class="navbar navbar-expand-lg navbar-dark bg-secondary fixed-top">
        <a class="navbar-brand" href="/">` + NETWORK_NAME + `</a>
        <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
            <span class="navbar-toggler-icon"></span>
        </button>

        <div class="collapse navbar-collapse" id="navbarSupportedContent">
            <ul class="navbar-nav mr-auto">` + serverNavigation + `</ul>
            <form class="form-inline" action="/redir" method="GET">
                <div class="input-group">
                    <select name="action" class="form-control">` + options + `</select>
                    <input name="proto" class="d-none" value="` + split[0] + `">
                    <input name="server" class="d-none" value="` + split[2] + `">
                    <input name="target" class="form-control" placeholder="Target" aria-label="Target" value="` + target + `">
                    <div class="input-group-append">
                        <button class="btn btn-warning" type="submit">&raquo;</button>
                    </div>
                </div>
            </form>
        </div>
    </nav>

    <div class="container">
    `))
}

// Print HTML footer to http response
func templateFooter(w http.ResponseWriter) {
	w.Write([]byte(`
		</div>
		<footer class="footer bd-footer text-muted">
			<div class="container-fluid p-3 p-md-5">
				<ul class="bd-footer-links">
					<li><a href="` + LG_HOME_PAGE_URL + `">Looking Glass</a></li>
					<li><a href="` + HOME_PAGE_URL + `" target="_blank">` + NETWORK_NAME + `</a></li>
				</ul>
				` + FOOTER_HTML + `
			</div>
		</footer> 
    </body>
</html>
    `))
}

// Write the given text to http response, and add whois links for
// ASNs and IP addresses
func smartWriter(w http.ResponseWriter, s string) {
	w.Write([]byte("<pre>"))
	for _, line := range strings.Split(s, "\n") {
		var isASes bool = false
		var lineFormatted string
		words := strings.Split(line, " ")

		for wordID, word := range words {
			if len(word) == 0 {
				continue
			}
			if wordID > 0 && (len(words[wordID-1]) == 0 || words[wordID-1][len(words[wordID-1])-1] == ':') {
				// Insert TAB if there are multiple spaces before this word
				lineFormatted += "\t"
			} else {
				lineFormatted += " "
			}

			if isIP(word) {
				// Add whois link to the IP, handles IPv4 and IPv6
				lineFormatted += "<a href=\"/whois/" + word + "\">" + word + "</a>"
			} else if len(strings.Split(word, "%")) == 2 && isIP(strings.Split(word, "%")[0]) {
				// IPv6 link-local with interface name, like fd00::1%eth0
				// Add whois link to address part
				lineFormatted += "<a href=\"/whois/" + strings.Split(word, "%")[0] + "\">" + strings.Split(word, "%")[0] + "</a>"
				lineFormatted += "%" + strings.Split(word, "%")[1]
			} else if len(strings.Split(word, "/")) == 2 && isIP(strings.Split(word, "/")[0]) {
				// IP with a CIDR range, like 192.168.0.1/24
				// Add whois link to first part
				lineFormatted += "<a href=\"/whois/" + strings.Split(word, "/")[0] + "\">" + strings.Split(word, "/")[0] + "</a>"
				lineFormatted += "/" + strings.Split(word, "/")[1]
			} else if word == "AS:" || word == "\tBGP.as_path:" {
				// Bird will output ASNs later
				isASes = true
				lineFormatted += word
			} else if isASes && isNumber(word) {
				// Bird is outputing ASNs, ass whois for them
				lineFormatted += "<a href=\"/whois/AS" + word + "\">" + word + "</a>"
			} else {
				// Just an ordinary word, print it and done
				lineFormatted += word
			}
		}
		lineFormatted += "\n"
		w.Write([]byte(lineFormatted))
	}
	w.Write([]byte("</pre>"))
}

// Output a table for the summary page
func summaryTable(w http.ResponseWriter, isIPv6 bool, data string, serverName string) {
	w.Write([]byte("<table class=\"table table-striped table-bordered table-sm\">"))
	for lineID, line := range strings.Split(data, "\n") {
		var row [6]string
		var rowIndex int = 0

		words := strings.Split(line, " ")
		for wordID, word := range words {
			if len(word) == 0 {
				continue
			}
			if rowIndex < 4 {
				row[rowIndex] += word
				rowIndex++
			} else if len(words[wordID-1]) == 0 && rowIndex < len(row)-1 {
				if len(row[rowIndex]) > 0 {
					rowIndex++
				}
				row[rowIndex] += word
			} else {
				row[rowIndex] += " " + word
			}
		}

		// Ignore empty lines
		if len(row[0]) == 0 {
			continue
		}

		if lineID == 0 {
			// Draw the table head
			w.Write([]byte("<thead>"))
			for i := 0; i < 6; i++ {
				w.Write([]byte("<th scope=\"col\">" + row[i] + "</th>"))
			}
			w.Write([]byte("</thead><tbody>"))
		} else {

			w.Write([]byte("<tr>"))

			// Add link to detail for first column
			if isIPv6 {
				w.Write([]byte("<td><a href=\"/ipv6/detail/" + serverName + "/" + row[0] + "\">" + row[0] + "</a></td>"))
			} else {
				w.Write([]byte("<td><a href=\"/ipv4/detail/" + serverName + "/" + row[0] + "\">" + row[0] + "</a></td>"))
			}

			// Draw the other cells
			for i := 1; i < 6; i++ {
				if i == 3 {
					if strings.ToLower(row[i]) == "up" {
						w.Write([]byte("<td><span class=\"badge badge-success\">" + row[i] + "</span></td>"))
					} else if strings.ToLower(row[i]) == "start" && strings.ToLower(row[5]) == "passive" {
						w.Write([]byte("<td><span class=\"badge badge-info\">" + row[i] + "</span></td>"))
					} else if strings.ToLower(row[i]) == "down" {
						w.Write([]byte("<td><span class=\"badge badge-secondary\">" + row[i] + "</span></td>"))
					} else {
						w.Write([]byte("<td><span class=\"badge badge-danger\">" + row[i] + "</span></td>"))
					}
				} else if i == 5 {
					w.Write([]byte("<td><pre>" + row[i] + "</pre></td>"))
				} else {
					w.Write([]byte("<td>" + row[i] + "</td>"))
				}
			}
			w.Write([]byte("</tr>"))
		}
	}
	w.Write([]byte("</tbody></table>"))
}
