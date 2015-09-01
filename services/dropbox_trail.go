package services
import (
	"github.com/drborges/rivers"
	"google.golang.org/appengine"
	"github.com/drborges/rivers/stream"
	"fmt"
	"net/http"
)

func DropboxDelta(req *http.Request) {

//	accessToken string, cursor string
//	accessToken := "4GFHgi3IL2IAAAAAAAAAClnkF5SyJKu7pcC64oNVg18r6tecpNkJdyY7spxDw9hF"
	accessToken := "Cdn2AKzQYPwAAAAAAAAT-9gnnBDOG68LO1Oms5qW9-4WLme9TX02EHLiQLgR2qoA"
	cursor := ""
	steamContext := rivers.NewContext()
	db := DropboxBuilder{steamContext}
	producer := db.DropboxProducer(appengine.NewContext(req), accessToken, cursor)


	err := rivers.NewWith(steamContext).From(producer).Each(func(data stream.T) { println(fmt.Sprintf("#### %+v", data))}).Drain()
	println(fmt.Sprintf("########## %v", err))
}