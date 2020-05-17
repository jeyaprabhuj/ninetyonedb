# NinetyOneDB

JSON Storage server inspired by couchdb rest api access written in Golang.


## Installation

Use the go get

```bash
git clone https://github.com/jeyaprabhuj/ninetyonedb.git

cd ninetyonedb

go build

eg: KVPORT=6060 ./ninetyonedb

```

## Usage

Start server by environment vairiable KVPORT=<port number>
```bash
./ninetyonedb
```


http://<ip>:<port>/http://localhost:8001/storages/

Initiate a HTTP Put with Request body JSON {"name":"<database/storage-name>"}

list all storages - http://<ip>:<port>/http://localhost:8001/storages/
  
Individual access - http://<ip>:<port>/http://localhost:8001/storages/<storage-name>

e.g. Consider database name is "Shop"

Individual access - http://<ip>:<port>/http://localhost:8001/storages/Shop

Client can be Postman ,any browser or client written in any language 

## Roadmap
Current status is reference implementation.
Equivalent implementation in Elixir will be also added as another repo.

First release:
Storage file compaction 

Second Release:
Support for selectors like couchdb

Third Release:
Simple replication and backup support

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
