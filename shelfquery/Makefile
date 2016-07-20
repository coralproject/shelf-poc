all: compile docker push

compile:
		CGO_ENABLED=0 go build -o relquery

docker:
		sudo docker build -t dwhitena/relquery .

push:
		sudo docker push dwhitena/relquery
