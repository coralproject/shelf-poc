all: compile docker push

compile:
		CGO_ENABLED=0 go build -o dummy

docker:
		sudo docker build -t dwhitena/dummy .

push:
		sudo docker push dwhitena/dummy
