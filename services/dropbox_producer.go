package services
import (
	"github.com/stacktic/dropbox"
	"github.com/drborges/rivers/stream"
	"github.com/drborges/rivers/producers"
	"golang.org/x/net/context"
)

type DropboxBuilder struct {
	Context stream.Context
	CurrentCursor string
}

func (b *DropboxBuilder) DropboxProducer(db *dropbox.Dropbox) stream.Producer {
	return &producers.Observable{
		Context: b.Context,
		Capacity: 1000,
		Emit: func(w stream.Writable) {
			for {
				dp, err := db.Delta(b.CurrentCursor, "")
				if err != nil {
					panic(err)
				}

				select {
				case <-b.Context.Closed():
					return
				default:

					println("@@@@@@ Before for with entries %+v", dp)
					for _, deltaEntry := range dp.Entries {
						if deltaEntry.Entry == nil {
							//Handle deleted file later
						} else {
							w <- deltaEntry.Entry
						}
					}
				}


				b.CurrentCursor = dp.Cursor.Cursor

				println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
				println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
				println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
				println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
				println("Next iteration............")

				if dp.HasMore == false {
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


