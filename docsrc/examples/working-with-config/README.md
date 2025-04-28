# Working with Configurations
{{config}}

This tour illustrates the basic configuration management. The library provides
an extensible framework to bring together configuration settings
and configuration target.

It covers the following basic scenarios:

- [`basic`]({{basic}}) Basic configuration management illustrating a simple configuration use case.
- [`generic`]({{generic}}) Handling of arbitrary configuration data.
- [`central`]({{central}}) Using central configuration files.
- [`configset`]({{config-sets}}) Using preconfigured config sets.
- [`write`]({{write-config}}) Providing new Config Object Types
- [`consume`]({{consume-config}}) Preparing Objects to be Configured by the Config Management
- [`applier`]({{config-appliers}}) Using Config Appliers for Generic Config Data Providers

## Running the example

You can call the main program with the name of the scenario as argument under `examples/working-with-config`.

## General Architecture

The configuration management can be found in sub package `config`.
Configuration is managed by a configuration context. This context is configured with
a set of context object types. A config type has a type name and provides a deserialization
of an appropriately types text representation of the setting of the config object.
With the type taken from the serialization format the context is able to transform
the content again to a configuration object.

A configuration object contains some configuration data. Its task is to apply these settings
to supported configuration target objects, for example other contexts.
A configuratuion object can be responsible for completely different target objects.

Configuration objects are then applied to a configuration context, which keeps a queue
of applied objects. A configurable target object then requests its configuration from
the configuration context.

It is important to keep in mind that this reverts the configuration direction. Applying
a configuration object to the configuration management does NOT configure
applicable objects. Instead, an object, which wants to be configured must ask the config
management to do so. THis has been designed tis way to keep the config management free
from references to any configurable object. This enables the Go garbage collection
to delete such objects if they are not used anymore. The config management does not
prevent such garbage collection. Both parts, the configuration context and the configurable
object are independently garbage collectable.

If an object want to be configured automatically, it requires a reference to the
configuration context (see [consumption tour]({{consume-config}}))

The context object supports incremental updates by providing
a target specific `Updater` object, which keeps track of already applied configuration
objects. This way, during runtime further configuration objects can be applied. The
configurable target object can request the context for those incremental updates using
its `Updater` object.

## Walkthrough

### Basic Configuration Management
{{basic}}

The complete example code can be found in [examples/working-with-config/01-basic-config-management.go](01-basic-config-management.go).

To use the configuration management an appropriate configuration context object
is required. The most simple way to achieve such an object is to use
the default context or a separate context based on the default context.

The default context is initialized with all the configuration object types
known to the library or attached later by a using library.

```go
{{include}{../../../examples/working-with-config/01-basic-config-management.go}{default context}}
```

The configuration context handles configuration objects.
A configuration object is any object implementing
the `config.Config` interface. The task of a config object
is to apply configuration to some target object.

One such object is our example config object defined in `examples/myconfig`.

```go
{{include}{../../../examples/working-with-config/01-basic-config-management.go}{my config}}
```

Here, we can specify some settings for out example application.
Typically, those config objects can be serialized to and deserialized from 
a YAML-based text representation.

```go
{{include}{../../../examples/working-with-config/01-basic-config-management.go}{marshal}}
```

The serialized format look as follows:

```yaml
{{execute}{go}{run}{../../../examples/working-with-config}{basic}{<extract>}{format}}
```

Like all the other manifest based descriptions this format always includes
a type field, which can be used to deserialize a specification into
the appropriate object.
This can be done by the config context, which keeps a set of known (registered)
config object types. It accepts YAML or JSON.

```go
{{include}{../../../examples/working-with-config/01-basic-config-management.go}{unmarshal}}
```

Regardless what variant is used (direct specification object or descriptor)
the config object can be added to a config context.

```go
{{include}{../../../examples/working-with-config/01-basic-config-management.go}{apply config}}
```

Every config object implements the
`ApplyTo(ctx config.Context, target interface{}) error` method.
It takes an object, which wants to be configured.
The config object then decides, whether it provides
settings for the given object and calls the appropriate
methods on this object (after a type cast).

Here is the code snippet from the apply method of our config object
config object ([examples/myconfig/type.go](../../../examples/myconfig/config.go)):

```go
{{include}{../../../examples/myconfig/config.go}{method apply}}
        ...
```

As can be seen a config object may support multiple config targets, and
configure those with different parts of its settings. Here, the
credentials are used to configure the credentials context, and the
service address to configure an object supporting the 
`{{execute}{go}{run}{../../info}{config target}}`
interface.


This way the config mechanism reverts the configuration
request, it does not actively configure something, instead
an object, which wants to be configured calls the config
context to apply pending configs.
To do this the config context manages a queue of config objects
and applies them to an object to be configured.

No we can request the configuration of an object by calling
the ` ApplyAllTo` method on the configuration context
for this object. Here, we use an `ExampleTarget` object which
implements the target interface of our config object.

```go
{{include}{../../../examples/working-with-config/01-basic-config-management.go}{configure}}
```

The result is then

```
{{execute}{go}{run}{../../../examples/working-with-config}{basic}{<extract>}{result}}
```

<!-------------------------------------------------------------------------------->

### Handling of Arbitrary Configuration
{{generic}}

The complete example code can be found in [examples/working-with-config/02-handle-arbitrary-config.go](02-handle-arbitrary-config.go).

The config management not only manages configuration objects for any
other configurable object, it also provides a configuration object of
its own. The task of this object is to handle other configuration objects
to be applied to a configuration object. The target of this configuration
object is therefore the configuration context itself.

```go
{{include}{../../../examples/working-with-config/02-handle-arbitrary-config.go}{config config}}
```

The generic config object holds a list of any other config objects,
or their specification formats.

We recycle our application config from the last example to get
a config object to be added to our generic config object.

```go
{{include}{../../../examples/working-with-config/02-handle-arbitrary-config.go}{sub config}}
```

Now, we can add this config object as nested configuration to
our generic config list.

```go
{{include}{../../../examples/working-with-config/02-handle-arbitrary-config.go}{add config}}
```

As we have seen in our previous example, config objects are typically
serializable and deserializable. This also holds for the generic config
object of the config context.

```go
{{include}{../../../examples/working-with-config/02-handle-arbitrary-config.go}{serialized}}
```

The result is a config object hosting a list (here with 1 entry)
of other config object specifications.

```
{{execute}{go}{run}{../../../examples/working-with-config}{generic}{<extract>}{result}}
```

The generic config object can be added to a config context, again, like
any other config object. If it is asked to configure a configuration
context it uses the methods of the configuration context to apply the
contained list of config objects.
Therefore, all config objects applied to a configuration context are
asked to configure the configuration context itself when queued to the
list of applied configuration objects.

```go
{{include}{../../../examples/working-with-config/02-handle-arbitrary-config.go}{apply config}}
```

If we now ask the context, again, to configure our application,
the nested config applies its value.

```go
{{include}{../../../examples/working-with-config/02-handle-arbitrary-config.go}{configure}}
```

The result is the same as in the previous example

```
{{execute}{go}{run}{../../../examples/working-with-config}{generic}{<extract>}{result}}
```


The very same mechanism is used to provide central configuration in a
configuration file for the OCM ecosystem, as will be shown in the next example.

<!-------------------------------------------------------------------------------->

### Central Configuration
{{central}}

The complete example code can be found in [examples/working-with-config/03-using-central-config.go](03-using-central-config.go).

Although the configuration of a context can
be done by a sequence of explicit calls to a configuration context
according to the mechanism
shown in the examples before, it provides a simple
library function, which can be used to configure a
context and all related other contexts with a single call
based on arbitrary central configuration files.

```go
{{include}{../../../examples/working-with-config/03-using-central-config.go}{central config}}
```

Here, we use the file `.appconfig`.
This file typically contains the serialization of such a generic
configuration specification shown in the previous example 
(or any other serialized configuration object),
enriched with specialized config specifications for
credentials, default repositories, or any
other configuration specifications.

Now, we can again configure our example target as shown before:

```go
{{include}{../../../examples/working-with-config/03-using-central-config.go}{configure}}
```

and achieve the same result as before.

```
{{execute}{go}{run}{../../../examples/working-with-config}{--dir}{../../../examples/working-with-config}{central}{<extract>}{result}}
```

#### Standard Configuration Files

Most important for such configurations are the credentials.
Because a library or application might embrace lots of network based remote
services or storage technologies, there are typically multiple technology
specific ways to configure credentials for command line tools.
Using the credentials settings shown in the [next-tour]({{working-with-credentials}}),
this library also provide a credential management. It is based on
this config management and provides an own config object.
It is possible to specify credentials for all
required purposes, and the configuration management provides
an extensible way to embed native technology specific ways
to provide credentials just by adding an appropriate type
of credential repository, which reads the specialized storage and
feeds it into the credential context. Those specifications
can be added via the credential configuration object to
the central configuration.

One such repository type is the Docker config type. It
reads a `dockerconfig.json` file and feeds in the credentials
to be used for OCI registries.

We first create the specification for a new credential repository of
type `dockerconfig` describing the default location
of the standard Docker config file.

```go
{{include}{../../../examples/working-with-config/03-using-central-config.go}{docker config}}
```

By adding the default location for the standard Docker config
file, all credentials provided by the `docker login` command
are available for our program, also. This way any required technology
can be embedded and this module provides a uniform way to access
this information regardless of its source.

A typical minimal <code>.appconfig</code> file can be composed as follows.
We add this config object to an empty generic configuration object
and print the serialized form. The result can be used as
default initial configuration file.

```go
{{include}{../../../examples/working-with-config/03-using-central-config.go}{default config}}
```

The result should looks as follows:

```yaml
{{execute}{go}{run}{../../../examples/working-with-config}{--dir}{../../../examples/working-with-config}{central}{<extract>}{config file}}
```

Because of the ordered map keys the actual output looks a little bit confusing: 
don't be worried about the location of the `type` field.


Besides from a file, such a config can be provided as data, also,
taken from any other source, for example from a Kubernetes secret.

```go
{{include}{../../../examples/working-with-config/03-using-central-config.go}{by data}}
```

<!-------------------------------------------------------------------------------->

#### Templating

The configuration library function does not only read the
config file, it also applies [*spiff*](https://github.com/mandelsoft/spiff)
processing to the provided YAML/JSON content. *Spiff* is an
in-domain yaml-based templating engine. Therefore, you can use
any spiff dynaml expression to define values or even complete
sub structures.

```go
{{include}{../../../examples/working-with-config/03-using-central-config.go}{spiff}}
```

This config object is not directly usable, because the cert value is not
a valid certificate. We use it here just to generate the serialized form.

```yaml
{{execute}{go}{run}{../../../examples/working-with-config}{--dir}{../../../examples/working-with-config}{central}{<extract>}{spiffconfig}}
```

If this is used with one of the above library functions, the finally generated
config object will contain the read file content, which is hopefully a
valid certificate.

<!-------------------------------------------------------------------------------->

### Working with Config Sets
{{config-sets}}

The complete example code can be found in [examples/working-with-config/04-config-sets.go](04-config-sets.go).

A `ConfigSet` represents, like configuration objects of the configuration context, a list of
configuration objects, but it is not itself a configuration object. Instead, it can be used
to configure named lists of configuration objects at a configuration context.

We configure our application config object, here:

```go
{{include}{../../../examples/working-with-config/04-config-sets.go}{config set}}
```

When adding a set to a context it gets assigned a name:

```go
{{include}{../../../examples/working-with-config/04-config-sets.go}{add config set}}
```

Once added to a context, sets can then be activated with a single API call by their assigned name.

```go
{{include}{../../../examples/working-with-config/04-config-sets.go}{activate}}
```

Now, the included config objects are applied to the configuration context, and
our application can be configured as already shown in the previous examples

```go
{{include}{../../../examples/working-with-config/04-config-sets.go}{configure}}
```

and provides the expected result:

```
{{execute}{go}{run}{../../../examples/working-with-config}{configset}{<extract>}{result}}
```

Config sets can not only be configured via the API, but the config object type
of the configuration management, also. Here, we configure a new set called `modified`.

```go
{{include}{../../../examples/working-with-config/04-config-sets.go}{configuring sets}}
```

This way preconfigured sets can be provided by central configuration sources and
enabled just by an API call using the set name.

But even the activation of previously defined sets can be handled by ut config object.

After activating the `modified` set with

```go
{{include}{../../../examples/working-with-config/04-config-sets.go}{activate modified}}
```

and reconfiguring our config target, the result is adapted, accordingly:

```
{{execute}{go}{run}{../../../examples/working-with-config}{configset}{<extract>}{modified result}}
```

<!-------------------------------------------------------------------------------->

### Providing new Config Object Types
{{write-config}}

The complete example code can be found in [examples/working-with-config/05-write-config-type.go](05-write-config-type.go).

So far, we just showed how to use config types to configure objects.
But the configuration management is highly extensible, and it is quite
simple to provide new config types, which can be used to configure
any new or existing object, which is prepared to consume configuration.

The next [chapter]({{consume-config}}) will show more elaborated way to
consume configuration and how to prepare an
object to be automatically configurable by
the configuration management. 

Now, we will show how new configuration object types can be
implemented and registered to be usable by the configuration management.

#### The Configuration Object Type

Typically, every kind of configuration object lives in its own package,
which always have the same layout. Our application congig example
object is defined in package [myconfig](../../../examples/myconfig/config.go)

A configuration object has a *type*, the configuration type. Therefore,
the package declares a constant `TYPE`. Here, we use the name `{{execute}{go}{run}{../../info}{my config type}}`.

It is the name of our new configuration object type.
To be globally unique, it should always end with a
DNS domain owned by the provider of the new type.

```go
{{include}{../../../examples/myconfig/config.go}{type name}}
```

Next, we need a Go type. `{{execute}{go}{run}{../../info}{config object struct}}` is the new Go type for the
config specification covering our example configuration. Because every config
object type uses its own package, always the same generic name can be used.
It just encapsulates our application configuration values. We use, the server
address and optionally credentials (basically this is not a good idea, because it 
could already be configured with other config objects. We just added this for
demonstration puropses)
used to configure the examples of our tour.

```go
{{include}{../../../examples/myconfig/config.go}{config type}}
```

Every config type structure must contain a field (and the appropriate methods)
for storing the config type name. This is done by embedding the
type `runtime.ObjectVersionedType` from the `runtime` package. This package
contains everything to work with specification objects and
serialization/deserialization.

Additional fields describe our desired configuration values.

A config type typically provides a constructor for a config object of
this type:

```go
{{include}{../../../examples/myconfig/config.go}{constructor}}
```

Additional setters can be used to configure the configuration object,
for example by mapping more complex objects into configuration attributes
prepared to be serialized.

```go
{{include}{../../../examples/myconfig/config.go}{setters}}
```

Getters can be used to prepare configuration attributes and make them available
to the application by converting values into more complex Go abstractions finally
usable by the application

```go
{{include}{../../../examples/myconfig/config.go}{getters}}
```

The utility function `runtime.CheckSpecification` can be used to
check a byte sequence to be a valid specification.
It just checks for a valid YAML document featuring a non-empty
`type` field:

```go
{{include}{../../../utils/runtime/utils.go}{check}}
```

The most important method to implement is `ApplyTo(_ cpi.Context, tgt interface{}) error`,
which must be implemented by all configuration objects.
Its task is to apply the described configuration settings to a dedicated
object.

```go
{{include}{../../../examples/myconfig/config.go}{method apply}}
```

It is free to configure any type of object, even multiple ones.
Therefore, it decides, whether it is able to handle a dedicated type of target
object and how to configure it. This way a configuration object
may apply is settings or even parts of its setting to any kind of target object.

Our configuration object supports two kinds of target objects:

If the target is a credentials context it configures the credentials to be
  used for tour application service (see also the [credential management tour]({{working-with-credentials}}).

But we want to accept more types of target objects. Therefore, we
introduce an own interface declaring the methods required for applying
some configuration settings (here, our service address).

```go
{{include}{../../../examples/myconfig/config.go}{target interface}}
```

By checking the target object against this interface, we are able
to configure any kind of object, as long as it provides the necessary
configuration methods.

Now, we are nearly prepared to use our new configuration, there is just one step
missing. To enable the automatic recognition of our new type (for example
in a central config file), we have to tell the configuration management
about the new type. This is done by an `init()` function in our config package.

Here, we call a registration function,
which gets called with a dedicated type object for the new config type.
A *type object* describes the config type, its type name, how
it is serialized and deserialized and some description.
We use a standard type object, here, instead of implementing
an own one. It is parameterized by the Go pointer type (`*{{execute}{go}{run}{../../info}{config object struct}}`) for
our specification object.

```go
{{include}{../../../examples/myconfig/config.go}{init}}
```

Without this registration the new config object could be used via API calls,
but the deserialization is not known by the configuration context. There, it could
not be used as part of a configuration file.

#### Using our new Config Object

After preparing a new special config type
we can feed it into the config management.

A usual, we gain access to our configuration context.

```go
{{include}{../../../examples/working-with-config/05-write-config-type.go}{default context}}
```

To setup our environment we create our new config based on the desired settings
and apply it to the config context.

```go
{{include}{../../../examples/myconfig/config.go}{method apply}}
```

Now, we should already be prepared to get the credentials configured
in our config object.

```go
{{include}{../../../examples/working-with-config/05-write-config-type.go}{query credentials}}
```

The default credential context uses the default configuration context
and is therefore able to determine the credentials for the requested
consumer scenario (more details about using the credential context are
shown in the [credential management tour]({{working-with-credentials}})).

```yaml
{{execute}{go}{run}{../../../examples/working-with-config}{write}{<extract>}{credentials}}
```


#### Using in the Central Configuration

Because of the registration of the new credential type, such a specification can
now be added to a serialized config file, also.

```go
{{include}{../../../examples/working-with-config/05-write-config-type.go}{in central config}}
```

The resulting config file looks as follows:

```yaml
{{execute}{go}{run}{../../../examples/working-with-config}{--dir}{../../../examples/working-with-config}{write}{<extract>}{config file}}
```

#### Applying to our Configuration Interface

Above, we added a new kind of target, the `{{execute}{go}{run}{../../info}{config target}}` interface.
By providing an implementation for this interface, we can directly request the
configuration of such an object using the config management.
We just provide a simple implementation for this interface, just storing the configured
repository specification.

```go
{{include}{../../../examples/working-with-config/01-basic-config-management.go}{target}}
```

The context management now is able to apply our config to such an object.

```go
{{include}{../../../examples/working-with-config/05-write-config-type.go}{request config}}
```

It finally sets the configured value:

```
{{execute}{go}{run}{../../../examples/working-with-config}{write}{<extract>}{result}}
```

This way any specialized configuration object can be added
by a user of the library. It can be used to configure
existing objects or even new object types, even in combination.

What is still required is a way
to implement new config targets, objects, which wants
to be configured and which autoconfigure themselves when
used. Our simple repository target is just an example
for some kind of ad-hoc configuration.
A complete scenario is shown in the next example.

<!-------------------------------------------------------------------------------->

{{consume-config}}
### Preparing Objects to be Configured by the Config Management

The complete example code can be found in [examples/working-with-config/06-write-config-consumer.go](06-write-config-consumer.go).

We already have our new acme.corp config object type,
and a target interface which must be implemented by a target
object to be configurable. The last example showed how
such an object can be configured in an ad-hoc manner
by directly requesting to be configured by the config
management.

Now, we want to provide an object, which configures
itself when used.
Therefore, we introduce a Go type `ServiceAddressProvider`.
It should be an object, which is
able to provide an address of an acme service.
It has a setter and a getter (the setter is
provided by our ad-hoc `ExampleTarget` implementation for an
object configurable with an address by our config object).

To be able to configure itself, the object must know about
the config context it should use to configure itself.

Therefore, our type contains an additional field `updater`.
Its type `cpi.Updater` is a utility provided by the configuration
management, which holds a reference to a configuration context
and is able to
configure an object based on a managed configuration
watermark. It remembers which config objects from the
config queue are already applied, and replays
the config objects applied to the config context
since the last update.

Finally, a mutex field is contained, which is used to
synchronize updates, later.

```go
{{include}{../../../examples/working-with-config/06-write-config-consumer.go}{type}}
```

For this type a constructor is provided, which initializes
the `updater` field with the desired configuration context.

```go
{{include}{../../../examples/working-with-config/06-write-config-consumer.go}{constructor}}
```

The magic now happens in the methods provided
by our configurable object.
The first step for methods of configurable objects
dependent on potential configuration is always
to update itself using the embedded updater, before
some action is executed potentially affected by configuration changes.

Please note, the configuration management reverses the
request direction. Applying a config object to
the config context does not configure dependent objects,
it just manages a config queue, which is used by potential
configuration targets to configure themselves.
The actual configuration action is always initiated
by the object, which wants to be configured.
The reason for this is to avoid references from the
management to managed objects. This would prohibit
the garbage collection of all configurable objects
as long as the configuration context exists.

```go
{{include}{../../../examples/working-with-config/06-write-config-consumer.go}{method}}
```

After defining our provider type we can now start to use it
together with the configuration management and our configuration object.

As usual, we first determine out context to use.

```go
{{include}{../../../examples/working-with-config/06-write-config-consumer.go}{default context}}
```

Now, we create our provider as configurable object by binding it
to the config context.

```go
{{include}{../../../examples/working-with-config/06-write-config-consumer.go}{object}}
```

If we ask now for a service address we will get the empty
answer, because nothing is configured, so far.

```go
{{include}{../../../examples/working-with-config/06-write-config-consumer.go}{initial query}}
```

In the next step, we apply our config from the last example. Therefore, we create and 
initialize the config object with our desired settings and apply it to the config
context.

```go
{{include}{../../../examples/working-with-config/06-write-config-consumer.go}{apply config}}
```

Without any further action, asking for the provider now will return the
configured address. The configurable object automatically catches the
new configuration from the config context.

```go
{{include}{../../../examples/working-with-config/06-write-config-consumer.go}{query}}
```

It returns the configured service address:

```yaml
{{execute}{go}{run}{../../../examples/working-with-config}{consume}{<extract>}{result}}
```

Additionally, we should also be prepared to get the credentials,
our config object configures the provider as well as
the credential context.

```go
{{include}{../../../examples/working-with-config/06-write-config-consumer.go}{query credentials}}
```

This gives us the following output:

```yaml
{{execute}{go}{run}{../../../examples/working-with-config}{consume}{<extract>}{credentials}}
```

<!-------------------------------------------------------------------------------->


### Using Config Appliers for Generic Config Data Providers
{{config-appliers}}

The complete example code can be found in [examples/working-with-config/07-config-appliers.go](07-config-appliers.go).

A typical config object bundles the configuration logic
with particular configuration fields.
This could be problematic if generic data stores 
should be used to store configuration data. Here, we
have a third element, the technical access to the
data repository. This would require to create
config objects for every combination of
repository, config fields and config logic.
To circumvent this, the configuration management
supports a decoupling of the data access from application
logic by introducing `ConfigAppliers`.

A `ConfigApplier` is restricted to the pure configuration logic.
It takes arbitrary config data
and applies it according to its own logic.
The storage access to retrieve the data is left to
a storage technology specific implementation of
a config object just reading the data and applying
a named config applier.

Such a config object features properties to identify 
the config attributes in a data storage and the name
of the config applier.

Our demo applier expects a string as configuration value
and configures our well-known `{{execute}{go}{run}{../../info}{config target}}`
interface used to configure our demo application.


```go
{{include}{../../../examples/working-with-config/07-config-appliers.go}{applier}}
```

Like configuration object types, Config appliers are registered
at a configuration context. Default appliers are globally
registered using an `init`function and implicitly used for all
non-initial contexts.

```go
{{include}{../../../examples/working-with-config/07-config-appliers.go}{init}}
```

Appliers are used together with config objects providing generic
configuration data. For test purposes the library provides
one such config object types using arbitrary inline data for configuration.

```go
{{include}{../../../config/extensions/data/type.go}}
```

We create such a config object for our service address and our config applier:

```go
{{include}{../../../examples/working-with-config/07-config-appliers.go}{config data}}
```

Now, our config target can be configured as usual resulting in

```yaml
{{execute}{go}{run}{../../../examples/working-with-config}{applier}{<extract>}{result}}
```