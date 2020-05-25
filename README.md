# xDoc

xDoc will traverse all of GitLab -- groups/subgroups and projects that it has access to and will
look in each repository for:

* a top-level file named `.xdoc.yaml` that allows for setting some scrape (and other) options.
* a directory named xdoc (see the `-dir` option)

It will both render and index these markdown files, the following endpoints are available:

* `/s`; exposed the bleve search endpoint.
* `/s/<query>`; directly search for `<query>`.
* `/s/?q=<query<`; GET endpoint for the search page, may be hit directly, by `/s/<query>` might be
  more convenient.
* `/r/<groupname>/<subgroupname>/<projectname>`; hitting this endpoint will render the `index.md` markdown file.
* `/r/<groupname>/<subgroupname>/<projectname>/<filename>`; hitting this endpoint will render the `<filename>.md` markdown file.
* `/r/<groupname>/<projectname>`; hitting this endpoint will render the `index.md` markdown file.
* `/r/<groupname>/<projectname>/<filename>`; hitting this endpoint will render the `<filename>.md` markdown file.

More endpoints may be added in the future; for instance listing all repositories or things like
that.

## .xdoc.yaml

This YAML file contains various options, such as the Git ref used to needs to be retrieved, the
markdown flavor used for parsing and in which language the pages are written. This file MUST exist
in the master branch of the repository.

~~~ yaml
lang: LANG
ref: GIT-REFERENCE
flavor: mmark|commonmark|gfm|codelab
~~~

## Trying It

Using the public gitlab.com instance you can run the following:

~~~ sh
./xdoc -group miekg -gitlab https://gitlab.com -dir xdoc
~~~

## Conventions

As markdown is not really expressive rendering a lot of documents in a sane way requires metadata
that is not present in the files themselves.

1. If an `index.md` file is present it will be rendered. If not found the directory is sorted
   alphabetically and the first entry is rendered.

## Limitations

Re-downloads everything every time. Need to store last commit and retrieve (implement poor man's
git).

Using a storage abstraction (+memory caching) would be nice. Right now _everything_ is cached in
memory.

## TODO

* Think about references to files in the xdoc directory; we need to find (rewrite target?) somehow -
  should be done in the markdown renderer or source conventions?
* Think about references to images, same story as for files?
* Header/footer and stuff like that.
* Add codelab support, integrate 'claat' or just use the bits for parsing.
