# Overview of Usage

## Qazaar Taxonomy and Rules

The projects utilize a series of taxonomy and rules that allow for constructing a picture of all the projects, and from there extracting actionable information.

The Taxonomy describes what meta-data all the documents can define.  By assembling these in a central location (or how to translate them from one form into another), the documentation has a guide for the shape of its meta-data.

The rules (needs a better name) then allows for turning the taxonomy into usable data.

* "Desires" come from requirements and test plans.  The rules around desires mostly comes from the allowed format of the extracted, computer-usable data.  These come in two forms.
    * Requirements - a thing that should exist.  For example, desire a service that allows editing the User Profile structure", or "Desire tests that validate that the Charlie Jones persona cannot change another user's password."
    * Validations - when something matches a set of traits, there must exist other things for it.  For example, all user input must be cleansed for SQL and XSS injections.
* "Cohesive Groups" define traits that indicate different things that belong to the same "cohesive group", and other traits the they must have in common, and violations of the other traits indicate an error.  For example, a "structure" similar group requires that they have a "type" set to "structure" and they belong to a cohesive group when they have a shared "name" value, then their contained "fields" must all align.

## Collecting Documentation

The collected documentation 
