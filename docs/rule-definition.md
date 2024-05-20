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

An expression may take the form of "present", meaning that the descriptor has some value assigned to it.  By implication, the 'not' grouping operation means "no value present".

For numeric type values, the allowed operations include only "equal" to a precise value, or "within" a range.  This implies "outside" and "not equal" when used with the 'not' grouping operation.  If a numeric descriptor includes multiple values, then this can also take the form "all within" or "any within", and "equal" here means "any value equal".

Non-numeric values must always take the form of either "contains some", "contains all", or "contains only".  Containment comparison may use equality or a non-backtracking regular expression.

Additionally, because all descriptor values internally use a collection of values, a rule may also include a "count" of values, which allows for numeric expressions against the number of values stored in the descriptor.  They may also use "unique" and "distinct" qualifiers if the key allows for non-distinct values.


### Self-Organizing Groups

Self-organizing group (SOG) rules define a method for constructing groups through their shared descriptor values, rather than through an explicit membership declaration.  For example, of the descriptor key for the SOG rule uses `structure-name` with equal values, then the items which include `structure-name` values equal to `user_profile` share one SOG while those hose value equal `checkout_cart` share another SOG.

A SOG contains zero or more descriptor matching descriptors.  They also must have one or more "shared value" descriptor.  The shared value descriptor takes general form "all members of the SOG have descriptor key A, and the value of A for each item in the SOG matches all other items' descriptor key A value in the SOG."

For the purposes of reporting, the descriptor keys and their shared values construct the identifier for the SOG.

The members of the SOG construct a super structure for the SOG itself.  The super structure contains all the descriptors for all members, with the values being a list of the members values.  These lists differ slightly from the above limitation on descriptor values, so that rules can check for uniqueness of member values (that is, no member may share the same value), as well as check for the union of the set (removes duplicates).  It also constructs meta-descriptors to count the number of members.

For the purposes of this document, a SOG represents one collection of members whose shared descriptors match a rule, and a SOG rule describes how to lump items into members of a SOG.

Because the SOGs construct a single super structure, SOG rules can operate on other SOGs.  Implementations should not allow for SOG rules to declare a recursive model of super structure generation, and so may require some additional restriction.


### Variable Values

(todo expand)
Rules may allow for an externally declared set of variable replacement values.


## Implications

Similar to how a rule includes a definition for whether an item must follow the rule, the rule includes implications for matching the rule.  These describe expectations.  An implication *fails* if the item does not follow the rule, and *passes* if it does.

Implications have a "level" associated with them, which must have one of a pre-defined set of values.  These may take the form of "error", "warning", "deprecated", or any number of other values, as needed by the system.  The evaluation of the implications may cause a halting error, or warnings, or may halt if the system encounters a certain criteria of accumulated implication failures.


### Conformity Implication

The simplest form of implication comes in the form of an "if-then" style rule set.  If a set of rules match an item, then another set of rules must match on the item.

This allows for enforcing setups, like "all requirements must have a requirement id."  It also allows for constructing experience into the system, such as "all input text must include SQL injection tests."

Conformity implications can apply to both simple descriptor matchers and SOGs.


### Convergent Implication

Convergent implications build upon the [SOGs](#self-organizing-groups) by enforcing restrictions within the group, by having the group itself define the restrictions.  This takes the general form of, "for all items in the SOG, they must have descriptor A match each other."

This takes advantage of the construction of a super structure from the members of the SOG.

Due to the lack of a hierarchy, ideas such as "all things that use the structure 'user_profile' must share fields with matching names and types" require one rule for declaring how fields must have specific keys match amongst members, and another rule for constructing a multi-staged SOG, one for declaring each implementation's collection of fields for the structure, and another to collect those SOG super structures into another super structure that enforces that the collection of field names must match between members.
