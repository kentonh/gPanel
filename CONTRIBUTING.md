# Contributing to gPanel

Todo general paragraph...

## Stylistic & Conventions Guidelines

We follow all naming and style conventions that are officially put out by the developers of Go themselves. All of these rules and conventions can be found in their [Effective Go](https://golang.org/doc/effective_go.html) tutorial.

As this project is open source we do encourage commenting wherever may be necessary and require comments for all function declarations. The reason this is important is because it will increase productivity to whoever wants to come along and either add to or improve your code later on.

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
  var answer int;

  if one < 0 || two < 0 {
    return answer, errors.New("One or both of the input integers are negative")
  } else {
    answer = one + two
    return answer, nil
  }
}
```

This is how a typical function block comment should look. A short description of the function, followed by parameter explanations, then return explanations, and finally, if applicable, a "Known Problems" section. The known problems section is important because it is an easy way to set a flag of sorts for other contributors to come and fix what you couldn't figure out or didn't have time to do.
