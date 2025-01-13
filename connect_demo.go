package main

import (
	"fmt"
	"go/types"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type ProxyServ struct {
	name string
	ip   uint16 //100<1000
}
type Client struct {
	dIp        uint16 // dIp 1<30 динам..
	nameDomain string
	proxyIp    uint16
	id         int
}

const maxIp = 30

var HashIp = map[uint16]uint16{}

var DOMAIN = [10]ProxyServ{{
	name: "google.com",
	ip:   uint16(rand.Intn(1000)),
}, {
	name: "alibaba.com",
	ip:   uint16(rand.Intn(1000)),
}, {
	name: "git.com",
	ip:   uint16(rand.Intn(1000)),
}, {
	name: "yandex.com",
	ip:   uint16(rand.Intn(1000)),
},
	{
		name: "google.com",
		ip:   uint16(rand.Intn(1000)),
	},
	{
		name: "google.com",
		ip:   uint16(rand.Intn(1000)),
	}, {
		name: "alibaba.com",
		ip:   uint16(rand.Intn(1000)),
	}, {
		name: "git.com",
		ip:   uint16(rand.Intn(1000)),
	}, {
		name: "yandex.com",
		ip:   uint16(rand.Intn(1000)),
	},
	{
		name: "google.com",
		ip:   uint16(rand.Intn(1000)),
	},
}

func createP(id int) Client {

	return Client{dIp: uint16(rand.Intn(maxIp)), id: id, nameDomain: DOMAIN[rand.Intn(len(DOMAIN))].name}
}
func clientFlow(in chan<- Client) {
	count := 10
	for count > 0 {
		time.Sleep(time.Second * 2)
		in <- createP(count)
		count--
	}
	close(in)
}
func IpDns(in <-chan Client, out chan<- Client) {
	for v := range in {
		print("Новое подключение...  ")
		WriteLog(v)                        //новый клиент не сможет подключиться, пока другой не получит ip
		serv, err := getServ(v.nameDomain) //вариант решения (здесь задержка)/ каждому клиенту своя горутина->клиент ждет только своего ответа
		if err == nil {
			v.proxyIp = serv.ip
			print("/id:" + strconv.Itoa(v.id) + "->" + "Получает адрес прокси-сервера:")
			WriteLog(serv)
			out <- v
		} else {
			panic(err)
		}
	}
	println("Поток клиентов завершен.")
	close(out)
}

// CashingIp выполняет прокси-сервер
func CashingIp(out <-chan Client, outR chan<- Client) {
	for v := range out {
		time.Sleep(time.Second)
		if ip, ok := HashIp[v.dIp]; ok {
			v.dIp = ip
			print("Клиент был подключен...")
		} else {
			ip = ProxyIp(v.dIp)
			HashIp[v.dIp] = ip
			v.dIp = ip

		}
		outR <- v
	}
	close(outR)
}
func Connect(outR <-chan Client) {
	for v := range outR {
		print("Подключение к серверу завершено для:")
		WriteLog(v)
	}
}
func ProxyIp(ip uint16) uint16 {
	return ip>>rand.Intn(5) + 1
}
func getServ(domain string) (*ProxyServ, error) {
	time.Sleep(time.Second * 6) //иммитация задержки получения прокси-сервера(ip) по домену
	var randIp []ProxyServ
	for _, v := range DOMAIN {
		if v.name == domain {
			randIp = append(randIp, v)
		}
	}
	if len(randIp) != 0 {
		return &randIp[rand.Intn(len(randIp))], nil
	}
	return nil, types.Error{Msg: "Unknown domain!"}
}

func main() {
	in := make(chan Client)
	out := make(chan Client)
	outR := make(chan Client)
	go clientFlow(in)
	go IpDns(in, out)
	go CashingIp(out, outR)
	go Connect(outR)
	time.Sleep(time.Second * 1000000)
}

func (dns ProxyServ) String() string {
	return "domain:" + dns.name + "/ip DNS:" + strconv.Itoa((int)(dns.ip))
}
func (p Client) String() string {
	str := "ip:" + strconv.Itoa((int)(p.dIp)) + "/id:" + strconv.Itoa(p.id) + "/domain:" + p.nameDomain

	return str

}
func WriteLog(wr fmt.Stringer) {
	log.Println(wr.String())
}
