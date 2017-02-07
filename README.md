# split_tests

Splits a test suite into groups of equal time, based on previous tests timings.

This is necessary for running the tests in parallel. As the execution time of test files might vary drastically, you will not get the best split by simply dividing them into even groups.

## Compatibility

This tool was written for Ruby and CircleCI, but it can be used with any file-based test suite on any CI.

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

### Naive split by file count

In the absence of prior test times, `split_tests` can still split files into even groups by count.

```
rspec $(split_tests)
```

## Arguments

``` plain
$./split_tests -help

  -circleci-branch string
        Current branch for CircleCI (or set CIRCLE_BRANCH) - required to use CircleCI
  -circleci-key string
        CircleCI API key (or set CIRCLECI_API_KEY environment variable) - required to use CircleCI
  -circleci-project string
        CircleCI project name (e.g. github/leonid-shevtsov/split_tests) - required to use CircleCI
  -glob string
        Glob pattern to find test files (default "spec/**/*_spec.rb")
  -help
        Show this help text
  -junit
        Use a JUnit XML report for test times
  -junit-path string
        Path to a JUnit XML report (leave empty to read from stdin)
  -line-count
        Use line count to estimate test times
  -split-index int
        This test container's index (or set CIRCLE_NODE_INDEX) (default -1)
  -split-total int
        Total number of containers (or set CIRCLE_NODE_TOTAL) (default -1)
```

## Compilation

This tool is written in Go.

* Install Go.
* Install the [Glide](https://glide.sh) dependency manager.
* Checkout the code
* `glide install`
* `make`

* * *  

(c) [Leonid Shevtsov](https://leonid.shevtsov.me) 2017
