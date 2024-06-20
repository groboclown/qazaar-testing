# Qazaar Tagging Ontology

## Defining What You Mean

Qazaar works by extracting "tagging" information from all kinds of sources.  This tagging information, by itself, can have vague meaning when left by itself, and can lead to inconsistent use.  When the tag has different meaning to different parts of the source, the tag itself loses value and Qazaar no longer provides meaningful feedback.

In order to attempt to avoid these kinds of situations, Qazaar requires an [*ontology*](https://en.wikipedia.org/wiki/Ontology_(information_science)) on the tags to enforce a common structure and meaning.

A well-developed tag in Qazaar contains certain features:

- A unique tag identifier, or "key".
- A description to help people understand how to use it.
- A value type, which may be *numeric*, *enumerated* (choose one of a group of values), or *free text*.
- Number of values allowed.  In all cases, a tag may have zero values; if a source does not declare a value for a tag, then that equates to declaring the tag with zero values.
- Whether the tag requires only distinct values within the collection (ignores duplicate values).


## What Gets Tagged

Qazaar tracks each tag located at a specific point in the source document.  Qazaar itself does not dictate how to declare the tag, or even what a "source document" means.  This grants authors flexibility in making these decisions as best suits their needs.

Qazaar, at the moment, does not have any tooling for creating the tag inventory.  However, we expect future work to construct tooling to help this, and eventually create some common methods to create this association.


## Sharing Tags Between Projects

Projects have their own evolution path, and the tags used by each project takes on its own, individual meaning.  However, you may need to combine the tagging from multiple projects into a single view.  Just mashing them together will muddle their meanings and create noise.

To help with this, Qazaar includes mechanisms for transforming one set of Ontology into another.  This allows for a path to merge multiple project tags into a single view.  Note that this requires effort to create and maintain, and the merged result highly depends on the individuals who maintain the transformations to understand the source and targets to create something meaningful.


## Ontology, Taxonomy, and Meronomy

Qazaar very explicitly avoids placing a hierarchy on the tagged source items.  While the source itself may construct a hierarchy, and tools that extract information from the source generally takes advantage of that, the extracted information contains no hierarchy.  This means concepts such as *taxonomy* and *meronomy* do not fully apply to Qazaar's tagging method.  Instead, the more general *ontology* applies.
