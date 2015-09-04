package services
import (
	"github.com/stacktic/dropbox"
	"github.com/drborges/rivers/stream"
	"github.com/drborges/rivers/producers"
	"golang.org/x/net/context"
	"fmt"
)


type DropboxDeltaProducerBuilder struct {
	Context       stream.Context
	CurrentCursor string
	Client        *dropbox.Dropbox
}

func (producer *DropboxDeltaProducerBuilder) Build() stream.Producer {
	return &producers.Observable{
		Context: producer.Context,
		Capacity: 1000,
		Emit: func(emitter stream.Emitter) {
			for {
				println(fmt.Sprint("Using Cursor value: %v", producer.CurrentCursor))
				page, err := producer.Client.Delta(producer.CurrentCursor, "")
				if err != nil {
					panic(err)
				}
				for _, deltaEntry := range page.Entries {
					if deltaEntry.Entry == nil {
						//Handle deleted file later
					} else {
						emitter.Emit(deltaEntry.Entry)
					}
				}

				producer.CurrentCursor = page.Cursor.Cursor

				println("--------Next iteration............")

				if page.HasMore == false {
					println("------ No more entries!!-----")
					return
				}


			}
		},
	}
}

const clientId = "m1ngt61jtc7wr7g"
const clientSecret = "34wbf1zkiu1xcnk"


func dropboxClient(context context.Context, accessToken string) *dropbox.Dropbox {
	db := dropbox.NewDropbox()
	db.SetContext(context)
	db.SetAppInfo(clientId, clientSecret)
	db.SetAccessToken(accessToken)

	return db
}


