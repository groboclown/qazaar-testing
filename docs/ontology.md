# Qazaar Tagging Ontology

## Defining What You Mean

Qazaar works by extracting "tagging" information from all kinds of sources.  This tagging information, by itself, can have vague meaning when left by itself, and can lead to inconsistent use.  When the tag has different meaning to different parts of the source, the tag itself loses value and Qazaar no longer provides meaningful feedback.

In order to attempt to avoid these kinds of situations, Qazaar requires an [*ontology*](https://en.wikipedia.org/wiki/Ontology_(information_science)) on the tags to enforce a common structure and meaning.

A well developed tag in Qazaar contains certain features:

- A unique tag identifier, or "key".
- A description to help people understand how to use it.
- A value type, which may be *numeric*, *enumerated* (choose one of a group of values), or *free text*.
- Number of values allowed.  In all cases, a tag may have zero values; if a source does not declare a value for a tag, then that equates to declaring the tag with zero values.
- Whether the tag requires only distinct values within the collection (ignores duplicate values).


## What Gets Tagged

Qazaar points the tags to a specific point in the source document.


## Sharing Tags Between Projects


## Ontology, Taxonomy, and Meronomy

Qazaar very explicitly avoids placing a hierarchy on the tagged source items.  While the source itself may construct a hierarchy, and tools that extract information from the source generally takes advantage of that, the extracted information contains no hierarchy.  This means concepts such as *taxonomy* and *meronomy* do not fully apply to Qazaar's tagging method.  Instead, the more general *ontology* applies.
