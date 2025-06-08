# Notes

## What am I trying to accomplish here?

There are several fairly-common operations I don't love Go's default solutions to:

- Equality comparison, especially for types that rely heavily on getter functions.
  - Arguably you shouldn't be using getters so much on Go methods ... it seems more-idiomatic to rely
    on public members for this kind of stuff, at least on types that you want to be able to directly
    compare.
  - On the other hand, I also don't want to have a lot of impediments to the getter pattern when it is
    the right thing to do for other reasons.
  - This lack of getter support is especially-acute if you want to compare interface values *as interfaces*
    and not do a full comparison of the underlying types to each other.
  - But, again, this feels a bit non-idiomatic ... interfaces in Go are normally for behavioral things.
  - If we take the getters out of the equation for a minute, then I think the only case is situations where 
    a straight equality comparison won't work for whatever reason?  For example with floating-point numbers
    where the less-significant digits may just not work out the same ay all the time?
  - I suppose there are also cases where the structs that you want to compare contain interfaces as value
    fields ... that's a tough case because you don't just want to compare the underlying types ... the whole
    point of an interface is to abstract that away.
  - OK, so there may be some value in this, although I'm not sure I have concrete cases beyond things like
    floating point comparison.
  - Another interesting point is something like comparisons of structs, where maybe order doesn't matter among
    the values.  Situational, but not purely theoretical.
  - Or maybe situations where you really need pointer equality vs. equality of values pointed to by pointers.
  - The best example of a case where getters matter is probably something like a function that returns an
    iterator over values.
- Diffing one value against another, especially in the context of testing.
  - This is similar to the equality case, but the difference is that you want to have some method of reporting
    what is different.
  - I suppose teh same logic is possible, for example if you want to control over what fields really need to
    be included in the diffing process?
  - For simplicity sake, let's consider diffing & equality to be basically the same thing, with equality being
    a special case of not caring abut the exact nature of the differences.
  - This is a real, non-theoretical problem in the context of testing.
- Validation is another point.
  - I think the behavior that I want for validation is something along the lines of: a value is valid if all of
    its children (i.e. struct fields) are valid, and if any validation on the level of the type as a whole (for
    example relationships between fields) is valid.
  - I could adopt a pattern where types that can be validated adopt a `Validate()` method to return an error for
    invalid values.  The only real limitation on this is that it is a pain to determine (and manually write code
    for) validation of all child fields.
  - The way I've approached this in the past is to write a utility "validate children" function that can use
    reflection to iterate over the children of a particular value and validate them.
- String descriptions of values.
  - This is very similar to validation, in that you could _theoretically_ write code on all of your types to
    provide wonderful outputs, but it is rarely done in practice.
  - This seems like another case where access to getters would be useful, for example if there's some computed
    value based on the fields of a struct that is really interesting in addition to the fields themselves.
  - Even theoretically simple things like formatting each field on its own line & handling indentation are not
    super straightforward.
  - This also seems like a case (again similar to validation) where some kind of reflective exploration of child
    values could be really useful.
  - This also seems like a case where some options would go a long way ... for example: do you want single-line
    vs. multi-line output?  Do you want field names in the formatting?
- How about sorting?
  - The standard way to accomplish this is to write your own comparison function for types that you want to sort,
    probably based on `cmp.Compare` or `cmp.Less` under the hood.
  - That might not be so bad, especially assuming that you don't have that many types that you need to sort.
  - It does seem like there's potential problems here, especially if you need to sort across a large number of
    of fields.
  - However, I think I could get a lot of mileage out of just writing a chain comparison function to break ties,
    and I suspect that (again) I'll have few-enough types that actually require this kind of sorting.
  - This also seems like a place where code generation might be a really strong solution.
- Another interesting testing case might be something like matchers.
  - The idea being that strict equality might
    be too strong of an assertion in some tests, and you might only care about (for example) a subset of the fields
    from a struct, or maybe you only care about a slice's length vs. its actual values.
  - This is logically-similar to the diffing case above, except in this case you only create a skeleton of the type
    that you want to match against.
  - This is a cool idea, especially given that you could build up composite matchers that mirror the nested nature
    of a type ... but is also sounds really complicated.