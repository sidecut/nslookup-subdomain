# nslookup-subdomain

See <https://golang.org/pkg/net/#hdr-Name_Resolution> for details on using environment variables to force use of the golang resolver,
which will return more than one domain name.

Example:

```bash
export GODEBUG=netdns=go    # force pure Go resolver
```

## TODO

[ ] Allow forcing golang resolver programmatically.
