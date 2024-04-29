# Data Sources for Qazaar

Testers can generate many different kinds of data for constructing tests.

## Test Plans

Test plans themselves usually do not include executable tests.  However, they describe approaches for testing to attempt to cover all the features and usage scenarios.

Some places that Qazaar should look for data:

* Markdown files.
* HTML files.
* Microsoft Word or OpenOffice Write documents.
* Spreadsheets.
* Wiki pages.

Data collected should identify specific test plan names, sections, or subsections, along with applicable test taxonomy classifiers.

Each project will have its own way of building these documents, so teasing out the information requires unique methods.

## Manual Test Cases


## Automated Test Cases


## Test Execution Results


## Code Coverage Results


## Feature Implications

When working on a feature for code, it often has implications of other required work and testing effort.  Some examples:

* The feature allows for an external user to communicate over the network to perform an action.  This usually takes the form of a REST API or gRPC communication, which means developing a specification for the action, creating contract tests, defining limitations on input fields and allowed data types, a minimal expectation for uptime and load capabilities, proper authorization and authentication for ensuring proper access to make the call, and more.
* A new user interface field allows for the end-user to enter some text.  This implies setting a maximum field length, a place to store the information such as a database, communication channels for sending the data, the data encoding, a label for the field, translations for the field label, and more.

A development organization can construct a set of usual checklists to ensure the feature story includes certain acceptance criteria.  This checklist implies testing as early as possible - before even a single line of code was written.  In addition, each one of those includes a set of tests to verify that the tests exist to cover the feature areas.

Feature implications indicate a series of *aspects* that cross-cut multiple parts of the application.  Using the above example, a user input field would require supporting UTF-8 characters, which has its own set of requirements.  You could picture structuring these implications as:

* Field "User Comment/User Text"
    * Has aspect: user-input text.
    * Has maximum character count of 4000.
* User Input Text
    * Requires: input sanitation;
    * Requires: UTF-8 character sets;
    * Requires: maximum character count.
* UTF-8 Character sets
    * Requires 32-bit character storage
    * Requires 32-bit character display
    * Requires BIDI-text display
* Input sanitation
    * Prevents XSS vulnerabilities
    * Prevents SQL Injection
* Database Table "User Comment"
    * Has field: "User Comment/User Text"
* User Interface Comments
    * Has field: "User Comment/User Text"

This structure allows for an implication of many tests, which, by including these in the coverage map, allows for better insight into whether the tests fully cover the required functionality.
