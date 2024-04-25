# Qazaar Data Exchange

The underpinnings for Qazaar.  The tools must unify the data into compatible data exchange formats.

The project uses JSON documents as the primary exchange format, along with JSON Schema descriptions of those formats.

## Basic Information

All of the information gathered by the data exchange format requires:

* **Schema version**: the data format itself has a schema, which may change over time.  By providing this schema information, tools may better use the included data.
* **Source version**: the data comes from a source, and that source should have a specific version.  The data may come from multiple sources, and each one of those has a version.  By pinning the collected data to the source + version, readers can better understand the age of the data.  It also allows for comparing changes.

## Schema Basics

All schema definitions currently use JSON schema 07-draft format.  All schema must use an Object as the top-level item, and allow author comments.

Every numeric property must include a minimum and maximum value.  Every string property must include a maximum length.  Every array property must include a maximum length, along with methods for defining how to construct data if the number of items exceeds that length.

Date-time information always includes a Timezone value.

## Supported Schema

### [Test Execution](schema/test-execution.v1.schema.json)

Describes the results of running zero or more tests.

Each test in the document includes data such as:

* The test case name or identifier
* The test case source + version
* The date/time/timezone of the test execution
* The time taken to run the test
* The number of times the test was run
* The execution result (pass, fail, etc)
* The source code lines or instructions executed by the test

In some cases, the document may include, or only contain, a summary of test execution information.  In this case, the summary should also list the identifying information for the summarized tests.
