# token-replace

A custom client for translating tokenized configuration files into configs with passwords from Hashicorp Vault.  For those applications that won't read secrets from their environments.

Basically turns this:

```
PASSWORD=%%%vault/path/of/secrets:JSON_KEY%%%
```

into:

```
PASSWORD=thepasswordyouwant
```
