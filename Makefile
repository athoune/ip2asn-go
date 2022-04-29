build: bin
	go build -o bin/ip2asn .

bin:
	mkdir -p bin

test:
	go test -cover \
		github.com/athoune/ip2asn-go/tsv

ip2asn-v4.tsv.gz:
	curl -O https://iptoasn.com/data/ip2asn-v4.tsv.gz

bench: ip2asn-v4.tsv.gz
	go test -benchmem -run=^$$ -bench . \
		github.com/athoune/ip2asn-go/tsv
