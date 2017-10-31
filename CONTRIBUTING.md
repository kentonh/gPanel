# Contributing to gPanel

A quick way to get started contributing to this project is to check out our active issues in the issues tab of the master repository. If you would like to contribute to the wiki (wiki contribution guidelines are found [here](#wiki_guidelines_anchor)), that is also an easy to get started contributing. A more advanced way of contributing is to take a look at our [features list](#features_list_anchor) and pick one that needs work, or start working on an unimplemented feature. Feel free to edit this document in your own fork and submit a pull request to add to, subtract from, or make minor adjustments to existing features on the list.

## Stylistic & Conventions Guidelines

We follow all naming and style conventions that are officially put out by the developers of Go themselves. All of these rules and conventions can be found in their [Effective Go](https://golang.org/doc/effective_go.html) tutorial as well as [this Godoc](https://blog.golang.org/godoc-documenting-go-code) on documenting Go Code.

As this project is open source we do encourage commenting wherever may be necessary and require comments for all function declarations. The reason this is important is because it will increase productivity to whoever wants to come along and either add to or improve your code later on.

Go has [testing baked in the language](https://golang.org/pkg/testing/) and we would like there to be a test for nearly, if not all, functions that are in our program. If you don't write tests for your functions that you contribute, please note in the pull request that you did not.

### Examples of Well-documented Code

#### Function Blocks

```go
// addTwoPositiveIntegers returns the sum of two positive integers.
// When the two integers are not positive, it returns the zero-value of an integer and an error.
// BUG(george-e-shaw-iv) This function cannot handle negative integers, this needs addressed
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

#### Types, Variables, and Constants

```go
// A descriptive, custom data type for storing integers to be used in various time operations.
type Seconds int

// Variables don't always need comments, their name should give a good enough hint to what it does.
var cake string = "lie"

// These constants are used in various math and physics operations
const(
  PI  float32 = 3.14
  TAO float32 = 6.28
)
```

#### Packages

Package names should always match their parent folder.

```go
// Package networking contains various functions used to communicate between networks and
// draw data from the client network.
package networking
```

## <a name="wiki_guidelines_anchor"></a>Internal Wiki Conventions & Guidelines

To-do conventions and guidelines...

#### Wiki pages that don't exist but need to:

* Supported content types

#### Wiki pages that exist but need more content:

* N/A

## <a name="features_list_anchor"></a>Features

Feel free to edit this document in your own fork and submit a pull request to add to, subtract from, or make minor adjustments to existing features on the list.

Key: __Implemented__ | __*Implemented, but needs work*__ | *Unimplemented* | ~~Removed~~ | Section Heading

* Administrative Web Hosting Panel
  * Accessibility
    * *Using a special port*
    * *Authentication*
    * *User system*
  * Clients
    * *Multi-client support*
    * *Configuration of new clients*
    * *Configuration of existing hosts*
* Client Web Hosting Panel
  * Accessibility
    * __Using port 2082__
    * __*Authentication*__
    * *User system*
  * Public Website Control
    * *Turning access on/off*
    * *Graceful shutdown*
    * *Maintenance mode*
    * *IP Filtering*
  * Statistics
    * *Various graphs for usage, bandwidth, etc*
    * *Click heat maps*
  * Diagnostics
    * *Smart logging*
    * *Alert system for fatal errors*
  * Mail Servers
  * Domains
    * *Configuration of domains*
    * *Configuration of sub-domains*
    * *TLS/SSL certificate support*
    * *Multi-domain support*
  * Remote Access
    * *SSH Access*
  * File Manager
    * *CRUD*
    * *Inline Editor*
    * *Managing Permissions*
* Public Website
  * Accessibility
    * __Port 3000 for Development__
    * *Port 80 for Production w/o SSL/TLS*
    * *Port 443 for Production w/ SSL/TLS*
  * __Serve Requests (duh)__
  * *Concurrent Processing*
  * Supported content types
    * __*All of the obvious ones (.jpg/.html/.css/etc)*__
    * *.go*
    * *.php*
* General
  * Deployment
    * *Binary*
    * *GUI Installation Helper*
  * Multi-infrastructure Support
    * __Windows__
    * __Linux__
    * __Mac__
    * *Amazon Web Services*
    * *DigitalOcean*
    * *Google Cloud Platform*
