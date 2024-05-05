# brew-metadata-collector
Metadata collection of formulae from the [HomeBrew core](https://github.com/Homebrew/homebrew-core).


## Configuration

The configuration is done in the `config.yml` file. The file contains the following fields:
   * `output_dir`: The directory where the metadata will be stored.
   * `core_repo_url`: The URL of the HomeBrew core repository.
   * `core_repo_branch`: The branch of the HomeBrew core repository.
   * `core_repo_dir`: The path to the core repository.
   * `core_repo_clone`: A boolean value indicating whether the core repository should be cloned or not.
   * `max_workers`: The maximum number of workers to use when reading the formulae.


## Export format of the metadata

The extracted metadata is represented in the [pkg-deps-fmt](https://github.com/joseph-grabinger/pkg-deps-fmt) and stored in either a CSV or TSV file.

