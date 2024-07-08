# Defining Rules

Qazaar structures the assembled data as catagories and labels, but it needs a way to match these against other sets of data.

Rules contain two parts - the [descriptors](#item-descriptors) an item must exhibit for the rule to [apply](#descriptor-matching) to it, and the [implications](#implications) for items whose traits apply to the rule.

This document describes the low level details for the workings of rules and descriptors.  Qazaar expects implementations of rule constructing to allow for a higher level language to more naturally express expectations for projects.


## Item Descriptors

Before we can describe how rules work, we must describe the things that what rules work with.  The tooling that runs the rules collect individual items with associated *descriptors*.  The rules, then, operate over the descriptors within individual items.

The item descriptors all have a form of a descriptor key and value.  The key defines:

- The allowed set of values, called an *ontology*.  The kinds of values that an ontology can associate with a specific key include *enumerated* (values must come from a limited list of possibilities), *free* (any value), or *numeric*
- Must also specify a range of values the key can contain.  It must always allow zero, as an undeclared descriptor key on an item equates to declaring the key with zero items.  The maximum might use "unlimited", although, because computers cannot store infinite data, it should instead use a numeric limit.
- Must specify whether the values within the key must be distinct.

Item descriptor keys do not form an explicit hierarchy, although they may represent an implicit one.


## Descriptor Matching Rules

### Matching Descriptors

A rule first contains a matching collection of descriptors.  This has one or more descriptor keys, along with a matching pattern.  This can also include a "not" expression, meaning that if a matching descriptor exists within the item, then the rule does not apply to it.

A matching descriptors collection follows basic boolean logic between expressions ('and', 'or', and 'not' operators with grouping).  Each expression matches a specific descriptor key against its value, not against values of other descriptors.

For numeric type values, the allowed operations include "within" a range.  This implies "equal" by setting the bounds to the same value, as well as "outside" and "not equal" when used with the 'not' grouping operation.  If a numeric descriptor includes multiple values, then this can also take the form "all within" or "any within", and "equal" here means "any value equal".

Non-numeric values must always take the form of either "contains some" (one or more of the elements in the matcher must exist in the matching descriptor's value), "contains all" (all the elements in the matcher must exist in the matching descriptor's value), or "contains only" (the matching descriptor's value must contain exactly the same values as in the matcher).  Containment comparison may use equality or a non-backtracking regular expression.

Additionally, because all descriptor values internally use a collection of values, a rule may also include a "count" of values, which allows for numeric expressions against the number of values stored in the descriptor.  They may also use the "distinct" qualifier to construct a new list with each value occurring at most once, if the key allows for non-distinct values.  To replicate the idea of "present" and "not present", the "count" of values should use a "not equal to zero" and "equal to zero", expression, respectively.


### Self-Organizing Groups

Self-organizing group (SOG) rules define a method for constructing groups through their shared descriptor values, rather than through an explicit membership declaration.  For example, of the descriptor key for the SOG rule uses `structure-name` with equal values, then the items which include `structure-name` values equal to `user_profile` share one SOG while those hose value equal `checkout_cart` share another SOG.

A SOG contains zero or more descriptor matching descriptors.  They also must have one or more "shared value" descriptor.  The shared value descriptor takes general form "all members of the SOG have descriptor key A, and the value of A for each item in the SOG matches all other items' descriptor key A value in the SOG."

For the purposes of reporting, the descriptor keys and their shared values construct the identifier for the SOG.

The members of the SOG construct a super structure for the SOG itself.  The super structure contains all the descriptors for all members, with the values being a list of the members values.  These lists differ slightly from the above limitation on descriptor values, so that rules can check for uniqueness of member values (that is, no member may share the same value), as well as check for the union of the set (removes duplicates).  It also constructs meta-descriptors to count the number of members.

For the purposes of this document, a SOG represents one collection of members whose shared descriptors match a rule, and a SOG rule describes how to lump items into members of a SOG.

Because the SOGs construct a single super structure, SOG rules can operate on other SOGs.  Implementations should not allow for SOG rules to declare a recursive model of super structure generation, and so may require some additional restriction.

SOGs may have the declaring rule also force alterations to its descriptors.  These rules may remove, add, or replace descriptor values.

While the matching engine considers SOG super structures as stand-alone entities, the user interface aspect must be able to relate the SOG super structure entity to the constituent member sources, and the rules that created it.


### Variable Values

(todo expand)
Rules may allow for an externally declared set of variable replacement values.


## Implications

Similar to how a rule includes a definition for whether an item must follow the rule, the rule includes implications for matching the rule.  These describe expectations.  An implication *fails* if the item does not follow the rule, and *passes* if it does.

Implications have a "level" associated with them, which must have one of a pre-defined set of values.  These may take the form of "error", "warning", "deprecated", or any number of other values, as needed by the system.  The evaluation of the implications may cause a halting error, or warnings, or may halt if the system encounters a certain criteria of accumulated implication failures.


### Conformity Implication

The simplest form of implication comes in the form of an "if-then" style rule set.  If a set of rules match an item, then another set of rules must match on the item.

This allows for enforcing setups, like "all requirements must have a requirement id."  It also allows for constructing experience into the system, such as "all input text must include SQL injection tests."

At a low level, conformity implications only apply to simple matching rules.  The SOG definition does not need these, as the matching descriptors for the SOG definition can create matching rules on the members, and by adding specific descriptor alterations to the generated SOG, simple matching rules can apply on items with that new descriptor.  Implementations may make a quality-of-life improvement by allowing these on the SOG definition, though.


### Convergent Implication

Convergent implications build upon the [SOGs](#self-organizing-groups) by enforcing restrictions within the group, by having the group itself define the restrictions.

These requires a descriptor evaluation to match between all members.  The supported matchers include 'all match' (the value for each member must match), and 'disjoint' (each member's value must not match any other member's value).


# Examples

## All Implementations of a Structure Share the Same Fields with the Same Types

Due to the lack of a hierarchy, ideas such as "all things that use the structure 'user_profile' must share fields with matching names and types" require a collection of rules.  However, some fields may be specific to an implementation, and if it has the `visibility = "private"` descriptor, it will not be included.

It expects these descriptor ontologies:

- `data-type` enum, whose potential values include `structure` and `field`.
- `structure` free string.
- `field-name` free string.
- `field-type` enum, whose potential values include `string` and `float`.
- `visibility` enum, whose potential values include `private`.

And that produces these rules:

- A simple rule to enforce field declaration restrictions.
  - Has matching descriptors:
    - `data-type = "field"`
  - Has conformity implications:
    - `count(structure) = 1`
    - `count(field-name) = 1`
    - `count(field-type) = 1`
- A self-organizing group relation for each field in the structure.
  - Has matching descriptors:
    - `data-type = "field"`
    - `visibility` does not include `"private"`
  - Has SOG shared value descriptors:
    - `structure`
    - `field-name`
  - Has SOG descriptor alterations:
    - set `sog-type` to `structure-field`
  - Has no convergence implication.
- A simple rule for the generated structure-field SOGs:
  - Has matching descriptors:
    - `sog-type = structure-field`
  - Has conformity implications:
    - `count(distinct(field-type)) = 1`; because the constructed SOG joins the members' values together into an array, having a 'distinct' on the field-type means that all the field-type values must be the same.
- A self-organizing group relation for each structure, built on top of the field SOGs.
  - Has matching descriptors:
    - `sog-type = "structure-field`
    - `count(structure) = 1`
  - Has SOG shared value descriptors:
    - `structure`
  - Has SOG descriptor alterations:
    - set `sog-type` to `structure`; this replaces the joined-together value from the members' super structure value to now be a different, single value.
  - Has convergence implications:
    - `field-name` values must all match.
