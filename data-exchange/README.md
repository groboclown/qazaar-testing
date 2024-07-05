# Qazaar Data Exchange

The underpinnings for Qazaar.  The tools must unify the data into compatible data exchange formats.

The project uses JSON documents as the primary exchange format, along with JSON Schema descriptions of those formats.

These schema exist to satisfy programs as the primary consumer and producer for documents, not people.  The project expects other programs to turn these documents into more usable data.

## Basic Information

All of the information gathered by the data exchange format requires:

* **Schema version**: the data format itself has a schema, which may change over time.  By providing this schema information, tools may better use the included data.
* **Source version**: the data comes from a source, and that source should have a specific version.  The data may come from multiple sources, and each one of those has a version.  By pinning the collected data to the source + version, readers can better understand the age of the data.  It also allows for comparing changes.


## Supported Schema

### [Ontology](schema/ontology.v1.schema.json)

Defines the allowed descriptors, and how documents must use them.

Each project must have an ontology to validate conformity for descriptor usage.  For cases with multiple projects that need joining together, use the [ontology transform](schema/ontology-transform.v1.schema.json) document.


### [Document Description](schema/document-description.v1.schema.json)

Describes the items and their descriptors, extracted from a document.  As this can get very exhaustive and repetitive, we expect tooling to construct these files from the document.


### [Rules](schema/rules.v1.schema.json)

Defines expectations and requirements around extracted documents.


### [Test Execution](schema/test-execution.v1.schema.json) (under development)

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


## Contributor Style Notes

* The files must have standard JSON formatting; the style does not permit JavaScript-like comments.
* The files use 2 spaces for indention.
* Files use UTF-8 encoding.
* Schema reference objects SHOULD use only a single line.
    Correct:
    ```json
    {
      "properties": {
        "value": {"$ref": "#/$defs/ValueType"}
      }
    }
    ```
    Incorrect:
    ```json
    {
      "properties": {
        "value": {
            "$ref": "#/$defs/ValueType"
        }
      }
    }
    ```
* Object types SHOULD exist in the `$defs` section.  Array and basic types can exist in-line, without references.  For common data types, even simple ones, they should break out into shared references.
* All schema definitions currently use JSON schema 07-draft format.  All schema must use an Object as the top-level item, and allow author comments.  They must exist as stand-alone documents, outside of the schema declaration.
* Every numeric property must include a minimum and maximum value.  Every string property must include a maximum length.  Every array property must include a maximum length, along with methods for defining how to construct data if the number of items exceeds that length.
* Date-time information always includes a Timezone value.
* Instead of using an object with arbitrary property names, the schema should use an array with a `{"key": "", "value": ""}` style approach.  This allows for controlling the size of the object.
* Where possible, schema types should use the [shared schema](schema/_shared.schema.json).  While not referenced directly, the references should instead copy into the schema.
* Object keys should use lower camel case.  This makes using them in programming languages easier.  Meta-data fields should have a `$` prefix.
* Every object should allow `$comment` and `$comments` fields.  In general it should also allow the `sources` field, to allow linking the item back to where the generating program discovered the value.
* Every object should mark itself as allowing no additional properties (`"additionalProperties": false`).  However, implementations should not fail if the data includes unknown fields, as that allows for better forward compatibility.
* Feel free to add your own `$comment` to the schema, but note that these exist solely for schema maintainers, not for end-users.
