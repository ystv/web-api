# Architecture

This document aims to assist developers on understanding the codebase.

This application follows the MVC architecture (sort of) at a high level.
With the `controllers` packages containing the code handling HTTP responses
and putting it into something Go can handle.
The service packages are
where all the business logic is contained.
Each package has a set domain.

The current aim with the service package is to have a set of core packages:
* vod (handling all Video on Demand queries)
* campus (campus related activities)

These core packages aim to not depend on each other. Then there are the
hybrid packages which use a mixture of the cores and their own business
code. Examples include:
* creator
* public