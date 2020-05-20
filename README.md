# xDoc

xDoc will traverse all of GitLab -- groups/subgroups and projects that it has access to and will
look in each repository for:

* a top-level file named `.xdoc.yaml` that allows for setting some scrape options.
* a directory named xdoc (TBD: make this an option)

It will both render and index these markdown files, the following endpoints are available:

* `/s`; exposed the bleve search endpoint.
* `/s/<query>`; directly search for `<query>`.
* `/r/<groupname>/<subgroupname>/<projectname>`; hitting this endpoint will render the `index.md` markdown file.
* `/r/<groupname>/<subgroupname>/<projectname>/<filename>`; hitting this endpoint will render the `<filename>.md` markdown file.
* `/r/<groupname>/<projectname>`; hitting this endpoint will render the `index.md` markdown file.
* `/r/<groupname>/<projectname>/<filename>`; hitting this endpoint will render the `<filename>.md` markdown file.

More endpoints may be added in the future; for instance listing all repositories or things like
that.

## Conventions

As markdown is not really expressive rendering a lot of documents in a sane way requires metadata
that is not present in the files themselves.

1. If an `index.md` file is present it will be rendered. If not found the directory is sorted
   alphabetically and the first entry is rendered.

## Limitations

Re-downloads everything every time. Need to store last commit and retrieve (implement poor man's
git).
