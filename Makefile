.PHONY: all

gigo_path := $(GOPATH)/src/github.com/LyricalSecurity/gigo

all: setup update test
setup:
	@[[ ! -d $(gigo_path) ]] && git clone https://github.com/drborges/gigo $(gigo_path) || true
	go get github.com/LyricalSecurity/gigo/...
	go install github.com/LyricalSecurity/gigo

test:
	goapp test ./... -v -run=$(grep)

build:
	goapp build ./...

update:
	GIGO_GO=goapp gigo install -r requirements.txt

delete-branches:
	git branch | grep -v master | xargs -I {} git branch -D {}

serve:
	goapp serve --host 0.0.0.0 app/app.yaml

# Deployment tasks
deploy:
	goapp deploy -oauth app

rollback-deploy:
	appcfg.py --oauth2_refresh_token=***REMOVED*** rollback app

# FFMPEG related
compress-video:
	@which ffmpeg || brew install ffmpeg
	ffmpeg -i $(in) -b 512k -vcodec h264 -acodec copy $(out)

video-duration:
	@which ffmpeg || brew install ffmpeg
	ffmpeg -i $(in) 2>&1 | grep Duration | cut -d ' ' -f 4 | sed s/,// | cut -d ':' -f 3 | cut -d '.' -f 1 | xargs -I {} echo "{} seconds"

# Handy curls
http-get:
	curl -XGET https://gothere-dot-staging-api-getunseen.appspot.com$(path)

http-post:
	curl -XPOST https://gothere-dot-staging-api-getunseen.appspot.com$(path) \
		-H "Content-type: application/json" \
		-d '$(json)'

broadcast:
	curl -XPOST https://gothere-dot-staging-api-getunseen.appspot.com/places/$(place)/broadcasts \
		-H "Content-type: application/json" \
		-d '{"url": "$(url)", "length": $(length)}'


# ls ~/Downloads | grep gothere | xargs -I {} make compress-video in="{}" out=$(echo {} | sed -E 's/.*-([0-9]){2}\.mov/\1\.mov/g')
