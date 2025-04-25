# Working with Configurations

This tour illustrates the basic configuration management. The library provides
an extensible framework to bring together configuration settings
and configuration target.

It covers the following basic scenarios:

- [`basic`](/examples/lib/tour/04-working-with-config/01-basic-config-management.go) Basic configuration management illustrating a simple configuration use case.


## Running the example

You can call the main program with the name of the scenario under `examples/config`.

## General Architecture

Configuration is managed by a configuration context. This context is configured with
a set of context object types. A config type has a type name and provides a deserialization
of an appropriately types text representation of the setting of the config object.
With the type taken from the serialization format the context is able to transform
the content again to a configuration object.

A configuration object contains some configuration data. Its task is to apply these settings
to supported configuratuon target objects, for example other contexts.
A configuratuion object can be resposnible for completely different target objects.

Configuratuon objects are then applied to a configuration context, which keeps a queue
a applied objects. A configurable target object then requests its configuration from
the configuration context. The context object supports incremental updates by providing
a target specific `Updater` object, which keeps track of already applied configuration
objects. This way, during runtime further configuration objects can be applied. The
configurable target object can request the context for those incremental updates using
its `Updater` object.

## Walkthrough

### Basic Configuration Management

To use the configuration management an appropriate configuration context object
is required. The most simple way to achieve such an object is to use
the default context or a separate context based on the default context.

The default context is initialized with all the configuration object types
known to the library or attached later by a using library.

```go
	ctx := config.DefaultContext()
```

The configuration context handles configuration objects.
A configuration object is any object implementing
the `config.Config` interface. The task of a config object
is to apply configuration to some target object.

One such object is our example config object defined in `examples/myconfig`.

```go
	cfg := myconfig.NewConfig("service.provider.com")
	cfg.Credentials.Username = "appuser"
	cfg.Credentials.Password = "apppass"
```

Here, we can specify some settings for out example application.
Typically, those config objects can be serialized to and deserialized from 
a YAML-based text representation.

```go
	spec, err := json.MarshalIndent(cfg, "  ", "  ")
	if err != nil {
		return errors.Wrapf(err, "marshal config")
	}

	fmt.Printf("this is our configuration object:\n")
	helper.Output("format", func() {
		fmt.Printf("  %s\n", string(spec))
	})

```

The serialized format looks as follows:

```yaml
  {
    "type": "example.config.acme.org",
    "serviceAddress": "service.provider.com",
    "credentials": {
      "username": "appuser",
      "password": "apppass"
    }
  }
```

Like all the other manifest based descriptions this format always includes
a type field, which can be used to deserialize a specification into
the appropriate object.
This can be done by the config context, which keeps a set of known (registered)
config object types. It accepts YAML or JSON.

```go
	o, err := ctx.GetConfigForData(spec, nil)
	if err != nil {
		return errors.Wrapf(err, "deserialize config")
	}

	if diff := deep.Equal(o, cfg); len(diff) != 0 {
		fmt.Printf("diff:\n%v\n", diff)
		return fmt.Errorf("invalid des/erialization")
	}
```

Regardless what variant is used (direct specification object or descriptor)
the config object can be added to a config context.

```go
	err = ctx.ApplyConfig(cfg, "explicit setting")
	if err != nil {
		return errors.Wrapf(err, "cannot apply config")
	}
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
func (c *MyConfigObject) ApplyTo(_ cpi.Context, tgt interface{}) error {
	switch t := tgt.(type) {
	// if the target is a credentials context
	// configure the credentials to be used for the
	// described OCI repository.
	case credentials.Context:
		// determine the consumer id for our target repository.
		if c.Credentials.Username != "" && c.Credentials.Password != "" {
			id := identity.GetConsumerId(c.ServiceAddress)

			// create the credentials.
			creds := c.GetCredentials()

			// configure the targeted credential context with
			// the provided credentials (see previous examples).
			t.SetCredentialsForConsumer(id, creds)
		}

	// if the target consumes an OCI repository, propagate
	// the provided OCI repository ref.
	case ConfigTarget:
		t.SetServiceAddress(c.ServiceAddress)

	// all other targets are ignored, we don't have
	// something to set at these objects.
	default:
		return cpi.ErrNoContext(TYPE)
	}
	return nil
}

        ...
```

As can be seen a config object may support multiple config targets, and
configure those with different parts of its settings. Here, the
credentials are used to configure the credentials context, and the
service address to configure an object supporting the `ConfigTarget`
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
	tgt := &ExampleTarget{}

	err = ctx.ApplyAllTo(tgt)
	if err != nil {
		return errors.Wrapf(err, "request configuration")
	}
```

The result is then

```
configured address is "service.provider.com"
```
