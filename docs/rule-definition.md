# Defining Rules

Qazaar structures the assembled data as catagories and labels, but it needs a way to match these against other sets of data.

Rules contain two parts - the [descriptors](#item-descriptors) an item must exhibit for the rule to [apply](#descriptor-matching) to it, and the [implications](#implications) for items whose traits apply to the rule.


## Item Descriptors

Before we can describe how rules work, we must describe the things that what rules work with.  The tooling that runs the rules collect individual items with associated *descriptors*.  The rules, then, operate over the descriptors within individual items.

The item descriptors all have a form of a descriptor key and value.  The key defines the allowed set of values, called a *taxonomy*.  The kinds of values that a taxonomy can associate with a specific key include *enumerated* (values must come from a limited list of possibilities), *free* (any value), or *numeric*, and must also specify how many distinct values the key can contain.


## Descriptor Matching

### Matching Descriptors

A rule first contains a matching collection of descriptors.  This has one or more descriptor keys, along with a matching pattern.  This can also include a "not" expression, meaning that if a matching descriptor exists within the item, then the rule does not apply to it.

A matching descriptors collection follows basic boolean logic between expressions ('and', 'or', and 'not' operators with grouping).  Each expression matches a specific descriptor key against its value, not against values of other descriptors.

An expression may take the form of "present", meaning that the descriptor has some value assigned to it.  By implication, the 'not' grouping operation means "no value present".

For numeric type values, the allowed operations include only "equal" to a precise value, or "within" a range.  This implies "outside" and "not equal" when used with the 'not' grouping operation.  If a numeric descriptor includes multiple values, then this can also take the form "all within" or "any within", and "equal" here means "any value equal".

Non-numeric values must always take the form of either "contains some", "contains all", or "contains only".  Containment comparison may use equality or a non-backtracking regular expression.


### Cohesive Group Descriptors

Cohesive groups define impromptu rules based on items with matching descriptors.  Rather than defining a set of fixed values for the matchers, they instead define descriptor keys which match between them.  For example, of the descriptor key for the cohesive group uses `structure-name` with equal values, then the items which include that descriptor key and their values equal all belong to the same cohesive group.  Multiple items may have different values, and each of those different values form separate cohesive groups.

Cohesive groups may include matching descriptors to narrow down the allowed items in the cohesive group.  However, members to a group still self-describe the constructed group, not the rules.

(need detail on how matching occurs between descriptor values)


## Implications

Similar to how a rule includes a definition for whether an item must follow the rule, the rule includes implications for matching the rule.  These describe expectations.

The matching expectation works just like the matching descriptor.


### Rigid Implication




### Shared Trait Implication
