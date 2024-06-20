# Code, Requirement, and Tag Coverage

A common part of modern software development includes collecting *code coverage* to identify areas of the software not executed by its automated tests.  This has limitations, but allows for developers and testers to better identify areas for improvements.

With the introduction of [Qazaar tagging](ontology.md), we have the ability to provide similar coverage information regarding any tagged document.  This means we need to describe the source tags, the test tags, and the coverage information.  Additionally, the coverage information provides a linkage between different aspects of the project.

- Use case documents can trace to requirements.  This can help describe development priority by ranking the use case priority.
- Updates to software code reflect the affected use cases and requirements, which help narrow down the list of corresponding manual tests.
- Code coverage numbers generated from running automated tests can use the source and test tags to indicate requirements covered by the automated tests.

By collecting the use case and requirement tags covered by the executed tests, the testers can now construct reports that accurately detail what amount of the product has no coverage, which helps to describe current risk in releasing.
