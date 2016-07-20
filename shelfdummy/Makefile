all: compile docker push

compile:
		CGO_ENABLED=0 go build -o shelfdummy

docker:
		sudo docker build -t dwhitena/shelfdummy .

push:
		sudo docker push dwhitena/shelfdummy
