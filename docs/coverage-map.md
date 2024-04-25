# Coverage Maps

A project should have a *feature list* that it should aim to fulfil.  It may not implement them all, but at least it has a roadmap of what the product owners want it to do.

However, that does not reveal the whole picture of what a tester needs to cover.  A system can have many other aspects to it than just the raw feature list.  It can require environmental resilience conditions, historical bugs (to prevent regressions), combinations of features that may trigger odd behavior, and so on.

To this end, a *coverage map* joins together the different features and other areas that the testing team plans to cover.

At a high level, a coverage map closely ties with the concept of a *test plan*.  A coverage map, though, implies a living organizational construct that ties active tests back to requirements.  It allows for continuing to build the plan deeper into details until it eventually describes the test cases.

As a map, it itself does not live in a single document.  It can spread across multiple documents and into automated test code, and even across projects, and include other coverage maps.  It has a unified concept of categories, tags, and labels.

A coverage map helps you define all the risks to releasing.  Anything that exists in that map that you have not tested indicates an *unknown risk*.  If it doesn't represent that, then you probably do not need it.

## Test Object Taxonomy

A *test object* refers to something within the coverage map that describes a test.  A test plan, a test case step, an automated test suite.  All of these and more.

A *test object* does not refer to a feature or a bug.  A test object covers those.

The coverage map allows a hierarchy of test objects, along with a *taxonomy classifiers* for the objects.  The *taxonomy* allows for all the ways to organize and cross-cut the test objects.

The taxonomy classifiers can cover all test objects in a hierarchial structure, except some underneath it may not apply to some classifiers, while others only apply to a single test object.

### Category

A *category* refers to a named item, which test objects provide zero or 1 value for, from a fixed list of known values.  Some categories may put a restriction to require 1 value.

For example:

* **test-type**: could have known values 'functional', 'performance', 'longevity', 'end-to-end', 'smoke'.
* **environment**: could have known values 'firefox' and 'chrome', or maybe 'android' and 'ios'.  As a detail, it would only apply to certain test runs, rather than the whole test case.

### Category Set

A *category set* is similar to a Category, but allows zero or more values.  The category set may include extra rules for values in the set, such as exclusionary values.

### Label

A *label* refers to a named item, and a value for the item.  The label does not define a list of allowed values.

### Label Set

### Tag Set

### Meta-data

Each one of these items can have corresponding meta-data.  For example, a "ticket" label could include meta-data that transforms the values into HTTP links to a ticketing system.


## Transformations

Your test plan may need to cover multiple hierarchies, each with their own taxonomy.  If you have lots of time and complete ownership of those hierarchies, you can restructure their taxonomy under a unified umbrella.  However, that's usually not the case.  Additionally, each hierarchy may have very good reasons to use something unique.

To support a single, unified vision while having different kinds of hierarchies, they can support transformations - either from the child to the parent or vice versa (using both leads to madness - you'll be mad at yourself for doing it).

### Simple Translation

In the most simple case, you may translate a label name to another, or map the values to another value.

### Complex Transformation

Example: System A uses a label set "bug", with each value referencing a ticket number, and a "feature" label set with each value referencing a feature ID.  System B uses a label set named "ticket", with a full ticket ID + number.

While you may find a way to query the ticket system for the full ID to discover whether it references a feature or bug, and transform from System B to System A, you may find it easier to transform from System A to System B by turning the bug into a ticket ID + add a tag for "bugfix" (it tests a bug fix), and turn the feature into a ticket ID for the feature source + add a tag for "feature".
