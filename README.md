![warcluster logo](https://mir-s3-cdn-cf.behance.net/project_modules/disp/28168026279733.563586a13f4ae.gif)

WarCluster - Real-time Massive Multiplayer Online Space Strategy Arcade Browser Social Game for Twitter! ヘ(◕。◕ヘ)
===
:warning: This project's support is discontinued. Venture forth at your own peril. It's ~3 years old and we had no prior production experience with [Golang](http://golang.org/). It's considered as a playground although it's working (we held numerous tournaments with ~100 simultaneous real players and dozens small scale inhouse tests) ;)

Related links:

 - The javascript client repo: https://github.com/vladimiroff/warcluster-site
 - Amazingly strong NodeJS AI player with fuzzy logic (written at 2 day hackathon) here: https://github.com/lepovica/WarCluster-AI  ლ(ಠ益ಠლ)
 - The awesome UI designer @Denitsa that also made the art of the game: https://www.behance.net/gallery/26279733/War-Cluster
 - The team behind this project: [humans.txt](https://github.com/vladimiroff/warcluster-site/blob/develop/public/humans.txt)

What can you expect to see in this ~3-year old repo:
- Golang and WebSockets :D
- Golang and Redis =D
- Golang and tests ＼(＾O＾)／
- Voronoi diagram! ... wait, [wut](https://en.wikipedia.org/wiki/Voronoi_diagram)?

### History of how we got here

long story short: We wanted to learn something new ♪┏(・o･)┛♪┗ ( ･o･) ┓♪, experiment with the unthinkable & foresee the unimaginable. So we decided to make a web-fuckin-browser MMO as a side-project... It must had breathtaking tactics, large scale strategizing, thoughtful meta gameplay & dangerously social diplomacy. This is how WarCluster was born as the first social game for twitter! It was greatly inspired by [other games](https://hackpad.com/WarCluster-inspiration-1BvEVX758Ti)

#### [Setting](https://mir-s3-cdn-cf.behance.net/project_modules/1400/b2e9fc26279733.563593bb12d20.png)

It has 6 different races. Each race has it's own unique color.
![six races](https://mir-s3-cdn-cf.behance.net/project_modules/max_1200/0cbb1626279733.56353fde29024.png)

Each color corresponds to unique twitter #hashtag used for race wide communication (eg. #WarClusterRed). There're only verbal alliences between players.

#### Winning

The game is time based. Each round is approximately 10 days. The player with the *most controlled planets* (i.e. biggest empire) at the end of the round is proclaimed winner and therefore Galactic Emperor (here's one of the Beta winners [@valexiev](https://trello-attachments.s3.amazonaws.com/56e9cf6ad708c73bd6d0d26b/1352x623/0028f23f18dda84e9b5414f6c92e6c07/galactic_emperor.png))

#### Gameplay

Each player starts with a solar system consisting of 9 neutral planets. His home planet (the one with the asteroid belt) is impregnable. Typically each colonized planet has owner and generates for him army pilots per minute. The bigger the planet the more it generates.

At some point the army can grow bigger than the allowed population of the planet. So when the numbers get red bad things start to happen and army pilots are beginnig to die by the minute. Different sizes of planets have different population caps.
![planet info](https://trello-attachments.s3.amazonaws.com/56e9cf6ad708c73bd6d0d26b/789x435/d90a33aec2efcb3e2ceb7d62b3607faa/Screenshot-from-2015-04-19-17-17-41.png)![planet](https://trello-attachments.s3.amazonaws.com/56e9cf6ad708c73bd6d0d26b/556x418/e0134266d20b2618e5f56c3a09881c82/Screenshot-from-2015-04-19-17-14-04.png)

A player can send army pilots to *attack, spy or support* other planets depending on their verbal agreements with its owner. The spaceships travel slow so moves are carefully plotted. Imagine it as a real time game of [Go](https://en.wikipedia.org/wiki/Go_%28game%29)
![warcluster red armada](https://trello-attachments.s3.amazonaws.com/56e9cf6ad708c73bd6d0d26b/1221x604/91b0f6877364a81aa7368e2c1afbb7d3/WarClusterRed-armada.png):

 - Upon *attacking* a player conquers a planet if his attacking army is `> 1` than the other army. Everything is lost otherwise.
 - Upon *spying* a player gets the number of the army on the targeted planet.
 - Upon *supporting* a player is donating his army to the planet's owner.

A player should spend extra effort to pick his enemies wiser and his allies patiently. Plot twists are around every hour.

There're different kind of ships indicating the size of the army. From left (smallest army, no more than 500 pilots) to right (biggest army, more than 6000 pilots). This is used to understand what's the expected army power that's approaching you.

![xs ships](https://trello-attachments.s3.amazonaws.com/56e9cf6ad708c73bd6d0d26b/182x164/81889186e2ec41c6cb423bf737e554e8/Screenshot-from-2015-04-19-18-22-34.png)![sm ships](https://trello-attachments.s3.amazonaws.com/56e9cf6ad708c73bd6d0d26b/212x180/c7c9a5a1fecd8b1adcbcb22cf4f5afa3/Screenshot-from-2015-04-19-18-23-40.png)![md ships](https://trello-attachments.s3.amazonaws.com/56e9cf6ad708c73bd6d0d26b/253x218/1a7fbb8f69c4a3c7b1a709dd92da83fa/Screenshot-from-2015-04-19-18-24-11.png)![lg ships](https://trello-attachments.s3.amazonaws.com/56e9cf6ad708c73bd6d0d26b/181x202/7f6b8001328c942ee54ecfa1741b0ae9/Screenshot-from-2015-04-19-18-24-46.png)![xxl ships](https://trello-attachments.s3.amazonaws.com/56e9cf6ad708c73bd6d0d26b/252x222/d6e47bf777cef2ebed607c53f78c005f/Screenshot-from-2015-04-19-18-34-12.png)

### Running the server
prerequisites:
- [Redis](http://redis.io/), so you have to install it and run the server  (`redis-server` is the most lazy way)
- [Golang](http://golang.org/). Install version 1.3 or higher. Decide which directory is going to be your go workspace directory, set __$GOPATH__ to it and create `src/` directory inside. I'm going to use `~/go`.

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

Add your twitter consumer/secret keys [here](https://github.com/vladimiroff/WarCluster/blob/develop/config/config.gcfg.default#L12-L13) in order to have working social twitter login.

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

#### Contributing:

Fork it ( • ∀•)–Ψ and make required changes. After that push your changes in branch, which is named according to the changes you did. Initiate the PR (づ￣ ³￣)づ

