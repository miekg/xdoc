# xDoc

xDoc will traverse all of GitLab -- groups/subgroups and projects that it has access to and will
look in each repository for:

* a top-level file named `.xdoc.yaml` that allows for setting some scrape options.
* a directory named xdoc (TBD: make this an option)

It will both render and index these markdown files, the following endpoints are available:

* `/search`; exposed the bleve search endpoint.
* `/groupname/subgroupname/projectname`; hitting this endpoint will render the `index.md` markdown file.
* `/groupname/subgroupname/projectname/filename`; hitting this endpoint will render the `<filename>.md` markdown file.

More endpoints may be added in the future; for instance listing all repositories or things like
that.

## Conventations

As markdown is not really expressive rendering a lot of documents in a sane way requires metadata
that is not present in the files themselves.
