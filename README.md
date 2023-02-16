# falcondb

a in-memory document based database for fast prototyping

- `falconDB` is designed to be a JSON document-based database that relies on key based access to achieve `O(1)` access time. 
- Fields can hold references to other documents, which are automatically resolved up to a certain depth on retrieval. 
- All of these documents are stored as actual JSON files on the local machine, allowing developers to easily read, debug, and modify the data without the need for external tools. 

That means you can do stuff like
* make ID based authentication services
* store user data
* simple application cache
* and much much more, without the hassle of setting up an entire database schema and having to deal with drivers!

*However*, `falconDB` does not have any aggregation frameworks, advanced queries, sharding, or support for storage distribution. It was not created with the intention of ever being a production ready database, and should not be used as such!

#### key principles
* easy to lookup &mdash; key-based lookup in `O(1)` time
* easy to debug &mdash; all documents are JSON files which are human readable
* easy to deploy &mdash; single binary with no dependencies. no language specific drivers needed!

## endpoints
#### `GET /`
```bash
# get all files in database index
curl localhost:4000/

# example output on 200 OK
# > {"files":["test","test2","test3"]}
```

#### `POST /`
```bash
# manually regenerate index
# shouldn't need to be done as each operation should update index on its own
curl -X POST localhost:4000/

# example output on 200 OK
# > regenerated index
```

#### `GET /:key`
```bash
# get document with key `key`
curl localhost:4000/key

# example output on 200 OK (found key)
# > {"example_field": "example_value"}
# example output on 404 NotFound (key not found)
# > key 'key' not found
```

#### `PUT /:key`
```bash
# creates document `key` if it doesn't exist
# otherwise, replaces content of `key` with given
curl -X PUT -H "Content-Type: application/json" \
            -d '{"key1":"value"}' localhost:4000/key

# example output on 200 OK (create/update success)
# > create 'key' successful
```

#### `DELETE /:key`
```bash
# deletes document `key`
curl -X DELETE localhost:4000/key

# example output on 200 OK (delete success)
# > delete 'key' successful
# example output on 404 NotFound (key not found)
# > key 'key' doest not exist
```

#### `GET /:key/:field`
```bash
# get `example_field` of document `key`
curl localhost:4000/key/example_field

# example output on 200 OK (found field)
# > "example_value"
# example output on 400 BadRequest (field not found)
# > err key 'key' does not have field 'example_field'
# example output on 404 NotFound (key not found)
# > key 'key' not found
```
#### `PATCH /:key/:field`
```bash
# update `field` of document `key` with content
# if field doesnt exist, create it
curl -X PATCH -H "Content-Type: application/json" \
              -d '{"nested":"json!"}' \
              localhost:4000/key/example_field

# example output on 200 OK (found field)
# > patch field 'example_field' of key 'key' successful
# example output on 404 NotFound (key not found)
# > key 'key' not found
```

## commands
```bash
falcondb help  # shows a list of commands
falcondb start # start a falcondb server on :4000 using folder `db`
falcondb shell # start an interactive falcondb shell
```

#### `falcondb start`
This command starts a new `falcondb` server which listens for requests port `:4000` and uses the default folder `db`. The API endpoints are listed [here](#markdown-header-endpoints).

You can change the directory with the `--dir <value>, -d <value>` flag.
```bash
# e.g.
falcondb --dir some/folder start # start a falcondb server using folder `some/folder`
falcondb -d . start           # start a falcondb shell in the current directory
```

You can also change the port the server is hosted on with the `--p <value>, -p <value>` flag.
```bash
# e.g.
falcondb start --port 8081  # start a falcondb server on port 8081
falcondb -d . start -p 4000 # start a falcondb server on port 4000 using current directory
```

#### `falcondb shell`
This command starts a new `falcondb` interactive shell using the default folder `db`. The interactive shell isn't designed to do everything the API does, rather it is more like a quick tool to explore the database by allowing easy viewing of the database index, lookup of documents, and deletion of documents.

Similar to the `falcondb` server, you can change the directory with the `--dir <value>, -d <value>` flag.
```bash
# e.g.
falcondb -d . shell # start a falcondb shell using current directory
```

## reference resolution
You can refer to other documents by using a reference of the form `REF::<key>`. For example, with the following two JSONs:
#### `ref.json`
```json
{
  "key": "REF::nested"
}
```

#### `nested.json`
```json
{
  "nestedKey": "nestedVal"
}
```
You end up with the following:
```json
{
   "key": {
      "nestedKey": "nestedVal"
   }
}
```
This can be done within arrays and maps, to any arbitrary depth for which references should be resolved! The API has a default resolving depth of 3 while the CLI has a default of 0 but this can be explicitly changed if needed. For example through the API:
#### `GET /:key`
```bash
# get document with key `key` and only up to 1 layer of references resolved
curl localhost:4000/key?depth=1
```
#### `GET /:key/:field`
```bash
# get `example_field` of document `key`, resolving up to 5 layers deep
curl localhost:4000/key/example_field?depth=5
```
## running `falconDB`
#### from source
0. `git clone https://github.com/themillenniumfalcon/falconDB.git`
1. `go install github.com/themillenniumfalcon/falconDB`
2. `falcondb`

#### via docker
0. `docker pull nishankpr/falcondb:latest`
1. `docker run -p 4000:4000 nishankpr/falcondb:latest st` # change -p 4000:4000 to different port if necessary

**Note:** the docker version only supports the REST API server, not the CLI

## building `falconDB` from source
0. `git clone https://github.com/themillenniumfalcon/falconDB.git`
1. `make build`
2. (optional) for cross-platform builds, run `make build-all`