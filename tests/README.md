# Layer0 Smoketests

We use (https://github.com/sstephenson/bats) for Layer0 smoketests; the docs are simple and so are the tests. 

# Local Config

* Environment Variables must be populated with the contents of `l0-setup endpoint -i`
* Requires an existing layer0 install (`>=v0.7.1`)

# Running

From the `layer0/tests` directory:

```
./test.sh
```

# Tips and Tricks

When adding new smoketests, always add a teardown section ala the existing tests. We suggest adding the `--wait` flag to commands that delete entities, as this ensures the tests will pass in a modular fashion. 
