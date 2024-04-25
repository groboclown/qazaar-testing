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
