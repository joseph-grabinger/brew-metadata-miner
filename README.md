# brew-metadata-miner
Mining of dependencies and associated metadata of formulae from the [HomeBrew core](https://github.com/Homebrew/homebrew-core).


## Configuration

The configuration is done in the `config.yml` file. The file contains the following fields:
   * `output_dir`: The directory where the output file will be stored.
   * `core_repo`:
       * `url`:  The URL of the HomeBrew core repository.
       * `branch`: The branch of the HomeBrew core repository.
       * `dir`: The path to the core repository.
       * `clone`: A boolean value indicating whether the core repository should be cloned or not.
   * `max_workers`: The maximum number of workers to use when reading the formulae.


## Export format of the metadata

The extracted metadata is  stored in a TSV file where it is represented in the following format: 

```sh
0  "<package_manager>"  "<name>"  "<license>"  "<namespace>/<username>/<repository>"  "<stable_archive_url>"  "<system_requirement>"
1  "<package_manager>"  "<name>"  "<license>"  "<type>"  "<system_restriction>"
...
```

A leading zero indicates a package line, whereas a leading one indicates a dependency line.


