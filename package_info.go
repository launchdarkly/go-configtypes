/*
Package configtypes provides types that may be useful in configuration parsing and validation, or
for other purposes where primitive types do not enforce the desired semantics.

# Opt types

The types beginning with "Opt" follow a general contract as described below; there are only a few
special cases that are documented for the specific types.

# Defined versus empty

The "Opt" types represent optional values of some wrapped type. Some of them allow any such value.
Others have additional validation, indicated by extra words after the name of the type: for
instance, OptIntGreaterThanZero is like OptInt except the value must be greater than zero.

Any instance can be in either a "defined" or an "empty" state. Defined means that it contains a
value of the wrapped type. Empty means no value was specified. The IsDefined() method tests this
state.

Instances have a GetOrElse() method that returns the wrapped value if defined, or the specified
alternative value if empty. The only exception to this is if there is already a natural way for
a value of the wrapped type to be undefined (for instance, if it is a pointer), in which case the
method is just Get().

If the wrapped type has mutable state (a slice, pointer, or map), it is always copied when
accessed.

# Opt type constructors

The zero value of these types is the empty state, so the way to declare a value in that state is
simply as an empty struct, OptFoo{} (all of these types are structs).

For types without validation, the NewOptFoo(valueType) constructor returns a defined OptFoo for
any given value.

	opt := OptBool{}          // opt is empty
	opt := NewOptBool(false)  // opt contains the value false
	opt := NewOptBool(true)   // opt contains the value true

For types with validationk the NewOptFoo constructor instead returns (optType, error). It will
never return an instance that wraps an illegal value, and there is no way for code outside of this
package to construct such an instance (without using tricks like reflection).

	opt, err := NewOptIntGreaterThanZero(3)   // err is nil, opt contains 3
	opt, err := NewOptIntGreaterThanZero(-3)  // err is non-nil, opt is empty

# Converting Opt types to or from strings

If the wrapped type is not string, there is also a NewOptFooFromString(string) constructor to
convert a string to this type. Again, if it is impossible for this conversion to fail, the
constructor returns just an instance; if it can fail, it returns (optType, error). Since Opt
types can always be empty, an empty string is always valid.

	opt, err := NewOptBoolFromString("")     // err is nil, opt is empty
	opt, err := NewOptBoolFromString("true") // err is nil, opt contains the value true
	opt, err := NewOptBoolFromString("bad")  // err is non-nil, opt is empty

The UnmarshalText method (encoding.TextUnmarshaler) behaves identically to the ...FromString
constructor. The encoding.TextUnmarshaler interface is recognized by many packages that do file
parsing (such as gcfg), so all such packages will automatically implement the appropriate
behavior for these types.

	var opt OptBool
	err := opt.UnmarshalText([]byte("true")) // err is nil, opt contains the value true

The String and MarshalText methods do the reverse, returning an empty string if empty or else
a string that is in the same format used by the parsing methods.

# Converting Opt types to or from JSON

These types also implement the json.Marshaler and json.Unmarshaler interfaces. An empty value
always corresponds to a JSON null; otherwise, the JSON mapping depends on the type, so for
instance a non-empty OptBool is always a JSON boolean.

# Opt types with multiple values

Some types, such as OptStringList, represent a collection of values. How this translates to
a text format depends on the type. The standard behavior is that if a parsing framework
uses encoding.TextUnmarshaler, it assumes that it will call UnmarshalText repeatedly for each
value with the same name (as gcfg does). In a context where it is not possible to have more
than one value with the same name (such as environment variables), the type may implement the
optional SingleValueTextUnmarshaler interface to indicate that an alternate format should be
used, such as a comma-delimited list.

# Reading from files or variables

Two common use cases are parsing a configuration file and reading values from environment variables.

The former can be done with any package that can interact with types via the encoding.TextUnmarshaler
interface, such as https://github.com/launchdarkly/gcfg. Attempting to set a field's string value
with this interface causes parsing to fail if the string format is not valid for the field's type,
for example if a string that is not an absolute URL is specified for a field of type OptURLAbsolute.

The VarReader type adapts the same functionality, but reads values from environment variables (or
from a name-value map). You can read values one at a time, specifying each variable name, or you can
use field tags to specify the variable names directly within the struct.

Since field names that do not appear in the parsed file or in the environment variables are not
modified, you can use both of these methods together: that is, read a configuration file that sets
some fields in a struct, and then allow environment variables to override other fields.

There is a limited ability to enforce that a field must have a value. Go has no way to prevent a
field or variable from being declared with a zero value for its type, so a struct with a required
field could always exist in an invalid state, but the Validate() function and VarReader will both
raise errors if a field that has a ",required" field tag was not set.
*/
package configtypes
