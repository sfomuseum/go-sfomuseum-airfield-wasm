# go-sfomuseum-airfield

Go package for working with airfield-related activities at SFO Museum (airlines, aircraft, airports).

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/sfomuseum/go-sfomuseum-airfield.svg)](https://pkg.go.dev/github.com/sfomuseum/go-sfomuseum-airfield)

Documentation is incomplete at this time.

## A note about "lookups"

As of this writing (October, 2021) most of the code in this package is focused around the code for "looking up" a record using one or more identifiers (for example an `ICAO` code or a SFO Museum database primary key) and/or to find the "current" instance of a given identifier (it turns out that IATA codes were re-used in the past). In all cases we're trying to resolve a given identifier back to a specific Who's On First ID for that "thing" in time and the code has been broken out in the topic and source specific subpackages.

All of those subpackages implement a generic `Lookup` interface which looks like this:

```
type Lookup interface {
	Find(context.Context, string) ([]interface{}, error)
	Append(context.Context, interface{}) error
}
```

There is an equivalent interface in the [go-sfomuseum-architecture](https://github.com/sfomuseum/go-sfomuseum-architecture) package which has all the same concerns as this package but is specific to the architectural elements at SFO. In time it would be best if both packages (`go-sfomuseum-architecture` and `go-sfomuseum-airfield`) shared a common "lookup" package/interface but that hasn't happened yet. There were past efforts around this idea in the [go-lookup*](https://github.com/search?q=org%3Asfomuseum+go-lookup) packages which have now been deprecated. If there is going to be a single common "lookup" interface package then the it will happen in a `github.com/sfomuseum/go-lookup/v2` namespace.

Which means that there is _a lot of duplicate code_ to implement functionality around the basic model in order to accomodate the different interfaces because an `icao.Aircraft` is not the same as a `sfomuseum.Aircraft` and a `sfomuseum.Aircraft` won't be the same as a `sfomuseum.Airport`. Things in the `sfomuseum` namespace tend to be more alike than not but there are always going to be edge-cases so the decision has been to suffer duplicate code (in different subpackages) rather than trying to shoehorn different classes of "things" in to a single data structure.

One alternative approach would be to adopt the [GoCloud `As` model](https://gocloud.dev/concepts/as/) but there is a sufficient level of indirection that I haven't completely wrapped my head around so it's still just an idea for now.

It's not great. It's just what we're doing today. The goal right now is to expect a certain amount of "rinse and repeat" in the short term while aiming to make each cycle shorter than the last.

### Data storage for "lookups"

The default data storage layer for lookup is an in-memory `sync.Map`. This works well for most cases and enforces a degree of moderation around the size of lookup tables. Another approach would be to use the [philippgille/gokv](https://github.com/philippgille/gokv) package (or equivalent) which is a simple interface with multiple storage backends. TBD..

## Airlines

### Adding a new airport to `sfomuseum/sfomuseum-data-enterprise`

Documentation to follow. Please consult the code for [cmd/create-airline](cmd/create-airline/main.go) in the meantime.

## Airports

### Adding a new airport to `sfomuseum/sfomuseum-data-whosonfirst`

Airport data is sourced frome `whosonfirst-data` repositories. The first thing to do is figure out which respository a given airport record is stored in. You can use the Who's On First Spelunker to look up this data:

* https://splelunker.whosonfirst.org/

_Note that some airports are stored in the `whosonfirst-data-admin-xy` repository because at the time of the "great splitting of the `whosonfirst-data-admin` repository in to per-country repositories" that airport's country of origin wasn't able to be determined. While importing airports that are in the `-admin-xy` repository is a good opportunity to move the record in to the correct repository that is not strictly necessary._ 

Although data is sourced from the Who's On First (WOF) project SFO Museum maintains local copies of WOF records in the `sfomuseum/sfomuseum-data-whosonfirst` repository. The `go-sfomuseum-whosonfirst` package was written to provide tools for importing data from WOF in to the `sfomuseum-data-whosonfirst` repository. For example:

```
$> cd /usr/local/sfomuseum/go-sfomuseum-whosonfirst
$> ./bin/import-feature \
	-reader-uri 'github://whosonfirst-data/whosonfirst-data-admin-{COUNTRY}' \
	{ID} {ID} {ID}
```

_Note also this example assumes that there is a checkout for the `sfomuseum-data-whosonfirst` repository in `/usr/local/data`. Consult the package documentation for details._

The `import-feature` tool will do a few things:

* It will retrieve the record for each (WOF) ID specified, as well any relevant ancestors for that ID (region, country)
* For each ID fetched it will create a corresonding JSON file in the `sfomuseum-data-whosonfirst/properties` folder. These JSON files contain any additional SFO Museum -specific properties or property values that should be overwritten (for example `wof:repo`).

Once the data has been imported make sure to commit your changes:

```
$> cd /usr/local/data/sfomuseum/sfomuseum-data-whosonfirst
$> git add {NEW FILES}
$> git commit -m "Add ..." {NEW FILES}
$> git push origin main
```

Commiting the changes is relevant to the `go-sfomuseum-airfield` package which provides pre-compiled lookup tables for things related to the SFO airfield (airlines, aircraft, airports). By default these tables are built by fetching the `sfomuseum/sfomuseum-data-whosonfirst` repository from GitHub. For example:

```
$> cd /usr/local/sfomuseum/go-sfomuseum-airfield
$> make compile
$> git commit -m "recompile data" .
$> git push origin main
```

Commiting the changes (to the `go-sfomuseum-airfield` package) is also relevant because a lot of other tools that use those lookup tables build them on the fly by fetching the serialized tables over the wire from GitHub; this allows us to update airfield data without involving the time-consuming process of updating every other package that uses `go-sfomuseum-airfield`.

### Adding a new airport to `whosonfirst/whosonfirst-data-admin-*`

Sometimes (not often) there are new airports which haven't been added the Who's On First (WOF) project yet. This is an example of how you might create a basic record for such a record, in this case [Istanbul Airport](https://spelunker.whosonfirst.org/id/1779770747/) in Turkey. The first step is to clone the `whosonfirst-data-admin-tr` repository:

```
$> git clone \
	--depth 1 \
	git@github.com:whosonfirst-data/whosonfirst-data-admin-tr.git \
	/usr/local/data/whosonfirst-data-admin-tr
```

The next step is to build a SQLite database, with the relevant spatial tables, that we can use to perform "point-in-polygon" operations to determine the new airport's parent and ancestors. Use the tools in the `go-whosonfirst-sqlite-features-index` to create this database:

```
$> cd /usr/local/whosonfirst/go-whosonfirst-sqlite-features-index

$> ./bin/sqlite-index-features \
	-all \
	-timings \
	-dsn /usr/local/data/whosonfirst-data-admin-tr.db \
	/usr/local/data/whosonfirst-data-admin-tr
```

For Turkey it takes about 3-4 minutes to create the `/usr/local/data/whosonfirst-data-admin-tr.db` database. Once the database has been created reference it when invoking the `wof-create` tool in the `whosonfirst/go-whosonfirst-exportify` package. For example:

```
$> cd /usr/local/whosonfirst/go-whosonfirst-exportify

$> ./bin/wof-create \
	-writer-uri repo:///usr/local/data/whosonfirst-data-admin-tr \
	-resolve-hierarchy \
	-spatial-database-uri 'sqlite://?dsn=/usr/local/data/whosonfirst-data-admin-tr.db' \
	-geometry '{"type":"Point","coordinates":[28.727778,41.262222]}' \
	-string-property 'properties.wof:placetype=campus' \
	-string-property 'properties.wof:country=TR' \
	-string-property 'properties.wof:name=Istanbul Airport' \
	-string-property 'properties.wof:repo=whosonfirst-data-admin-tr' \
	-int-property 'properties.mz:is_current=1' \
	-string-property 'properties.edtf:inception=2018-10-29' \
	-string-property 'properties.edtf:cessation=..' \
	-string-property 'properties.src:geom=wikipedia'
```

This will create a new, and minimal, record for the Istanbul Airport which can then be updated by hand as necessary. For testing and debugging purposes you can emit the new record to STDOUT but assigning the `-writer-uri` flag like this:

```
$> ./bin/wof-create -writer-uri stdout:// {OTHER OPTIONS}
```	

Commit the new record and then import it in to the `sfomuseum-data-whosonfirst` repository as described above. It's worth noting that I have write permissions on all the `*-data` repositories discussed so far. If you don't have write permissions all of these operations can also be accomplished using forks of the relevant repositories (which can then submit PRs upstream).

## See also

* https://github.com/sfomuseum-data/sfomuseum-data-aircraft
* https://github.com/sfomuseum-data/sfomuseum-data-enterprise
* https://github.com/sfomuseum-data/sfomuseum-data-whosonfirst
