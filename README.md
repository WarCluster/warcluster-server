# Nothing to see here, move along


## But if you really like to know

Real-time massively multiplayer online space strategy arcade browser game! A
lot of words describing the first social game on Twitter that relies on total
annihilation! :D

http://warcluster.com ! `>_< \m/`


## Installation

Warcluster is using [Redis](http://redis.io/), so you have to install it and
run the server  (`redis-server` is the most lazy way).

Now it's time for the real deal. Install [Go](http://golang.org/) version 1.3
or higher. Decide which directory is going to be your go workspace directory,
set __$GOPATH__ to it and create `src/` directory inside. I'm going to use
`~/go`.

    $ mkdir -p go/src
    $ export GOPATH=$HOME/go
    $ export PATH=$PATH:$GOPATH/bin

It's a good idea to place these two _export_ commands somewhere in your .bashrc
file, though. Otherwise you would have to execute them on each re-login in your
system. Now it's time to fetch warcluster and install its requirements.

    $ git clone git@github.com:Vladimiroff/WarCluster.git $GOPATH/src/warcluster
    $ cd $GOPATH/src/warcluster/
    $ go get -v github.com/stretchr/testify
    $ go get -v

If you run redis on your localhost without any custom configuration and you're
okay with the game running on port 7000, then you should be able to run the
server without any modifications. Otherwise, copy `config/config.gcfg.default`
to `config/config.gcfg` and make your changes.

Just to be sure, everything is set up propery run the tests:

    $ go test ./...

If you see no errors, then __Hurray!__ Let's actually compile server:

    $ go install

If this command runs without errors you should have a `bin/` directory in your
`$GOPATH` and now you can run it:

    $ warcluster
