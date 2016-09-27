echo "((((Prototypical launcher for coral project services))))"
echo ""
echo "Repo: https://github.com/coralproject/shelf"
echo ""
echo "Requirements:"
echo "GOPATH must be set"
echo "GOBIN must be set"
echo "mongod must be installed"
echo ""

if [ ! -d "$GOPATH/src/github.com/coralproject/shelf" ]; then
	echo "Getting github.com/coralproject/shelf"
	go get github.com/coralproject/shelf
else 
	echo "Shelf found, no need to get it"
fi


mongod&
source $GOPATH/src/github.com/coralproject/shelf/config/localhost.cfg

echo "Building Corald"
cd $GOPATH/src/github.com/coralproject/shelf/cmd/sponge
go install
lsof -i tcp:16180 | awk 'NR>1 {print $2}' | xargs kill
corald&

echo "Building Sponge CLI"
cd $GOPATH/src/github.com/coralproject/shelf/cmd/sponge
go install

echo "Building Sponge Daemon"
cd $GOPATH/src/github.com/coralproject/shelf/cmd/sponged
go install 
lsof -i tcp:16181 | awk 'NR>1 {print $2}' | xargs kill
sponged&

echo "Building Wire CLI"
cd $GOPATH/src/github.com/coralproject/shelf/cmd/wire
go install 

echo "Building Xenia CLI"
cd $GOPATH/src/github.com/coralproject/shelf/cmd/xenia
go install 

echo "Building Xenia Daemon"
cd $GOPATH/src/github.com/coralproject/shelf/cmd/xeniad
go install 
lsof -i tcp:16182 | awk 'NR>1 {print $2}' | xargs kill
xeniad&

echo "Building Coral Daemon"
cd $GOPATH/src/github.com/coralproject/shelf/cmd/corald
go install 
