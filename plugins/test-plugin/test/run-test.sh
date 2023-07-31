[ -f "passwordfile" ] && rm -f passwordfile 
[ -d "00/" ] && rm -rf 00/
[ -d "test00/" ] && rm -rf test00/ 
[ -d "01/" ] && rm -rf 01/
[ -d "test01/" ] && rm -rf test01/ 
[ -d "02/" ] && rm -rf 02/
[ -d "test02/" ] && rm -rf test02/

mkdir -p test00 test01 test02 00/keystore 01/keystore 02/keystore  00/geth 01/geth 02/geth  00/plugins 01/plugins 02/plugins 


cp ../engine.go test00/ 
cp ../engine.go ../main.go ../hooks.go ../tracer.go ../live_tracer.go  test01/
cp ../engine.go ../shutdown.go test02/
cd test00/ 
go build -buildmode=plugin -o ../00/plugins
cd ../
cd test01/ 
go build -buildmode=plugin -o ../01/plugins
cd ../
cd test02/ 
go build -buildmode=plugin -o ../02/plugins
cd ../

cp UTC--2021-03-02T16-47-49.510918858Z--f2c207111cb6ef761e439e56b25c7c99ac026a01 00/keystore
cp UTC--2021-03-02T16-47-39.492920333Z--4204477bf7fce868e761caaba991ffc607717dbf 01/keystore
cp UTC--2021-03-02T16-47-59.816632526Z--2cb2e3bdb066a83a7f1191eef1697da51793f631 02/keystore

cp nodekey00 00/geth/nodekey
cp nodekey01 01/geth/nodekey
cp nodekey02 02/geth/nodekey

echo -n "supersecretpassword" > passwordfile

$GETH init --datadir=./00 genesis.json
$GETH init --datadir=./01 genesis.json
$GETH init --datadir=./02 genesis.json

# miner node
$GETH --cache.preimages --config config00.toml --authrpc.port 8552 --port 64480 --verbosity=0 --nodiscover --networkid=6448 --datadir=./00/ --mine --miner.etherbase f2c207111cb6ef761e439e56b25c7c99ac026a01 --unlock f2c207111cb6ef761e439e56b25c7c99ac026a01 --http --http.api eth,debug,net --http.port 9545 --password passwordfile --allow-insecure-unlock &
pid0=$!

sleep 1
# passive node
$GETH --cache.preimages --config config01.toml --authrpc.port 8553 --port 64481 --verbosity=3 --syncmode=full --nodiscover --networkid=6448 --datadir=./01/ --unlock 4204477bf7fce868e761caaba991ffc607717dbf --miner.etherbase 4204477bf7fce868e761caaba991ffc607717dbf --password passwordfile --ws --ws.port 8546 --ws.api eth,admin --http --http.api eth,debug,net --http.port 9546 --allow-insecure-unlock &

sleep 1 

# shutdown node
$GETH --config config02.toml --authrpc.port 8556 --port 64484 --verbosity=0 --syncmode=full --nodiscover --networkid=6448 --datadir=./02/ --unlock 2cb2e3bdb066a83a7f1191eef1697da51793f631 --miner.etherbase 2cb2e3bdb066a83a7f1191eef1697da51793f631 --password passwordfile --ws --ws.port 8548 --ws.api eth,admin --http --http.api eth,debug,net --http.port 9547 --allow-insecure-unlock &
pid1=$!

sleep 5

if ps -p $pid1 > /dev/null; then
  kill $pid1
fi

sleep 255

if ps -p $pid0 > /dev/null; then
  kill $pid0
fi

wait

rm -f passwordfile
rm -rf 00/ 01/ 02/ test00/ test01/ test02/


