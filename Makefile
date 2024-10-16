# 0) basic build pre deploy
build:
	@go build -o bin src/main.go

go_live_reload:
	@air

templ_live_reload:
	@air2 -c .air.templ.toml

# 1) watch .go files but no templ etc. then triggers templ generate by changing fake.templ file
air:
	@air -c .air.go_and_triger_templ.toml

# 2) temp gnerate to watch templ files and to perform hot reload in browser
devtg:
	@templ generate --watch --proxy="http://0.0.0.0:10001" --cmd="go run ./src/main.go"

# 2a) temp gnerate 2 plus stderr redirect /dev/pts/13
devtg2:
	@templ generate --watch --proxy="http://0.0.0.0:10001" --cmd="go run ./src/main.go" 2>/dev/pts/13
devtg3:
	@templ generate --watch --proxy="http://0.0.0.0:10001" --proxyport="7332" --cmd="go run ./src/main.go" 2>/dev/pts/11

templ:
	@templ generate --watch --proxy="http://0.0.0.0:10001" --proxyport="7332" --cmd="go run ./src/main.go" 2>/dev/pts/$(t)

# 3) tailwind to watch changes in views folder for html, templ, js
devtw:
	@tailwindcss -i assets/input.css -o assets/output.css --watch

# 2+3) NE RADI: combined templ and tailwind but tailwind is not working well with some classes
devtt:
	@templ generate --watch --proxy="http://0.0.0.0:10000" --cmd="go run ./src/main.go & tailwindcss -i input.css -o assets/output.css"
