# Measures and Metrics

As part of the execution of tests, they generate large amounts of *metrics* - a labeled, numeric value, usually with an accompanying unit.  For example, unit tests can generate:

* total test cases executed
* test cases failed
* test cases passed
* test cases skipped
* code coverage

Some of these can have accompanying measures to better qualify the data.  For example, code coverage can, at the most detailed level, reference a 0 or 1 value for a single instruction - did a test execute the instruction or not?  Computations exist for generating higher level values, such as total percentage of a file or module.
