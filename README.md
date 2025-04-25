# A Generic and Extensible Configuration and Credential Framework for Go

This module provides some generic and extensible  frameworks for
- managing configurations by bringing together configuration providers and configuration consumers is a generic and uniform way.
  Any kinds of configuration settings can be combined ans consistently
  forwarded to appropriate configuration targets.  This way, although any used library can use own configuration settings, the zusing program can use a uniform way to configure all those different used environments.
- managing credentials by bringing together arbitrary credential sources, like config files or a vault credential repository and any kind of credential consumer. Hereby, programs don't nod to bother about reading various credential sources and forwarding appropropriate credentials to used libraries. They just have to use the crednetial mangement and pass a generic credential context to potential credential cosumer. This context then is able to provide appropriate crednetials according to the needs of the consumer.
- data contexts responsible for a particular technical realm bundling the  access API for this realm with configuration settings.

A tour demonstrates how to use those frameworks:
- [configuration contexts](examples/working-with-config/README.md)
