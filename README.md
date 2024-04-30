# Qazaar

## The Test Organizer Tool

## A Strong Dose of Quality Courage

Qazaar gives you the capability to organize your test plans and test cases together in a cohesive view, without telling you how to do it.

A quality team generally exists to answer these questions:

* "Can we ship it?"
* "What's at risk if we ship it now?"
* "What's left to do so we can ship it?"
* "What do we need to cut so we can ship it?"

Your test plans, test cases, bug lists, pass/fail reports, and code coverage numbers all hint at answers.  Qazaar aims to unify those to give you an answer.

A solid tool should also help guide the quality team to understand how to update the tests as the project changes its features, and new bugs come in.

If the tool does its job right, it can then examine the source code for changes since the last release, compare that to the coverage report, link tests to code, and provide suggestions for what tests to rerun.


## General Approach

Qazaar concerns itself first and foremost with the definition of data exchange formats.  With a single format to contain all the various forms of information, other tooling can build upon it to generate flexible capabilities.

Secondly, Qazaar attempts to define many different ways functional requirements express themselves, along with how tests indicate which requirements they cover.  The combination of these two allow for better understanding what has coverage and what the tests miss.


## A Use Case

Let's say we have a 3-tier web site (HTML + JavaScript user interface, a server, and a database).  We have a new feature to allow a user to edit their profile to include a bit of flavor text.  This allows users to get to know each other.

With this feature, we introduce some requirements to our product:

* A user can edit their own profile flavor text.  The flavor text can include up to 4,000 characters.
* Other authenticated users can read the profile flavor text of any other user.

The team reviews these requirements, and decides that this turns into changes to the existing product:

* User Interface Updates:
    * Add a new Flavor Text text area for Edit Profile actions.
    * Include the Flavor Text in the data sent to the server on a save action.
    * Add a new Flavor Text display area to the View Profile page.
* Server Updates:
    * Receive the Flavor Text field from Update Profile requests.
    * Send the Flavor Text field on Get Profile requests.
* Database Updates:
    * Add a new text blob column to the User Profile table.

Based on the feature and the requirements, the development team codifies these with:

* Field: User Profile - Flavor Text
    * Parent Record: UserProfile/username
    * Data type: text area
    * Editable: true
    * Maximum size: 4,000
    * Edit Access: owner
    * Read Access: any authorized user
    * Edit Label Name: "Describe Yourself"
    * View Label Name: "About Myself"

However, the product already has existing, *cross-cutting* requirements across the product:

* User Interface:
    * All text areas must include a label, whose value must pull from the supported localized translations.
        * Codified as: If has "Label", then Value must include localization.
    * All editable text fields must support Unicode characters.
        * Codified as: If "Data type" = "text area", then content = Unicode
    * All text areas must contain zero or more characters, and fewer than *per-field maximum* characters.
        * Codified as:
            * If "Data type" = "text area", then field must include 'maximum size' > 0.
            * If "Data type" = "text area", then UI field includes maximum length = field 'maximum size'.
    * The user interface must scrub the displayed text to not allow for cross-site scripting (XSS) vulnerabilities.
        * Codified as: If "Data type" = "text area", then UI escapes characters safe for rendering.
* Server:
    * All text fields received must use UTF-8 encoding.
    * The text fields must not exceed *per-field maximum* characters.
    * All text sent to the database must use proper parameter SQL commands to avoid SQL injection attacks.
    * Requests for saving data require *per-request level* authorized access.
    * Requests for sending data require *per-request level* authorized access.
    * Requests must respond within 2 seconds under heavy load.
* Database:
    * Columns storing more than 2000 Unicode characters must use binary blobs, encoded as UTF-32BE.
    * Binary blobs must set *per-field maximum* byte count, times 4 (because 4 bytes store 1 32-bit character).

On top of this, the development team has a series of items that it created because of discovered bugs and limitations in the software they depend on:

* All Unicode text must support 32-bit characters, such as emoji.
    * Codified as: If content = Unicode, then character range = valid unicode characters. (You can find the precise allowable values elsewhere; it's complex)
* All Unicode text can support up to 4 diacritic marks per displayed glyph.  The base character counts as one character, and each diacritic mark counts as another.

From these various data sources, we can now add meta-data to our tests to see what our tests cover.


## Forcing Functions

This tool came about because, in the opinion of the authors, the state of art for test case management tools revolves around keeping people locked into a single vendor, rather than working with a team's own processes.

Specifically, the features desired by the team include:

* A tool, not a workflow.
* A tool to organize the tests and test documentation.
* Let the source control handle branching and merging; that's what they do, so let them do it.  This means no central database, and working disconnected.
* Don't require automation authors to go into the tool to point it at the automation.
* Make it easy to turn manual steps into automation.
* Don't try to be an everything tool; let the team pick the right tools for the job.
* If I have a report, I need it to contain information I can act upon.  "80% Covered" is a nice number, but doesn't tell me whether we can release or not.
* Many different ways to categorize and itemize tests.  My team may have 'load tests' and 'functional tests' while another team wants 'attack vectors', but a test can only be one, while a test can cover both 'login' and 'logout' features, and I may want to find tests that someone marked as 'flaky'.  Not only do I want this, but I also want a sane way to keep track of these.
* I'd like to have it connect to my other tracking systems.  However, just because some tool doesn't have support for another tool, I still want to connect them.


## License

Qazaar uses the [Apache-2.0 license](LICENSE) for all its software.  Contributors must release their software under that license.  Currently, the project also releases its documentation under the same license.

```
   Copyright 2024 Qazaar-Testing Members

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
```
