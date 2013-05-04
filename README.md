# Nothing to see here, move along

## But if you really like to know
Real-time massively multiplayer online space strategy arcade browser game!
A lot of words describing the first social game on Twitter that relies on total annihilation! :D

http://warcluster.com ! `>_< \m/`


## Installation

Warcluster is using [Redis](http://redis.io/), so you have to install it and run the server, using `redis-server` command.

Now it's time for the real deal. Install [Go](http://golang.org/) version 1.0 or higher. Decide which directory is going
to be your go workspace directory, set __$GOPATH__ to it and create `src/` directory inside. I'm going to use `~/go`.

    $ mkdir -p go/src
    $ export GOPATH=$HOME/go

Great. Now it's time to fetch the warcluster and install its requirements.

    $ git clone git@github.com:altras/WarCluster.git go/src/warcluster
    $ cd go/src/warcluster/
    $ go get

Just to be sure, everything is set up propery run the tests:

    $ go test ./...

If you see no errors, then __Hurray!__
