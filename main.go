package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/miekg/dns"
	"google.golang.org/appengine"
)

func init() {
	http.HandleFunc("/", handler)
}

func nslookup(name, server string) (*dns.Msg, error) {
	dnsClient := new(dns.Client)
	msg := dns.Msg{
		MsgHdr: dns.MsgHdr{
			RecursionDesired: true,
		},
		Question: []dns.Question{{
			Name:   name,
			Qtype:  dns.TypeANY,
			Qclass: dns.ClassINET,
		}},
	}
	in, rtt, err := dnsClient.Exchange(&msg, server)
	if err != nil {
		log.Printf("Error was not nil: %v", err)
		return in, err
	}
	log.Printf("RTT: %v", rtt)
	if rcode, ok := dns.RcodeToString[in.Rcode]; ok {
		log.Printf("Rcode: %v", rcode)
	}
	if opcode, ok := dns.OpcodeToString[in.Opcode]; ok {
		log.Printf("Opcode: %v", opcode)
	}
	log.Printf("in: %+v", in)
	log.Printf("*in: %+v", *in)
	log.Printf(`error msg:"%+v"`, in)
	return in, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	server := r.URL.Query().Get("server")

	responseContext := struct {
		Body   string
		Name   string
		Server string
	}{}

	responseContext.Name = name
	responseContext.Server = server

	tmpl, err := template.ParseFiles("./index.html.template")
	if err != nil {
		log.Printf("parsing template failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	msg, err := nslookup(name, server)
	if err != nil {
		responseContext.Body = fmt.Sprintf("%s", err)
	} else {
		responseContext.Body = fmt.Sprintf("%s", msg)
	}
	err = tmpl.Execute(w, responseContext)
	if err != nil {
		log.Printf("%s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func main() {
	appengine.Main()
}
