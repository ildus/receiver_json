all:
	node rules.js > ./rules.json
	go build