all: compile docker push

compile:
		CGO_ENABLED=0 go build -o dummydirect

docker:
		sudo docker build -t dwhitena/dummydirect .

push:
		sudo docker push dwhitena/dummydirect
