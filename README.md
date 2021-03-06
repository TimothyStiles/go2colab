# go2colab

Scientists (my [main project's](https://github.com/TimothyStiles/poly) users) love jupyter notebook tutorials

[pkg.dev.go](https://pkg.go.dev/)'s runnable playground [doesn't support file IO](https://pkg.go.dev/io/ioutil#example-ReadFile) but I love example tests and the fact that they update with every new release (guaranteeing docs never rot and I never need to manually update them).

So I'm writing this program that rips an entire project's example tests that include the keyword "tutorial" and packages them into jupyter notebooks that get pushed to a hosting service then update tests in godocs with associated urls. 

Is this a dirty hack to get runnable IO examples? Maybe, but I'll be danged if I need to maintain tutorials on top of my libraries.

### Spec
Given a `Go` package repo url `go2colab` will:

- [x] Store the `Url`
- [x] Extract and store the and store `Repo`, `Owner`, `Branch/Commit`, and `Host`'s names from the `Url`
- [x] if the `Branch/Commit` isn't supplied default to commit associated latest release tag
- [x] Clone the git repo into a temporary directory
- [x] If no `Branch/Commit` is supplied extract the latest release tag from repo with its associated commit.
- [x] Parse `go.mod` for `Go` version
- [x] Find every `example_test.go` file path in the repo and its subpackages and store those in `Repo.Paths`
- [x] Iterate through example paths and extract every example test with the word "tutorial" in its definition and parse it into a `Tutorial` struct.
- [x] For each `Tutorial`:
  - [x] Initalize a `Notebook` struct with required metadata. (an [autorunning notebook cell](https://coding-stream-of-consciousness.com/2018/11/13/jupyter-auto-run-cells-on-load/) that sets up the go kernel and env)
  - [x] Convert `Tutorial.Source` to a `Notebook.Cell` and append it to `Notebook.Cells`
  - [ ] If flag - Write the notebook to a file in same directory as source example test
  - [ ] If flag - Push the notebook to hosting service with Go runtime
  - [ ] If flag - Update godocs with the tutorial's url
  - [ ] If flag - Update godocs with the tutorial's url

### What this is not (unless y'all start throwing money at my [github sponsors](https://github.com/sponsors/TimothyStiles/)).

- This does not convert colab notebooks to Go example tests
- This is not a doc hosting solution
- This is not a product
