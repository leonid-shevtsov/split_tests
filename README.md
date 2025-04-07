# split_tests

Splits a test suite into groups of equal time, based on previous tests timings.

This is necessary for running the tests in parallel. As the execution time of test files might vary drastically, you will not get the best split by simply dividing them into even groups.

## Compatibility

This tool was written for Ruby and CircleCI, but it can be used with any file-based test suite on any CI.
Since then, CircleCI has introduced built-in test splitting. Also since then, the tool has been applied on
GitHub Actions, that does not provide native test splitting.

There is a [split-tests GitHub Action](https://github.com/marketplace/actions/split-tests) using this tool available on the Actions Marketplace.

It is written in Golang, released as a binary, and has no external dependencies.

## Usage

Download and extract the latest build from the releases page.

### Using the CircleCI API

Get an API key and set `CIRCLECI_API_KEY` in the project config.

```
rspec $(split_tests -circle-project github.com/leonid-shevtsov/split_tests)
```

(The tool returns the set of files for the current split, joined by spaces.)

### Using a JUnit report

```
rspec $(split_tests -junit -junit-path=report.xml -split-index=$CI_NODE_INDEX -split-total=$CI_NODE_TOTAL)
```

Or, if it's easier to pipe the report file:

```
rspec $(curl http://my.junit.url | split_tests -junit -split-index=$CI_NODE_INDEX -split-total=$CI_NODE_TOTAL)
```

### Naive split by line count

If you don't have test times, it might be reasonable for your project to assume runtime proportional to test length.

```
rspec $(split_tests -line-count)
```

### Apply bias

Often a specific split will not just run the test suite, but also a linter or some other quicker checks. In this case you can use the `bias` argument to balance the split better:

```
# account for 20-second linter run in split 0
split_tests -bias 0=20 -junit ...
```

The effect is that the split algorithm will assume an external delay of 20 seconds for the 0th split, and will reduce its assigned load by 20 seconds (as best it can.) Bias can be negative, too.

Don't forget to specify the same bias configuration on all runners, not just the ones that have bias.

This works best when you have real test timings (JUnit or CircleCI mode.) For splits by line count, you can still find the right bias empirically - although splits by line count are never perfectly balanced anyway.

### Naive split by file count

In the absence of prior test times, `split_tests` can still split files into even groups by count.

```
rspec $(split_tests)
```

## Arguments

````plain
$./split_tests -help

  -bias string
        Set bias for specific splits (if one split is doing extra work like running a linter).
        Format: [split_index]=[bias_in_seconds],[another_index]=[another_bias],...
  -circleci-branch string
        Current branch for CircleCI (or set CIRCLE_BRANCH) - required to use CircleCI
  -circleci-key string
        CircleCI API key (or set CIRCLECI_API_KEY environment variable) - required to use CircleCI
  -circleci-project string
        CircleCI project name (e.g. github/leonid-shevtsov/split_tests) - required to use CircleCI
  -exclude-glob string
        Glob pattern to exclude test files. Make sure to single-quote.
  -glob string
        Glob pattern to find test files. Make sure to single-quote to avoid shell expansion. (default "spec/**/*_spec.rb")
  -help
        Show this help text
  -junit
        Use a JUnit XML report for test times
  -junit-path string
        Path to a JUnit XML report (leave empty to read from stdin; use glob pattern to load multiple files)
  -line-count
        Use line count to estimate test times
  -split-index int
        This test container's index (or set CIRCLE_NODE_INDEX) (default -1)
  -split-total int
        Total number of containers (or set CIRCLE_NODE_TOTAL) (default -1)
```

## Compilation

This tool is written in Go and uses Go modules.

- Install Go
- Checkout the code
- `make`

---

(c) [Leonid Shevtsov](https://leonid.shevtsov.me) 2017-2020
````
