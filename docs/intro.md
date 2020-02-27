# Introduction to Arr.ai

| Status: **INCOMPLETE DRAFT** |
|-|

Arr.ai is many things, but it is first and foremost a data representation and
transformation language. This tutorial-style introduction will guide you through
the basic concepts underpinning arr.ai's model of data and computation and will
offer a teaser into some of it's more advanced capabilities.

## About the name

The domain name arr.ai was available and there was some irony in the fact that a
language called arr.ai doesn't have arrays (though it kind of does; see below).

## Some lexical conventions

Arr.ai has a rich syntax, which we won't dive into just yet. A few elements are
worth covering upfront to aid comprehension below.

1. **Identifiers:** Parameter and variables names start with `_`, `@`, `$` or a
   Unicode letter, and may continue with a sequence of any of these and Unicode
   decimal numbers.

   Examples: `x`, `$y`, `Username`, `i0`, `@j12`, `apple_π`

   The identifier `.` is a special case. It is often used as a default argument
   in transform expressions.

2. **Keywords:** The following names are predefined and cannot be reassigned as
   parameter or variable names: `true`, `false`, `let`

3. **Comments:** Comments start with a `#` and end at the end of the line.

   Example: `# Comment on comments.`

4. **Offset collections:** In the string `"hello"`, the first character, `h`, is
   at position zero. In the alternate form `12▸"hello"`, the `h` is at position
   12 and the remaining characters occupy positions 13&ndash;16. This is known
   as an offset-string or a string with holes. While this syntax is allowed in
   arr.ai, `▸` isn't available on most keyboards, so the syntax `"12\>hello"`
   represents the same offset-string. When printing such strings, arr.ai will
   normally use the more compact representation with `▸`. This form may also be
   used to represent offset arrays: `[12\> 1, 2, 3]`.

## Data

We start with data, because:

> *Bad programmers worry about the code. Good programmers worry about data
> structures and their relationships.* &mdash; [Linus
> Torvalds](https://www.goodreads.com/quotes/1188397-bad-programmers-worry-about-the-code-good-programmers-worry-about)

Arr.ai's data model is remarkably simple, having only three kinds of value, all
of them immutable:

1. **Numbers** are 64-bit binary floats.
2. **Tuples** associate names with values.
3. **Sets** hold sets of values.

Let's be clear about what the above means. Arr.ai has no arrays. It also has no
strings, Booleans, maps, functions, packages, pointers, structs, classes or
streams. Arr.ai has numbers, tuples and sets. There is nothing else.

But let's also be clear that this is far less restrictive than it might at first
seem. You can in fact represent:

* Arrays: `[]`, `[2, 4, 8]`
* Strings: `""`, `"hello"`
* Booleans: `true`, `false`
* Maps: `{}`, `{"a": 42}`, `{1: 34, 2: 45, 3: 56}`
* Functions:
  * Functions are unary: `\x 1 / x`
  * Binary functions don't exist, but `\x \y //.math.sqrt(x^2 + y^2)` is a unary
    function that takes a single parameter, `x`, and returns a unary function.
    The returned function takes a single parameter, `y`, and returns the
    hypotenuse of a right triangle with sides *x* and *y*.
* Packages: `//.math.sin(1)`, `//./myutil/work(42)`, `//github.com/foo/bar`

All of the above forms are syntactic sugar for specific combinations of numbers,
tuples and sets. For example, the string `"hello"` is a shorthand for the
following set:

```text
{
   (@: 1, @char: 101),
   (@: 2, @char: 108),
   (@: 4, @char: 111),
   (@: 3, @char: 108),
   (@: 0, @char: 104),
}
```

(Order doesn't matter in a set. It's the `@` attribute that determines the
position of each character in the string being represented.)

## Data transformation

Arr.ai is an expression language, which means that every arr.ai program, no
matter how complex, is a single expression evaluating to a single value. You can
play with the language on the command line by running `arrai eval <expression>`,
with `e` being a shortcut for `eval`, e.g.:

```bash
$ arrai e 42
42
$ arrai e '//.math.pi'
3.141592653589793
$ arrai e '[1, (a: 2), {3, 4, 5}]'
[1, (a: 2), {3, 4, 5}]
$ arrai e '[1, (a: 2), {3, 4, 5}](1)'  # Arrays are functions.
(a: 2)
$ arrai e '"hello"(3)'                 # So are strings.
108
$ arrai e '"hello" => (@:.@, @item:.@char)'
[104, 101, 108, 108, 111]
$ arrai e '[104, 101, 108, 108, 111] => (@:.@, @char:.@item)'
"hello"
$ arrai e '{
   (@: 1, @char: 101),
   (@: 3, @char: 108),
   (@: 0, @char: 104),
   (@: 4, @char: 111),
   (@: 2, @char: 108),
}'
"hello"
```

The last example underscores the point made earlier that strings are in fact
sets of tuples. There is no semantic distinction between the two forms.

## Expressions

### Literals

#### Core literals

The core syntax for literals can expresses numbers, tuples and sets.

1. Numbers: `0`, `1`, `-2`, `3.45e-6`, `7+8.9i`,
   `9969216677189303386214405760200`

   Integer components may be written in the following forms:

   * Decimal: `123`
   * Hexadecimal: `0x7b`
   * Octal: `0o173`
   * Binary: `0b1111011`

2. Tuples: `()`, `(a:1)`, `('t.m.o.l.': 42)`, `(x: (a: (), b: 2), y: -3i)`

   Like structs in the C family of languages, names are not values in their own
   right. They cannot be stored in variables or data structures and therefore
   cannot be manipulated as values. They serve only to specify which element of
   a tuple is being specified or retrieved.

   Unlike C structs, names can be any sequence of characters, with string syntax
   allowing characters not permitted in identifiers. Also unlike C structs,
   tuples do not have to conform to definitions stipulating the available fields
   or the types of values they can hold. A tuple can have any fields and each
   fields can hold any value of any type.

3. Sets: `{}`, `{1, 2, 3}`, `{(a:1, b:2), (a:4, b:7)}`, `{2, {}, (c:4)}`

#### Sugared literals

As explained earlier, many other structures are expressible beyond just numbers,
tuples and sets. It is important to remember that these other structures are
simply special arrangements of the base types. They do, however, give arr.ai the
flavor and power of much richer type systems while retaining a remarkably simple
data model. Also, because these sugared forms are all just the base types in
disguise, all of the expressive machinery designed for numbers, tuples and sets
can be applied to strings, arrays, etc.

##### Boolean syntax

Arr.ai takes a leaf out of the C89 playbook and omits Boolean types from the
base type systems. Nonetheless, `false` and `true` are defined in the core
language as aliases for the following sets.

1. `false = {}`
2. `true = {()}`

These are not the only values that may be used in logical operations. All values
can be tested for "trueness". Most values are considered "true". The only
exceptions are `0`, `()` and `{}`.

##### String syntax

TODO

##### Array syntax

TODO

##### Relation syntax

TODO

### Logic expressions

Arr.ai supports operations on "true" and "false" values. The values `0`, `()`
and `{}` are considered "false", while all other values are "true".

1. `expr1 if testexpr else expr2` evaluates to `expr1` if `testexpr` is "true",
   or `expr2` otherwise.
2. `expr1 && expr2` evaluates to `expr1` if it is "true" or `expr2` otherwise.
3. `expr1 || expr2` evaluates to `expr1` if it is "false" or `expr2` otherwise.

All above expressions exhibit short-circuit behaviours, which means that that
`expr2` will be evaluated if its value is needed. While the arr.ai language has
no side-effects, short-circuit behaviour is still needed to terminate recursion.

### Arithmetic expressions

Arr.ai supports operations on numbers.

1. Unary: `+`, `-`
2. Binary:
   1. Well known: `+`, `-`, `*`, `/`, `%` (modulo), `^` (power)
   2. Modulo-truncation: `-%` (`x -% y = x - x % y`)

### Structure access expressions

1. Tuple attribute: `tuple.attr`
2. Dot variable attribute: `.attr = (.).attr`
3. Function element:
   1. `[2, 4, 6, 8](2) = 6`, `"hello"(1) = 101`
   2. `{"red": 0.3, "green": 0.5, "blue", 0.2}("green") = 0.5`
4. Function slice:
   1. `[1, 1, 2, 3, 5, 8](2:5) = [2, 3, 5]`
   2. `[1, 2, 3, 4, 5, 6](1:5:2) = [2, 4]`

### Binding expressions

The following operators bind `name` to something related to `expr` (details
below) and evaluates expression `transform` with `name` in scope. The effect is
to transform `expr`.

1. **`let name = expr transform`** or **`expr -> \name transform`**: Transforms
   `expr`.
2. **`expr => \name transform`**: Transforms each element of set `expr` and
   evaluates to the set of results.
3. **`expr >> \name transform`**: Transforms each item of keyed-collection
   `expr` and evaluates to the key-collection of results, with each result being
   associated with the same key that the original item was. This works for any
   binary relation with an `@` attribute, which includes strings, arrays,
   functions and other structures.
4. **`expr :> \name transform`**: Binds `name` to each value in tuple `expr`,
   evaluates `transform` and reassociates each result with the corresponding
   name, producing a new tuple.

If `expr` is omitted, `.` is assumed.

If `\name` is omitted, `\.` is assumed.

### Relations

Relations are sets of tuples with a common set of names across all tuples. They
are analogous to SQL tables. Numerous operators exist that work on these
structures:

1. **Join:** TODO: describe
2. TODO: complete

### Functions

There are several flavors of functions. All functions are binary relations with
one attribute called `@`. The other attribute can have any name, including the
empty name, `''`. The following are some examples of functions.

1. **Strings:** `"hello"(2) = 108` (`l`)
2. **Arrays:** `[10, 15, 20, 25, 30](3) = 25`
3. **Lambda functions:** `\x 2 * x`

Unlike most other languages, arr.ai are no concept of named functions, either at
file level or any other scope. All functions are anonymous. A function can, of
course, be bound to a name via `let` or `->`, but, since it cannot refer to this
name at the moment of assignment, this presents a challenge for implementing
recursion. This problem is solved by a couple of functions in the standard
library:

1. **`//.fn.fix`** is a fixed-point combinator. It is typically used to
   transform non-recursive functions into recursive ones, e.g.:

   ```arrai
   let factorial = //.fn.fix \factorial \n 1 if n < 2 else n * factorial(n - 1)
   factorial(6)
   ```

2. **`//.fn.fixt`** is a variant of `fix` that operates on tuples of functions
   instead of a single function. This allows mutual recursion, e.g.:

   ```arrai
   let eo = //.fn.fixt((
      even: \t \n n == 0 || t.odd (n - 1),
      odd:  \t \n n != 0 && t.even(n - 1),
   ))
   eo.even(6)
   ```

In future, these functions will be available through syntactic sugar, something
like:

```arrai
let rec factorial = \n 1 if n < 2 else n * factorial(n - 1)
```

```arrai
let rec (
   even = \n n == 0 || odd (n - 1),
   odd  = \n n != 0 && even(n - 1),
)
even(6)
```

### Packages

External libraries may be accessed via package references.

1. **`//.`** Is the root of the standard library. It provides access to many
   packages providing a wide range of useful capabilities. The following is a
   small sample of the full set:
   1. **`//.math`:** math functions and constants such as `//.math.sin`
      and `//.math.pi`.
   2. **`//.str`:** string functions such as `//.str.upper` and
      `//.str.join`.
   3. **`//.fn`:** higher order functions such as `//.fn.fix` and `//.fn.fixt`.
2. **`//./path`** provides access to other arrai files relative to the current
   arrai file's parent directory (current working directory for expressions such
   as the `arrai eval` source that aren't associated with a file).
3. **`///path`** provides access to other arrai files relative to the root of
   the current module (TODO: explain modules).
4. **`//hostname/path`** provides access to arrai files in remote packages,
   e.g.: `//github.com/foo/bar`.

### Tuples vs Maps

It may not be immediately obvious why tuples and maps exist as distinct kinds of
values. Firstly, there is a practical reason: maps can have any kind of value as
keys:

```text
{
   "x":                 "red",
   [1, 2]:              "green",
   (a: [3], b: {5, 6}): "blue",
}
```

A more important distinction is that tuples should be used to capture various
known dimensions of a concept, whereas maps are more appropriate to map from an
arbitrary or unbounded set of values to some associated values. For example, a
collection of cars by license plate should be modeled as a map, since the set of
license plates is unbounded. The details of each car, however, form a closed set
of known attributes, which should be expressed as tuples:

```text
# Map
{
   "ILVME-23": (        # Tuple
      make:  "Porsche",
      model: "911",
      year:  1964,
   ),
   "ZUM-888": (         # Tuple
      make:  "Bugatti",
      model: "Veyron",
      year:  2005,
   ),
}
```

## Macros

Arr.ai has a macro system. The following example expresses a URL as a strongly
typed value:

```bash
$ arrai e '//.web.url{https://me@foo.com/bar?x=42}'
(
   source: "https://me@foo.com/bar?x=42",
   scheme: "https",
   authority: (
      userinfo: [8▸"me"],
      host: 11▸"foo.com",
   ),
   path: [19▸"bar"],
   search: {23▸"x": [25▸"42"]},
)
```

Another example is representing JSON:

```bash
$ arrai e '//.encoding.json{{"x": 1, "y": [2, 3], "z": null}}'
{
   "x": 1,
   "y": [2, 3],
   "z": (),
}
```

(Arr.ai has no counterpart for JSON null, so it uses the empty tuple as a
proxy.)

Macros are invoked via the syntax `macro{content}`. The content inside the macro
invocation is subject to a grammar defined by the macro itself, not regular
arr.ai syntax. Each macro can support its own grammar for the kind of content it
supports.

## Grammars

Arr.ai supports encoding of grammars directly inside the language. These
grammars may then be used to parse other content.

Example:

```bash
$ arrai e '//.grammar.lang.wbnf{expr -> @:[-+] > @:[/*] > \d+;}{1+2*3}'
("": [+], @rule: expr, expr: [(expr: [("": 1)]), ("": [*], expr: [("": 2), ("": 3)])])
```

The primary use of grammars is in the macro system. However, grammars are
themselves data structures, and can be transformed as such, allowing interesting
additions such as compositing, subsetting and otherwise transforming grammars.