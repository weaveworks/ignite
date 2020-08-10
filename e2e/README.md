This is the e2e test package for Ignite

These tests are verified in Continuous Integration builds.

Run using make:
```
make e2e
make e2e E2E_REGEX=TestVolume
make e2e-nobuild E2E_REGEX=TestVolume
```

How to run the test suite manually:
```
sudo IGNITE_E2E_HOME=$PWD $(which go) test ./e2e/. -v -count 1 -run Test
```
