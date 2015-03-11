all:
	python rules.py --json > ./rules.json
	go build