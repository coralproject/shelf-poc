all: compile docker push

compile:
		CGO_ENABLED=0 go build -o shelfquery

docker:
		sudo docker build -t dwhitena/shelfquery .

push:
		sudo docker push dwhitena/shelfquery
