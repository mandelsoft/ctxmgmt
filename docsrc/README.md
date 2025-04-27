# A Generic and Extensible Configuration and Credential Framework for Go

This module provides some generic and extensible  frameworks for
- managing configurations by bringing together configuration providers and configuration consumers in a generic and uniform way.
  Any kinds of configuration settings can be combined and consistently
  forwarded to appropriate configuration targets.  This way, although any used library may use own configuration settings, the using program can use a single uniform way to configure all those different used environments.
- managing credentials by bringing together arbitrary credential sources, like config files or a vault credential repository and any kind of credential consumer. Hereby, programs don't nod to bother about reading various credential sources and forwarding appropriate credentials to used libraries. They just have to use the credential management and pass a generic credential context to potential credential consumer. This context then is able to provide appropriate crednetials according to the needs of the consumer.
- data contexts responsible for a particular technical realm bundling the access API for this realm with configuration settings. Both, the configuration and the credential management are implemented ad data contexts using this framework. Hereby, the credential management context uses and incorporates a configuration context for its configuration needs

## Example 

Using the provided frameworks allows to bundle
all required configuration in one central configuration file.
It may contain one or any number of arbitrary, but typed
configuration data structures in form of a YAML document.
In our example this is a `.appconfig` file

```yaml
{{include}{../examples/working-with-credentials/.appconfig}}
```

The credential management offers a config object of type
`{{execute}{go}{run}{./info}{credentials config type}}`. 
It can be used to configure supported credential repositories, here a docker config file. Arbitrary credential repositories types can be 
added by using libraries. Out-of-the box, HasiCorp Vault, Docker and NPM config files are supported. Additionally,
explicit credential setting can be configured, here, credentials for a consumer type `{{execute}{go}{run}{./info}{consumer type}}` 
used in out example.

A second config object (`{{execute}{go}{run}{./info}{my config type}}`), implemented by this example,
covers the configuration parameters for our demo application.

The application program now does not need to bother
with config or credential handling anymore.
When starting, it just uses the config file to configure
the used credential context:

```go
{{include}{../examples/working-with-credentials/00-application-scenario.go}{main}}
```

This context is then passed around, and can be used by all participants to gain access to the required information.

The application constructor uses an own configuration structure, which can be served by our application config object. Calling `ApplyAllTo` applies all pending
configuration objects applicable for the given configuration target.

```go
{{include}{../examples/application/application.go}{application}}
```

This way the address of the service instance to use is configured. But it does not need to care about the service client. It also gets access to the context and is responsible on its own to retrieve the required additional information from it.

It creates a specification for the required credentials, consisting of a consumption type (our `{{execute}{go}{run}{./info}{consumer type}}`) and additional fields concretely describing the scenario, here, the host and port name.


```go
{{include}{../examples/service/identity/identity.go}{consumer id}}
```

Using this specification it just asks the credential management to determine the best matching credential setting extracted from one of the configured credential sources.
The credential context incorporated a configuration context, which is used to implicitly configure and update the used credential sources. The caller does not need to care about this. It just asks for credentials meeting the given constraints.

```go
{{include}{../examples/service/service.go}{service}}
```

Finally, our main program just prints the configuration 
configured for the application and the service client.

```
{{execute}{go}{run}{../examples/working-with-credentials}{--dir}{../examples/working-with-credentials}{application}{<extract>}{app info}}
```

The complete example can be found in [examples/working-with-credentials/00-application-scenario.go](examples/working-with-credentials/00-application-scenario.go).


## The Tour
A tour demonstrates how to use those frameworks:
- [using configuration contexts]({{config}})
