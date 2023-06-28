all: run


build:
	cd  ./agentFlow;  go build .
	

	
run:
	cd  ./manager;  go run .
	
