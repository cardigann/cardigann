Thanks for submitting a pull request! If you aren't submitting a new indexer definition, you can delete all this text and write a summary of your change.

For a new indexer, please follow this checklist:

- [ ] Run `cardigann test definitions/trackername.yml` and include the output here:

```
Definition file for alpharatio parsed OK ✓
Testing indexer alpharatio
Indexer alpharatio passed ✓
``` 

- [ ] Post the version of cardigann you tested this with (from the footer of the web interface or `cardigann --version`)
- [ ] Make sure to add the indexer to the list in the README

If you run into issues with testing, remember you can run a debug test with `cardigann test --debug --cachepages definitions/trackername.yml`. 
