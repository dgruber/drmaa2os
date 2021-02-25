# Contributing

It is really great that you are reading this! There are so many awesome 
workload managers around where a standardized Go API for job management
would be useful, but only a few of them get currently attention inside
the drmaa2os project. I'm happy to accept any contribution to the project.

## Contribution Requirements

In order to be able to accept contributions there are just a few requirements.

- Testing. Tests should be written using the great _ginkgo_ library. It helps
a lot to reduce test code and instead produce more test cases.
- Contributor agreement. Any (significant) code change requires a contributors
license agreement. Please find them in the contributing subirectory.
- Please check if the functionality is according to the DRMAA2 standard 
spec. If there is no matching concept it can be added as an extension. Most 
structs define a ExtensionList (technically a string map).

## How to Send a Signed Agreement?

The contributor agreement can be sent as email to info @ gridengine. org.

Any question can be sent to this email address as well.

## Submitting Changes

If you plan to add something please use the standard git workflow for that:
- Clone repo
- Create new branch
- Add changes / ideally summarize as one commit
- Create a pull request
