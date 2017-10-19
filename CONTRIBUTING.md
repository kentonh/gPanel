# Contributing to gPanel

A quick way to get started contributing to this project is to check out our active issues in the issues tab of the master repository. If you would like to contribute to the wiki (wiki contribution guidelines are found [here](#wiki_guidelines_anchor)), that is also an easy to get started contributing. A more advanced way of contributing is to take a look at our [features list](#features_list_anchor) and pick one that needs work, or start working on an unimplemented feature. Feel free to edit this document in your own fork and submit a pull request to add to, subtract from, or make minor adjustments to existing features on the list.

## Stylistic & Conventions Guidelines

We follow all naming and style conventions that are officially put out by the developers of Go themselves. All of these rules and conventions can be found in their [Effective Go](https://golang.org/doc/effective_go.html) tutorial.

As this project is open source we do encourage commenting wherever may be necessary and require comments for all function declarations. The reason this is important is because it will increase productivity to whoever wants to come along and either add to or improve your code later on.

Go has [testing baked in the language](https://golang.org/pkg/testing/) and we would like there to be a test for nearly, if not all, functions that are in our program. If you don't write tests for your functions that you contribute, please note in the pull request that you did not.

### Examples of Well-documented Code

#### Function Blocks

```go
/*
  A function used to add two integers who must be positive

  @param one int
    The first integer term of the to-be sum, must be positive
  @param two int
    The second integer term of the to-be sum, must be positive

  @return int
    The sum of the two input integers
  @return error
    The resulting error, if there is none it returns as nil

  Known Problems: Can't handle negative integers
*/
func addTwoPositiveIntegers(one int, two int) (int, error) {
  var answer int

  if one < 0 || two < 0 {
    return answer, errors.New("One or both of the input integers are negative")
  } else {
    answer = one + two
    return answer, nil
  }
}
```

This is how a typical function block comment should look. A short description of the function, followed by parameter explanations, then return explanations, and finally, if applicable, a "Known Problems" section. The known problems section is important because it is an easy way to set a flag of sorts for other contributors to come and fix what you couldn't figure out or didn't have time to do.

## <a name="wiki_guidelines_anchor"></a>Internal Wiki Conventions & Guidelines

To-do...

## <a name="features_list_anchor"></a>Features list

Feel free to edit this document in your own fork and submit a pull request to add to, subtract from, or make minor adjustments to existing features on the list.

Key: __Implemented__ | __*Implemented, but needs work*__ | *Unimplemented* | ~~Removed~~ | Section Heading

* Administrative Web Hosting Panel
  * Accessibility
    * Using a special port
    * Authentication
    * User system
  * Clients
    * Multi-client support
    * Configuration of new clients
    * Configuration of existing hosts
* Client Web Hosting Panel
  * Accessibility
    * Using port 2082
    * Authentication
    * User system
  * Public Website Control
    * Turning access on/off
    * Graceful shutdown
    * Maintenance mode
    * IP Filtering
  * Statistics
    * Various graphs for usage, bandwidth, etc
    * Click heat maps
  * Diagnostics
    * Smart logging
    * Alert system for fatal errors
  * Mail Servers
  * Domains
    * Configuration of domains
    * Configuration of sub-domains
    * TLS/SSL certificate support
    * Multi-domain support
  * Remote Access
    * SSH Access
  * File Manager
    * CRUD
    * Inline Editor
    * Managing Permissions
* Public Website
  * Serve Requests (duh)
  * Concurrent Processing
  * Supported content types
    * All of the obvious ones (.jpg/.html/.css/etc)
    * .go
    * .php
* General
  * Deployment
    * Binary
    * GUI Installation Helper
  * Multi-infrastructure Support
    * Windows
    * Linux
    * Mac
    * Amazon Web Services
    * DigitalOcean
    * Google Cloud Platform
