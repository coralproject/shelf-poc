all: compile docker push

compile:
		CGO_ENABLED=0 go build -o shelfrunner

docker:
		sudo docker build -t dwhitena/shelfrunner .

push:
		sudo docker push dwhitena/shelfrunner
