package flags

//
//import (
//	"errors"
//	"log"
//	"strconv"
//	"strings"
//)
//
//type NetAddress struct {
//	Host string
//	Port int
//}
//
//func (a *NetAddress) String() string {
//	return a.Host + ":" + strconv.Itoa(a.Port)
//}
//
//func (a *NetAddress) Set(s string) error {
//	hp := strings.Split(s, ":")
//	if len(hp) != 2 {
//		return errors.New("need address in a form host:port")
//	}
//	port, err := strconv.Atoi(hp[1])
//	if err != nil {
//		log.Print(err)
//		return err
//	}
//	a.Host = hp[0]
//	a.Port = port
//	return nil
//}
